package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	//lint:ignore ST1001 importing common shared across all dje8 cmds
	. "damien.live/dje8/pkg/common"
	"damien.live/dje8/pkg/ucodebuilder"
)

var OpCodes = [256]string{
	"NOP", //00
	"LDA", //01
	"ADD", //02
	"SUB", //03
	"STA", //04
	"LDI", //05
	"JMP", //06
	"JC ", //07
	"JZ ", //08
	"---", //09
	"---", //0a
	"---", //0b
	"---", //0c
	"---", //0d
	"OUT", //0e
	"HLT", //0f
}

var Program = []byte{ // Pre ASM Code
	0x05, 0x01, // 0 LDI 1
	0x0e,             // 1 OUT
	0x02, 0x00, 0x22, // 2 ADD 34
	0x07, 0x00, 0x20, // 3 JC 32
	0x0e,             // 4 OUT
	0x04, 0x00, 0x21, // 5 STA 33
	0x01, 0x00, 0x22, // 6 LDA 34
	0x02, 0x00, 0x21, // 7 ADD 33
	0x07, 0x00, 0x20, // 8 JC 32
	0x0e,             // 9 OUT
	0x04, 0x00, 0x22, // a STA 34
	0x01, 0x00, 0x21, // b LDA 33
	0x06, 0x00, 0x03, // c JMP 3
	0x0f, // d HLT 0
	0x01, // e 1
	0x01, // f 1
}

// define registers
var ProgramCounter uint16 = 0
var MemoryAddressRegister uint16 = 0
var StackPointer uint16 = 0
var InstructionRegister uint8 = 0
var AccumulatorRegister uint8 = 0
var InternalRegister uint8 = 0
var ArithmeticLogicUnit uint8 = 0
var FlagsRegister Flag = 0
var ControlWord Control = 0
var ClockPulse uint8 = 0

var AddressBus uint16 = 0
var DataBus uint8 = 0xff // data bus is pulled high when inactive (can be used as a source of -1)
var MemorySpace []byte
var ROMAddress uint16
var ControlROM []Control

var EmulationHeaderPaddingSize int = 14

