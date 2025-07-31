package main

// TODO :
//   1. label pointers (the second half of labels)
//	 2. label arithmetic i.e. `JMP label+2`
//   3. strings in data i.e. `stringvar: 'Hello, World!',0`
//   4. operand checking

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"damien.live/dje8/pkg/common"
)

type AsmToken struct {
	Pointer string
	Val     byte
}

var labelMap = make(map[string]uint16)
var tokens []AsmToken
var currentByte uint16

var filename string
var paddedSize int = 0
var paddingByte ByteValue = 0x00
var mode ModeValue = 'x'

func main() {
	// parse args... only requirement is filename
	flag.Parse()
	if strings.TrimSpace(filename) == "" {
		fmt.Fprintf(os.Stderr, "Error: filename required\n")
		os.Exit(1)
	}

	// Read the file into a byte slice
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	// Convert the byte slice to a string and split by lines and tokenize
	err = TokenizeLines(strings.Split(string(data), "\n"))
	if err != nil {
		fmt.Printf("problem parsing file: %e\n", err)
		os.Exit(1)
	}

	// debug out
	for k, v := range labelMap {
		fmt.Printf("label %s = 0x%04x\n", k, v)
	}

	if mode == 'x' {
		printHexDump()
	}
	if mode == 'b' {
		writeBinFile()
	}

}

func TokenizeLines(lines []string) error {
	var fields []string
	var numOperandBytesRemaining int
	for _, line := range lines {
		fields = append(fields, strings.Fields(strings.Split(line, ";")[0])...)
	}
	token := AsmToken{}
	for _, field := range fields {
		if isLabel(field) {
			labelMap[strings.Trim(field, ":")] = uint16(len(tokens))
		} else if isOpCode(field) {
			// TODO hone error checking of operand length
			if numOperandBytesRemaining > 0 {
				return fmt.Errorf("not enough bytes before next opcode (%s) at byte %d", field, currentByte)
			}
			token.Val = byte(common.OpCodeLookup[field])
			numOperandBytesRemaining = numberOfOperandBytes(common.OpCodeLookup[field])
			token = nextToken(token)
		} else if isRawData(field) {
			// if the field is a byte or word:
			//    update bytesremaining if necessary
			//    store the byte or bytes separately creating new tokens
			isTwoBytes := false
			if strings.HasPrefix(field, "0x") && (len(field) > 4) {
				isTwoBytes = true
			}
			data, err := strconv.ParseUint(field, 0, 16)
			if err != nil {
				return fmt.Errorf("error parsing data byte %d = '%s' : %e", currentByte, field, err)
			}
			if data > 255 {
				isTwoBytes = true
			}
			if !isTwoBytes {
				numOperandBytesRemaining--
				token.Val = byte(data)
				token = nextToken(token)
			} else {
				numOperandBytesRemaining -= 2
				token.Val = byte(data >> 8)
				token = nextToken(token)
				token.Val = byte(data & 0xff)
				token = nextToken(token)
			}
			if numOperandBytesRemaining < 0 {
				numOperandBytesRemaining = 0
			}
		} else {
			// default: the field is a pointer to a label
			//   store it in the token for the second pass
			token.Pointer = field
			token = nextToken(token)
		}
	}
	return nil
}

func nextToken(token AsmToken) AsmToken {
	tokens = append(tokens, token)
	currentByte++
	return AsmToken{}
}

func isLabel(token string) bool {
	return strings.HasSuffix(token, ":")
}

func isOpCode(token string) bool {
	_, valid := common.OpCodeLookup[token]
	return valid
}

// Raw data is decimal, hex prefixed with 0x, or a single-quoted character
// TODO - if the character resolves to multi-byte, we could have a problem
// TODO - need an easy way to store strings of characters (0 terminated)
func isRawData(token string) bool {
	result, err := regexp.MatchString("^([0-9]+)|(0x[0-9A-Fa-f]+)|('.')$", token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "You broke the REGEXP again, apparently: %e\n", err)
		os.Exit(1)
	}
	return result
}

func numberOfOperandBytes(op common.OpCode) int {
	// TODO
	return 0
}

func writeBinFile() {
	var bytes []byte
	for _, token := range tokens {
		bytes = append(bytes, token.Val)
	}
	if len(bytes) < paddedSize {
		for i := len(bytes); i < paddedSize; i++ {
			bytes = append(bytes, byte(paddingByte))
		}
	}
	os.WriteFile(strings.Join([]string{filename, ".bin"}, ""), bytes, fs.ModePerm)
}

func printHexDump() {
	var printedBytes []byte
	for i, token := range tokens {
		if i%16 == 0 {
			fmt.Printf("%08x  ", i)
		}
		fmt.Printf("%02x ", token.Val)
		if unicode.IsPrint(rune(token.Val)) {
			printedBytes = append(printedBytes, token.Val)
		} else {
			printedBytes = append(printedBytes, '.')
		}
		if (i+1)%8 == 0 {
			fmt.Print(" ")
		}
		if (i+1)%16 == 0 || i+1 == len(tokens) {
			fmt.Printf("\033[61G|%s|\n", string(printedBytes))
			printedBytes = []byte{}
		}
	}
	fmt.Printf("%08x\n", len(tokens))
}

// *** CLI FLag Stuff ***
type ByteValue byte
type ModeValue rune

func init() {
	const (
		modeUsage = "x - output bytes to the console in a format similar to hexdump\n" +
			"b - output bytes in a binary file"
		paddingByteUsage = "byte to use as padding if outputting binary file"
		paddedSizeUsage  = "size in bytes to pad if outputting binary file\n" +
			"will not be padded if the size is smaller than the number of bytes generated"
		filenameUsage = "required: the name of the file containing the code to be assembled"
	)
	flag.Var(&mode, "m", modeUsage)
	flag.Var(&paddingByte, "p", paddingByteUsage)
	flag.IntVar(&paddedSize, "s", 0, paddedSizeUsage)
	flag.StringVar(&filename, "f", "", filenameUsage)
}

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
