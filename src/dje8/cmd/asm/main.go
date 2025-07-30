package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"unicode"

	//lint:ignore ST1001 importing common shared across all dje8 cmds
	. "damien.live/dje8/pkg/common"
)

type ByteValue byte

func (v *ByteValue) String() string {
	return strconv.FormatUint(uint64(*v), 16)
}

func (v *ByteValue) Set(s string) error {
	if temp, err := strconv.ParseUint(s, 0, 8); err != nil {
		return err
	} else {
		*v = ByteValue(temp)
	}
	return nil
}

type ModeValue rune

func (v *ModeValue) String() string {
	return string(*v)
}

func (v *ModeValue) Set(s string) error {
	if strings.HasPrefix(s, "X") || strings.HasPrefix(s, "x") {
		*v = 'x'
	} else if strings.HasPrefix(s, "B") || strings.HasPrefix(s, "b") {
		*v = 'b'
	} else if strings.HasPrefix(s, "A") || strings.HasPrefix(s, "a") {
		*v = 'a'
	} else {
		return fmt.Errorf("cannot process %s into mode", s)
	}
	return nil
}

var filename string
var paddedSize int
var paddingByte ByteValue
var mode ModeValue = 'x'

func init() {
	const (
		modeUsage        = "a output bytes to the console in roughly the same organization as the source asm\nx to output bytes to the console in a format similar to hexdump\nb to output bytes in a binary file"
		filenameUsage    = "the name of the file containing the code to be assembled"
		paddedSizeUsage  = "size in bytes the binary file will be padded to. will not be padded if the size is smaller than the number of bytes generated"
		paddingByteUsage = "byte to use as padding"
	)
	flag.Var(&mode, "m", modeUsage)
	flag.Var(&paddingByte, "p", paddingByteUsage)
	flag.IntVar(&paddedSize, "s", 0, paddedSizeUsage)
	flag.StringVar(&filename, "filename", "", filenameUsage)
	flag.StringVar(&filename, "f", "", filenameUsage)
}

func main() {
	flag.Parse()
	if strings.TrimSpace(filename) == "" {
		fmt.Fprintf(os.Stderr, "Error: filename required\n") // Print error message to standard error
		os.Exit(1)                                           // Exit with a non-zero status code (e.g., 1 for general error)
	}
	fmt.Printf("mode = %c; paddedSize = %d; paddingByte = %d; filename = %s\n", mode, paddedSize, paddingByte, filename)

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

	switch mode {
	case 'a':
		// Iterate and print lines
		for _, line := range lines {
			fmt.Println(line)
		}
	case 'x':
		var bytes []byte
		for _, line := range lines {
			for s := range strings.SplitSeq(line, " ") {
				b, _ := strconv.ParseInt(s, 0, 8)
				bytes = append(bytes, byte(b))
			}
		}
		var chars []byte
		for i, b := range bytes {
			chars = append(chars, b)
			if i%16 == 0 {
				fmt.Printf("0x%04x:", i)
			}
			if i%8 == 0 {
				fmt.Print("  ")
			}
			fmt.Printf(" %02x", b)
			if (i+1)%16 == 0 {
				fmt.Print("\033[65G|")
				for _, c := range chars {
					if unicode.IsPrint(rune(c)) {
						fmt.Print(string(c))
					} else {
						fmt.Print(".")
					}
				}
				fmt.Println("|")
				chars = []byte{}
			}
		}
		if len(chars) > 0 {
			fmt.Print("\033[65G|")
			for _, c := range chars {
				if unicode.IsPrint(rune(c)) {
					fmt.Print(string(c))
				} else {
					fmt.Print(".")
				}
			}

		}

		fmt.Println("|")
		fmt.Println()

	case 'b':
		var bytes []byte
		for _, line := range lines {
			for s := range strings.SplitSeq(line, " ") {
				b, _ := strconv.ParseInt(s, 0, 8)
				bytes = append(bytes, byte(b))
			}
		}
		if len(bytes) < paddedSize {
			for i := len(bytes); i < paddedSize; i++ {
				bytes = append(bytes, byte(paddingByte))
			}
		}
		os.WriteFile(strings.Join([]string{filename, ".bin"}, ""), bytes, fs.ModePerm)

	}
}