func main() {
	fmt.Println("***** DJE-8 Emulator *****")
	fmt.Println()

	DebugPrintConsts()

	MemorySpace = make([]byte, 65536)

	ControlROM = ucodebuilder.BuildUcode() // make([]Control, 65536)

	// Test code
	// r := rand.New(rand.NewSource(time.Now().Unix()))
	// r.Read(MemorySpace) // Fill memory with 64K of randomness
	// for i := range len(ControlROM) { // Fill control ROM with randomness (except for HLT and Reserved bits)
	// 	ControlROM[i] = Control(r.Uint32()) &^ (HLT | EX3 | EX2 | EX1 | EX0)
	// }
	// End Test Code

	fmt.Println("***** DJE-8 Simulation Starting *****")
	copy(MemorySpace[0:], Program)
	fmt.Println()
	PrintEmulationHeaderPadding()

	// main Fetch-Decode-Execute loop
	for {
	InstructionCycle:
		for ClockPulse = range 16 {
			ControlWord = ControlROMLookup(ClockPulse)
			// TODO: check for interrupts
			// TODO: check bus arbiter
			// AddressBus = 0
			// DataBus = 0xff
			PrintSnapshot()
			fmt.Println()
			time.Sleep(time.Millisecond * 100)
			// fmt.Scanln()
			switch ClockPulse {
			case 0: // FETCH
				AddressBus = ProgramCounter        // CO
				MemoryAddressRegister = AddressBus // MI
			case 1: // DECODE
				DataBus = MemorySpace[MemoryAddressRegister] // RO
				InstructionRegister = DataBus                // II
				ProgramCounter++                             // CU
			default: // EXECUTE
				if ControlWord&HLT != 0 {
					fmt.Println("*** HALT signal received. System halted.")
					os.Exit(0)
				}

				// ***** OUT SIGNALS FIRST *****
				// Address Bus OUT Signals
				if ControlWord&CO2 != 0 {
					AddressBus = ProgramCounter
				}
				if ControlWord&PO2 != 0 {
					AddressBus = StackPointer
				}
				if ControlWord&RO2 != 0 {
					AddressBus = uint16(MemorySpace[MemoryAddressRegister])<<8 | uint16(MemorySpace[MemoryAddressRegister+1])
				}

				// Data Bus OUT Signals
				if ControlWord&AO != 0 {
					DataBus = AccumulatorRegister
				}
				if ControlWord&RO != 0 {
					DataBus = MemorySpace[MemoryAddressRegister]
				}
				// *** ALU Start
				var Mode ALUMode = ALUMode((ControlWord / AU0) & 0xf) // Getting just ALU flags in the lower half of one byte
				switch Mode {
				case NOP:
					// Any ALU mode other than NOP or CMP puts ALU contents on the data bus.
					// Only arithmentic ALU modes and CMP load the ZCNV flags.
					//
					// This implementation of the emulator simply performs the arithmetic and updates the flags only
					// on the cycles that the ALU bits are set even though the hardware implementation is likely
					// always performing a computation.
				case ADD:
					DataBus = AddAndSetFlags(AccumulatorRegister, InternalRegister, false)
				case SUB:
					DataBus = SubtractAndSetFlags(AccumulatorRegister, InternalRegister, false)
				case ADC:
					DataBus = AddAndSetFlags(AccumulatorRegister, InternalRegister, FlagsRegister&CarryFlagC != 0)
				case SBC:
					DataBus = SubtractAndSetFlags(AccumulatorRegister, InternalRegister, FlagsRegister&CarryFlagC != 0)
				case AND:
					DataBus = AccumulatorRegister & InternalRegister
				case OR:
					DataBus = AccumulatorRegister | InternalRegister
				case NOT:
					DataBus = ^AccumulatorRegister
				case NEG:
					DataBus = -AccumulatorRegister
					if DataBus == 0 {
						FlagsRegister |= ZeroFlagZ
					} else {
						FlagsRegister &= (^ZeroFlagZ)
					}
					if DataBus&0x80 != 0 {
						FlagsRegister |= NegativeFlagN
					} else {
						FlagsRegister &= (^NegativeFlagN)
					}
				case INC:
					DataBus = AccumulatorRegister + 1
					if DataBus < AccumulatorRegister {
						FlagsRegister |= CarryFlagC
						FlagsRegister |= OverflowFlagV
					} else {
						FlagsRegister &= (^CarryFlagC)
						FlagsRegister &= (^OverflowFlagV)
					}
					if DataBus == 0 {
						FlagsRegister |= ZeroFlagZ
					} else {
						FlagsRegister &= (^ZeroFlagZ)
					}
					if DataBus&0x80 != 0 {
						FlagsRegister |= NegativeFlagN
					} else {
						FlagsRegister &= (^NegativeFlagN)
					}
				case DEC:
					DataBus = AccumulatorRegister - 1
					if DataBus > AccumulatorRegister {
						FlagsRegister |= CarryFlagC
						FlagsRegister |= OverflowFlagV
					} else {
						FlagsRegister &= (^CarryFlagC)
						FlagsRegister &= (^OverflowFlagV)
					}
					if DataBus == 0 {
						FlagsRegister |= ZeroFlagZ
					} else {
						FlagsRegister &= (^ZeroFlagZ)
					}
					if DataBus&0x80 != 0 {
						FlagsRegister |= NegativeFlagN
					} else {
						FlagsRegister &= (^NegativeFlagN)
					}
				case CMP:
					_ = SubtractAndSetFlags(AccumulatorRegister, InternalRegister, false)
				}
				// *** ALU End

				// ***** IN SIGNALS NEXT *****
				// Address Bus IN Signals
				if ControlWord&CI2 != 0 {
					ProgramCounter = AddressBus
				}
				if ControlWord&MI2 != 0 {
					MemoryAddressRegister = AddressBus
				}
				if ControlWord&RI != 0 {
					MemorySpace[MemoryAddressRegister] = DataBus
				}

				// Data Bus IN Signals
				if ControlWord&AI != 0 {
					AccumulatorRegister = DataBus
				}
				if ControlWord&BI != 0 {
					InternalRegister = DataBus
				}
				if ControlWord&CIH != 0 {
					ProgramCounter = uint16(DataBus)<<8 | (ProgramCounter & 0x00ff)
				}
				if ControlWord&CIL != 0 {
					ProgramCounter = uint16(DataBus) | (ProgramCounter & 0xff00)
				}
				// if ControlWord&MIH != 0 {
				// 	MemoryAddressRegister = uint16(DataBus)<<8 | (MemoryAddressRegister & 0x00ff)
				// }
				// if ControlWord&MIL != 0 {
				// 	MemoryAddressRegister = uint16(DataBus) | (MemoryAddressRegister & 0xff00)
				// }
				if ControlWord&II != 0 {
					InstructionRegister = DataBus
				}

				// Increments and decrements
				if ControlWord&CU != 0 {
					ProgramCounter++
				}
				if ControlWord&CU2 != 0 {
					ProgramCounter += 2
				}
				if ControlWord&MU != 0 {
					MemoryAddressRegister++
				}
				if ControlWord&MU2 != 0 {
					MemoryAddressRegister += 2
				}
				if ControlWord&PU != 0 {
					StackPointer++
				}
				if ControlWord&PU2 != 0 {
					StackPointer += 2
				}
				if ControlWord&PD != 0 {
					ProgramCounter--
				}
				if ControlWord&PD2 != 0 {
					ProgramCounter -= 2
				}

				// End instruction cycle last
				if ControlWord&STR != 0 {
					break InstructionCycle
				}

			}
		}
	}

}

