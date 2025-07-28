package common

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

//go:generate stringer -type=Flag
type Flag uint8

// Flags stored in the FlagsRegister
// Last character is the abbreviation of flag
const (
	OverflowFlagV Flag = 1 << iota
	NegativeFlagN
	CarryFlagC
	ZeroFlagZ
	InterruptFlagI
	ReservedFlag0_
	ReservedFlag1_
	ReservedFlag2_
	HiBitFlag = ReservedFlag2_ // *Const to enable traversing list
	LoBitFlag = OverflowFlagV  // *Const to enable traversing list
)

//go:generate stringer -type=Control
type Control uint32

// Control signals found in the ControlROM
const (
	HLT          Control = 0x80000000 >> iota // HALT
	AI                                        // Accumulator Register In
	AO                                        // Accumulator Register Out
	BI                                        // Internal Register In
	II                                        // Instruction Register In
	CI2                                       // Program Counter In from Address Bus
	CIL                                       // Program Counter LSB In from Data Bus
	CIH                                       // Program Counter MSB In from Data Bus
	CO2                                       // Program Counter Out
	CU                                        // Program Counter Increment
	CU2                                       // Program Counter Increment by 2
	MI2                                       // Memory Address Register In From Address Bus (removed MIL & MIH)
	MU                                        // Memory Address Register Increment
	MU2                                       // Memory Address Register Increment by 2
	RI                                        // Write to Memory
	RO                                        // Read from Memory to Data Bus
	RO2                                       // Read 2 bytes from memory to Address Bus
	PO2                                       // Stack Pointer Out
	PU                                        // Stack Pointer Increment
	PU2                                       // Stack Pointer Increment
	PD                                        // Stack Pointer Decrement
	PD2                                       // Stack Pointer Decrement
	AU3                                       // ALU Mode Bit 3
	AU2                                       // ALU Mode Bit 2
	AU1                                       // ALU Mode Bit 1
	AU0                                       // ALU Mode Bit 0
	FL                                        // Flags Register In
	EX3                                       // Reserved
	EX2                                       // Reserved
	EX1                                       // Reserved
	EX0                                       // Reserved
	STR                                       // Step Counter Reset
	HiBitControl = HLT                        // *Const to enable traversing list
	LoBitControl = STR                        // *Const to enable traversing list
)

//go:generate stringer -type=ALUMode
type ALUMode uint8

// The different ALU modes found in the ControlWord
// The first 8 are logical (don't affect status flags)
// and the second 8 are arithmetic (do affect status flags)
const (
	NOP       ALUMode = iota // 0000 no output
	AND                      // 0001 bitwise and
	OR                       // 0010 bitwise or
	NOT                      // 0011 bitwise inversion
	X3_                      // 0100 reserved
	X2_                      // 0101 reserved
	X1_                      // 0110 reserved
	X0_                      // 0111 reserved
	ADD                      // 1000 add
	SUB                      // 1001 subtract
	ADC                      // 1010 add considering carry flag
	SBC                      // 1011 subtract considering carry flag
	NEG                      // 1100 arithmetic negation
	INC                      // 1101 increment
	DEC                      // 1110 decrement
	CMP                      // 1111 compare
	LoALUMode = NOP          // *Const to enable traversing list
	HiALUMode = CMP          // *Const to enable traversing list
)

func DebugPrintConsts() {
	for i := HiBitFlag; i >= LoBitFlag; i = i >> 1 {
		// fmt.Printf("%-15s = %08b\n", i, i)
		printConstDebugString(i)
	}
	fmt.Println()
	for i := HiBitControl; i >= LoBitControl; i = i >> 1 {
		// fmt.Printf("%-15s = %032b\n", i, i)
		printConstDebugString(i)
	}
	fmt.Println()
	for i := LoALUMode; i <= HiALUMode; i++ {
		// fmt.Printf("%-15s = %032b\n", i, i)
		printConstDebugString(i)
	}
	fmt.Println()
}

func printConstDebugString[T constraints.Unsigned](con T) {
	fmt.Printf(strings.Join([]string{"%32s = %0", strconv.Itoa(reflect.TypeOf(con).Bits()), "b\n"}, ""), con, con)
}
