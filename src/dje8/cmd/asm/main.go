package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	//lint:ignore ST1001 importing common shared across all dje8 cmds
	. "damien.live/dje8/pkg/common"
)

var filename string

func init() {
	const (
		usage = "the name of the file containing the code to be assembled"
	)
	flag.StringVar(&filename, "filename", "", usage)
	flag.StringVar(&filename, "f", "", usage)
}

func main() {
	flag.Parse()
	if strings.TrimSpace(filename) == "" {
		fmt.Fprintf(os.Stderr, "Error: filename required\n") // Print error message to standard error
		os.Exit(1)                                           // Exit with a non-zero status code (e.g., 1 for general error)
	}

	// Read the file into a byte slice
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	// Convert the byte slice to a string and split by lines
	lines := strings.Split(string(data), "\n")

	// Iterate over lines and do the following:
	// - Remove comments and extra whitespace
	// - Replace Mnemonics with OpCodes
	// - Break all words into bytes
	// - Put labels in a table
	// - Deal with directives

	// Iterate and print lines
	for i, line := range lines {
		fields := strings.Fields((strings.SplitN(line, ";", 2))[0])

		//OpCodeLookup[fields[0]]
		for j, s := range fields {
			token8, err := strconv.ParseUint(s, 0, 8)
			if err != nil {
				token16, err := strconv.ParseUint(s, 0, 16)
				if err != nil {
					tokenOp, ok := OpCodeLookup[s]
					if !ok {
						fmt.Fprintf(os.Stderr, "Error parsing %s from line %d", s, i)
						os.Exit(1)
					}
					fields[j] = fmt.Sprintf("0x%02x", uint8(tokenOp))
				} else {
					fields[j] = fmt.Sprintf("0x%02x 0x%02x", token16>>8, token16&0x00ff)
				}
			} else {
				if len(s) > 4 {
					fields[j] = fmt.Sprintf("0x%02x 0x%02x", token8>>8, token8&0x00ff)
				} else {
					fields[j] = fmt.Sprintf("0x%02x", token8)
				}
			}
		}
		lines[i] = strings.Join(fields, " ")
	}

	// Iterate and print lines
	for _, line := range lines {
		fmt.Println(line)
	}

}
