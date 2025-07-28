package ucodebuilder

import (
	//lint:ignore ST1001 importing common shared across all dje8 cmds
	. "damien.live/dje8/pkg/common"
)

func BuildUcode() []Control {

	ControlROM := [16][16]Control{
		{CO2 | MI2, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                                                   // NOP
		{CO2 | MI2, RO | II | CU, CO2 | MI2, RO2 | MI2 | CU2, RO | AI | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                   // LDA
		{CO2 | MI2, RO | II | CU, CO2 | MI2, RO2 | MI2 | CU2, RO | BI, AU3 | AI | FL | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},       // ADD
		{CO2 | MI2, RO | II | CU, CO2 | MI2, RO2 | MI2 | CU2, RO | BI, AU3 | AU0 | AI | FL | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // SUB
		{CO2 | MI2, RO | II | CU, CO2 | MI2, RO2 | MI2 | CU2, RI | AO | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                   // STA
		{CO2 | MI2, RO | II | CU, CO2 | MI2, RO | AI | CU | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                            // LDI
		{CO2 | MI2, RO | II | CU, CO2 | MI2, RO2 | CI2 | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                               // JMP
		{CO2 | MI2, RO | II | CU, CU2 | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                                             // JC
		{CO2 | MI2, RO | II | CU, CU2 | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                                             // JZ
		{CO2 | MI2, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                                                   // ---
		{CO2 | MI2, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                                                   // ---
		{CO2 | MI2, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                                                   // ---
		{CO2 | MI2, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                                                   // ---
		{CO2 | MI2, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                                                   // ---
		{CO2 | MI2, RO | II | CU, STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                                                   // OUT
		{CO2 | MI2, RO | II | CU, HLT | STR, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                                             // HLT
	}

	var Ucode [65536]Control

	for i := range 16 {
		for j, Instr := range ControlROM {
			if i&int(CarryFlagC) != 0 && j == 7 {
				Instr = ControlROM[6]
			}
			if i&int(ZeroFlagZ) != 0 && j == 8 {
				Instr = ControlROM[6]
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
