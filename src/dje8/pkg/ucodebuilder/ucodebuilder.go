package ucodebuilder

import (
	//lint:ignore ST1001 importing common shared across all dje8 cmds
	. "damien.live/dje8/pkg/common"
)

func BuildUcode() []Control {

	ControlROM := [80][16]Control{
		/* NOP */ {COW | MIW, RO | II | CU, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* STOA */ {COW | MIW, RO | II | CU, COW | MIW, ROW | MIW | CUW, RI | AO | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* STOZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* STOM */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* LODI */ {COW | MIW, RO | II | CU, COW | MIW, RO | AI | CU | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* LODA */ {COW | MIW, RO | II | CU, COW | MIW, ROW | MIW | CUW, RO | AI | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* LODZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* LODM */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},

		/* NEG */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ASL */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ASR */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* NOT */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* LSL */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* LSR */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ROL */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ROR */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},

		/* ADDI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ADDA */ {COW | MIW, RO | II | CU, COW | MIW, ROW | MIW | CUW, RO | BI, AU3 | AI | FL | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ADDZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ADDM */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* SUBI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* SUBA */ {COW | MIW, RO | II | CU, COW | MIW, ROW | MIW | CUW, RO | BI, AU3 | AU0 | AI | FL | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* SUBZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* SUBM */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},

		/* ADCI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ADCA */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ADCZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ADCM */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* SBCI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* SBCA */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* SBCZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* SBCM */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},

		/* ANDI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ANDA */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ANDZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ANDM */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ORI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ORA */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ORZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* ORM */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},

		/* XORI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* XORA */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* XORZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* XORM */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* CMPI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* CMPA */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* CMPZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* CMPM */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},

		/* BEQ */ {COW | MIW, RO | II | CU, CUW | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // TODO currently absolute, make 8-bit relative
		/* BNE */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* BCS */ {COW | MIW, RO | II | CU, CUW | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // TODO currently absolute, make 8-bit relative
		/* BCC */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* BMI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* BPL */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* BVS */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* BVC */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},

		/* SEI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* JMP */ {COW | MIW, RO | II | CU, COW | MIW, ROW | CIW | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* JMPZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* JSR */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* JSRZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* RTS */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* INT */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* RTI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},

		/* CLZ */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* RSV1 */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* CLC */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* RSV2 */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* CLN */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* RSV3 */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* CLV */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* RSV4 */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},

		/* CLI */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* RSV5 */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* RSV6 */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* RSV7 */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* RSV8 */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* PUSH */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* POP */ {COW | MIW, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		/* HALT */ {COW | MIW, RO | II | CU, HLT | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	var Ucode [65536]Control

	for i := range 16 {
		for j, Instr := range ControlROM {
			if i&int(CarryFlagC) != 0 && j == int(BCS) {
				Instr = ControlROM[JMP]
			}
			if i&int(ZeroFlagZ) != 0 && j == int(BEQ) {
				Instr = ControlROM[JMP]
			}
			idx := i<<12 + j*16
			copy(Ucode[idx:], Instr[:])
		}
	}

	// var idx uint32 = 0
	// for i := range 256 {
	// 	for j, Instr := range ControlROM {
	// 		if j == 7 && i/16&int(CarryFlagC) != 0 {
	// 			Instr = ControlROM[6]
	// 		}
	// 		if j == 8 && i/16&int(ZeroFlagZ) != 0 {
	// 			Instr = ControlROM[6]
	// 		}
	// 		for _, Word := range Instr {
	// 			Ucode[idx] = Word
	// 			idx++
	// 		}
	// 	}
	// }

	return Ucode[:]
}
