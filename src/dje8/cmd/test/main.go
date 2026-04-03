package main

import (
	"fmt"
	"os"
)

var max = 256

func main() {

	sieve := make([]int, max)
	for i := range max {
		sieve[i] = 1
	}
	sieve[0] = 0
	sieve[1] = 0
	current := 1
	for {
		for sieve[current] != 1 {
			current++
			if current == max {
				printAndExit(sieve)
			}
		}
		for i := 2 * current; i < max; i += current {
			sieve[i] = 0
		}
		current++
	}

}

func printAndExit(sieve []int) {
	for i := range len(sieve) {
		if sieve[i] == 1 {
			fmt.Printf("%3d ", i)
		} else {
			fmt.Print("--- ")
		}
		if (i+1)%16 == 0 {
			fmt.Println()
		}
	}
	os.Exit(0)
}
