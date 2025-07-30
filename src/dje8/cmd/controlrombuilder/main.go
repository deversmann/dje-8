package main

import (
	"fmt"

	//lint:ignore ST1001 importing common shared across all dje8 cmds
	. "damien.live/dje8/pkg/common"
	"damien.live/dje8/pkg/ucodebuilder"
)

func main() {
	fmt.Println("***** DJE-8 ControlROM Builder *****")
	fmt.Println()
	DebugPrintConsts()

	Ucode := ucodebuilder.BuildUcode()

	for i, ControlWord := range Ucode {
		if i%16 == 0 {
			fmt.Printf("\n0x%04x:", i)
		}
		if ControlWord == 0 {
			fmt.Print(" 0x0")
		} else {
			fmt.Printf(" 0x%08x", uint32(ControlWord))
		}
	}

	fmt.Println()
}
