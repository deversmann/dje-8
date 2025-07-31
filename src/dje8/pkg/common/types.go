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
	CIW                                       // Program Counter In from Address Bus
	CIL                                       // Program Counter LSB In from Data Bus
	CIH                                       // Program Counter MSB In from Data Bus
	COW                                       // Program Counter Out
	CU                                        // Program Counter Increment
	CUW                                       // Program Counter Increment by 2
	MIW                                       // Memory Address Register In From Address Bus (removed MIL & MIH)
	MU                                        // Memory Address Register Increment
	MUW                                       // Memory Address Register Increment by 2
	RI                                        // Write to Memory
	RO                                        // Read from Memory to Data Bus
	ROW                                       // Read 2 bytes from memory to Address Bus
	POW                                       // Stack Pointer Out
	PU                                        // Stack Pointer Increment
	PUW                                       // Stack Pointer Increment by 2
	PD                                        // Stack Pointer Decrement
	PDW                                       // Stack Pointer Decrement by 2
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
	ALUNOP    ALUMode  = iota // 0000 no output
	ALUAND                    // 0001 bitwise and
	ALUOR                     // 0010 bitwise or
	ALUXOR                    // 0011 bitwise xor
	ALUNOT                    // 0100 bitwise inversion
	ALUX2_                    // 0101 reserved
	ALUX1_                    // 0110 reserved
	ALUX0_                    // 0111 reserved
	ALUADD                    // 1000 add
	ALUSUB                    // 1001 subtract
	ALUADC                    // 1010 add considering carry flag
	ALUSBC                    // 1011 subtract considering carry flag
	ALUNEG                    // 1100 arithmetic negation
	ALUINC                    // 1101 increment
	ALUDEC                    // 1110 decrement
	ALUCMP                    // 1111 compare
	LoALUMode = ALUNOP        // *Const to enable traversing list
	HiALUMode = ALUCMP        // *Const to enable traversing list
)

//go:generate stringer -type=OpCode
type OpCode uint8

const (
	NOP OpCode = iota
	STOA
	STOZ
	STOM
	LODI
	LODA
	LODZ
	LODM

	NEG
	ASL
	ASR
	NOT
	LSL
	LSR
	ROL
	ROR

	ADDI
	ADDA
	ADDZ
	ADDM
	SUBI
	SUBA
	SUBZ
	SUBM

	ADCI
	ADCA
	ADCZ
	ADCM
	SBCI
	SBCA
	SBCZ
	SBCM

	ANDI
	ANDA
	ANDZ
	ANDM
	ORI
	ORA
	ORZ
	ORM

	XORI
	XORA
	XORZ
	XORM
	CMPI
	CMPA
	CMPZ
	CMPM

	BEQ
	BNE
	BCS
	BCC
	BMI
	BPL
	BVS
	BVC

	SEI
	JMP
	JMPZ
	JSR
	JSRZ
	RTS
	INT
	RTI

	CLZ
	RSV1
	CLC
	RSV2
	CLN
	RSV3
	CLV
	RSV4

	CLI
	RSV5
	RSV6
	RSV7
	RSV8
	PUSH
	POP
	HALT

	FirstOpCode = NOP
	LastOpCode  = HALT
)

var OpCodeLookup map[string]OpCode

func init() {
	OpCodeLookup = make(map[string]OpCode)
	for i := FirstOpCode; i <= LastOpCode; i++ {
		OpCodeLookup[i.String()] = i
	}
}

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