func ControlROMLookup(microStep uint8) Control {
	// Control ROM Address calc:
	// Z C N V OPCODEXX STEP
	// where:
	//	Z = zero flag
	//	C = carry flag
	//	N = negative flag
	//  V = overflow flag
	//  OPCODEXX = current 8-bit instruction op code
	//  STEP = current 4-bit microcode step

	ROMAddress = (uint16(InstructionRegister))<<4 | ((uint16(FlagsRegister))&0xf)<<12 | (uint16(microStep))
	return Control(ControlROM[ROMAddress])
}

func AddAndSetFlags(op1 uint8, op2 uint8, carryIn bool) uint8 {
	var result uint16 = uint16(op1) + uint16(op2)
	if carryIn {
		result++
	}
	// if the 9th bit of the result before truncation isn't 0, we carried.
	if result&0x0100 != 0 {
		FlagsRegister |= CarryFlagC
	} else {
		FlagsRegister &= (^CarryFlagC)
	}
	// the the MSb of the operands is the same and the MSb of the result is different, we overflowed the sign bit
	if AccumulatorRegister&0x80 == InternalRegister&0x80 && AccumulatorRegister&0x80 != DataBus&0x80 {
		FlagsRegister |= OverflowFlagV
	} else {
		FlagsRegister &= (^OverflowFlagV)
	}
	if result == 0 {
		FlagsRegister |= ZeroFlagZ
	} else {
		FlagsRegister &= (^ZeroFlagZ)
	}
	if result&0x80 != 0 {
		FlagsRegister |= NegativeFlagN
	} else {
		FlagsRegister &= (^NegativeFlagN)
	}

	return uint8(result)
}

