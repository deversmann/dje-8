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

var Program = []byte{ // Pre ASM Code
	byte(LODI), 0x01, // 0 LODI 1
	byte(NOP),              // 1 OUT(NOP)
	byte(ADDA), 0x00, 0x22, // 2 ADDA 34
	byte(BCS), 0x00, 0x20, // 3 BCS 32??
	byte(NOP),              // 4 OUT(NOP)
	byte(STOA), 0x00, 0x21, // 5 STOA 33
	byte(LODA), 0x00, 0x22, // 6 LODA 34
	byte(ADDA), 0x00, 0x21, // 7 ADDA 33
	byte(BCS), 0x00, 0x20, // 8 BCS 32??
	byte(NOP),              // 9 OUT(NOP)
	byte(STOA), 0x00, 0x22, // a STOA 34
	byte(LODA), 0x00, 0x21, // b LODA 33
	byte(JMP), 0x00, 0x03, // c JMP 3
	byte(HALT), // d HALT 0
	0x01,       // e 1
	0x01,       // f 1
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
			time.Sleep(time.Millisecond * 1)
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
				if ControlWord&COW != 0 {
					AddressBus = ProgramCounter
				}
				if ControlWord&POW != 0 {
					AddressBus = StackPointer
				}
				if ControlWord&ROW != 0 {
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
				case ALUNOP:
					// Any ALU mode other than NOP or CMP puts ALU contents on the data bus.
					// Only arithmentic ALU modes and CMP load the ZCNV flags.
					//
					// This implementation of the emulator simply performs the arithmetic and updates the flags only
					// on the cycles that the ALU bits are set even though the hardware implementation is likely
					// always performing a computation.
				case ALUADD:
					DataBus = AddAndSetFlags(AccumulatorRegister, InternalRegister, false)
				case ALUSUB:
					DataBus = SubtractAndSetFlags(AccumulatorRegister, InternalRegister, false)
				case ALUADC:
					DataBus = AddAndSetFlags(AccumulatorRegister, InternalRegister, FlagsRegister&CarryFlagC != 0)
				case ALUSBC:
					DataBus = SubtractAndSetFlags(AccumulatorRegister, InternalRegister, FlagsRegister&CarryFlagC != 0)
				case ALUAND:
					DataBus = AccumulatorRegister & InternalRegister
				case ALUOR:
					DataBus = AccumulatorRegister | InternalRegister
				case ALUNOT:
					DataBus = ^AccumulatorRegister
				case ALUNEG:
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
				case ALUINC:
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
				case ALUDEC:
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
				case ALUCMP:
					_ = SubtractAndSetFlags(AccumulatorRegister, InternalRegister, false)
				}
				// *** ALU End

				// ***** IN SIGNALS NEXT *****
				// Address Bus IN Signals
				if ControlWord&CIW != 0 {
					ProgramCounter = AddressBus
				}
				if ControlWord&MIW != 0 {
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
				if ControlWord&CUW != 0 {
					ProgramCounter += 2
				}
				if ControlWord&MU != 0 {
					MemoryAddressRegister++
				}
				if ControlWord&MUW != 0 {
					MemoryAddressRegister += 2
				}
				if ControlWord&PU != 0 {
					StackPointer++
				}
				if ControlWord&PUW != 0 {
					StackPointer += 2
				}
				if ControlWord&PD != 0 {
					ProgramCounter--
				}
				if ControlWord&PDW != 0 {
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
	fmt.Printf("    IR:  0x%02x (%4s)  Step: 0x%x   F: %s (0x%02x)\n", InstructionRegister, OpCode(InstructionRegister), ClockPulse, formatFlagByte(FlagsRegister), uint8(FlagsRegister))
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