func SubtractAndSetFlags(op1 uint8, op2 uint8, carryIn bool) uint8 {
	var result uint8 = op1 - op2
	if carryIn {
		result--
	}
	if ((^op1&op2)|(^(op1^op2)&result))>>7 != 0 {
		FlagsRegister |= CarryFlagC
	} else {
		FlagsRegister &= (^CarryFlagC)
	}
	// the the MSb of the operands is the same and the MSb of the result is different, we overflowed the sign bit
	if AccumulatorRegister&0x80 == InternalRegister&0x80 && AccumulatorRegister&0x80 != DataBus&0x80 {
		FlagsRegister |= OverflowFlagV
	} else {
		FlagsRegister &= (^OverflowFlagV)
	}
	if result == 0 {
		FlagsRegister |= ZeroFlagZ
	} else {
		FlagsRegister &= (^ZeroFlagZ)
	}
	if result&0x80 != 0 {
		FlagsRegister |= NegativeFlagN
	} else {
		FlagsRegister &= (^NegativeFlagN)
	}

	return uint8(result)
}

func PrintEmulationHeaderPadding() {
	for range EmulationHeaderPaddingSize {
		fmt.Println()
	}
}

func PrintSnapshot() {
	fmt.Printf("\033[%dA", EmulationHeaderPaddingSize)
	fmt.Printf("    PC:  0x%04x             A: 0x%02x (%3d)\n", ProgramCounter, AccumulatorRegister, AccumulatorRegister)
	fmt.Printf("    MAR: 0x%04x             B: 0x%02x (%3d)\n", MemoryAddressRegister, InternalRegister, InternalRegister)
	fmt.Printf("    IR:  0x%02x (%3s)  Step: 0x%x   F: %s (0x%02x)\n", InstructionRegister, OpCodes[InstructionRegister], ClockPulse, formatFlagByte(FlagsRegister), uint8(FlagsRegister))
	fmt.Printf("    ROM Lookup: 0b%016b (0x%04x)\n", ROMAddress, ROMAddress)
	fmt.Printf("    Control Wd: %s\n", formatControlWord(ControlWord))
	fmt.Print(formatControlWordLabels("                "))
	fmt.Println()
	fmt.Printf("    AddrBus: 0b%016b (0x%04x)  DataBus: 0b%08b\n", AddressBus, AddressBus, DataBus)
	fmt.Printf("    RAM:")
	for i := range 64 {
		fmt.Printf(" %02x", MemorySpace[i])
		if (i+1)%8 == 0 {
			fmt.Print(" ")
		}
		if (i+1)%32 == 0 {
			fmt.Print("\n        ")
		}
	}
	fmt.Println()
}

func formatFlagByte(Flags Flag) string {
	FlagsChars := []rune{}
	for i := HiBitFlag; i >= LoBitFlag; i = i >> 1 {
		Name := i.String()
		if Flags&i != 0 {
			FlagsChars = append(FlagsChars, ([]rune(Name))[:1]...)
		} else {
			FlagsChars = append(FlagsChars, '-')
		}
	}
	return string(FlagsChars)
}

func formatControlWord(ControlWord Control) string {
	s := fmt.Sprintf("%032b", ControlWord)
	runes := []rune{}
	for i, ch := range []rune(s) {
		runes = append(runes, ch, ' ')
		if (i+1)%8 == 0 {
			runes = append(runes, ' ')
		}
	}
	return string(runes)
}

func formatControlWordLabels(Prefix string) string {
	var retval strings.Builder
	ControlNames := [][]rune{}
	for i := HiBitControl; i >= LoBitControl; i = i >> 1 {
		ControlNames = append(ControlNames, []rune(Control(i).String()))
	}
	for i := range 3 {
		retval.WriteString(Prefix)
		for idx, runes := range ControlNames {
			if len(runes) > i {
				retval.WriteRune(runes[i])
			} else {
				retval.WriteRune(' ')
			}
			retval.WriteRune(' ')
			if (idx+1)%8 == 0 {
				retval.WriteRune(' ')
			}
		}
		retval.WriteRune('\n')
	}
	return retval.String()
}
