package main

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

type asmToken struct {
	value   byte
	pointer string
	lineNo  uint16
}

type field struct {
	content string
	lineNo  uint16
}

var tokens []asmToken
var fields []field
var currentAddress uint16 = 0
var labelMap map[string]uint16 = make(map[string]uint16)

var filename string
var paddedSize int = 0
var paddingByte ByteValue = 0x00
var mode ModeValue = 'x'

func main() {
	flag.Parse() // parse args... only requirement is filename
	if strings.TrimSpace(filename) == "" {
		die("Error: filename required\n")
	}

	fileBytes, err := os.ReadFile(filename) // read file
	if err != nil {
		die(fmt.Sprintf("Problem reading file: %e\n", err))
	}
	for i, lineStr := range strings.Split(string(fileBytes), "\n") { // parse into fields
		for fieldStr := range strings.FieldsSeq(strings.ReplaceAll(strings.Split(lineStr, ";")[0], ",", " ")) {
			fields = append(fields, field{fieldStr, uint16(i + 1)})
		}
	}
	for i := 0; i < len(fields); i++ { // first pass process each field except label pointers; map label addresses
		currentField := fields[i]
		if strings.HasPrefix(currentField.content, "#") { // DIRECTIVE
			if currentField.content == "#org" {
				isStartingDirective := (i == 0)
				i++
				currentField = fields[i]
				address, err := strconv.ParseUint(currentField.content, 0, 16)
				if err != nil {
					die(fmt.Sprintf("error parsing #org directive address (%s) on line %d: %e\n", currentField.content, currentField.lineNo, err))
				}
				if !isStartingDirective {
					if uint16(address) < currentAddress {
						die(fmt.Sprintf("invalid #org directive on line %d, would result in a negative offset\n", currentField.lineNo))
					}
					tokens = append(tokens, make([]asmToken, uint16(address)-currentAddress)...) // padding
				}
				currentAddress = uint16(address)
			} else {
				die(fmt.Sprintf("unknown assembler directive (%s) on line %d\n", currentField.content, currentField.lineNo))
			}
		} else if strings.HasSuffix(currentField.content, ":") { // LABEL
			labelMap[currentField.content[:len(currentField.content)-1]] = currentAddress
		} else if _, found := common.OpCodeLookup[currentField.content]; found { // INSTRUCTION
			tokens = append(tokens, asmToken{byte(common.OpCodeLookup[currentField.content]), "", currentField.lineNo})
			currentAddress++
		} else if isData(currentField.content) {
			if strings.HasPrefix(currentField.content, "'") {
				tokens = append(tokens, asmToken{[]byte(currentField.content)[1], "", currentField.lineNo})
				currentAddress++
			} else {
				isTwoBytes := (strings.HasPrefix(currentField.content, "0x") && len(currentField.content) > 4)
				parsed, err := strconv.ParseUint(currentField.content, 0, 16)
				if err != nil {
					die(fmt.Sprintf("error parsing data (%s) at line %d: %e\n", currentField.content, currentField.lineNo, err))
				}
				isTwoBytes = isTwoBytes || parsed > 255
				tokens = append(tokens, asmToken{byte(parsed), "", currentField.lineNo}) // one byte or LSB
				currentAddress++
				if isTwoBytes {
					tokens = append(tokens, asmToken{byte(parsed >> 8), "", currentField.lineNo}) // MSB
					currentAddress++
				}
			}
		} else { // must be pointer
			if strings.ContainsAny(currentField.content, "<>") {
				tokens = append(tokens, asmToken{0, currentField.content, currentField.lineNo})
				currentAddress++
			} else {
				tokens = append(tokens, []asmToken{{0, "<" + currentField.content, currentField.lineNo}, {0, ">" + currentField.content, currentField.lineNo}}...)
				currentAddress += 2
			}
		}
	}

	for i := range len(tokens) { // second pass, replace labell pointers with addresses of labels
		pointer := tokens[i].pointer
		if pointer != "" {
			offsetIdx := strings.LastIndexAny(pointer, "+-")
			isMSB := strings.HasPrefix(pointer, ">")
			var offset int64
			var err error
			if offsetIdx > 0 {
				offset, err = strconv.ParseInt(pointer[offsetIdx:], 0, 16)
				if err != nil {
					die(fmt.Sprintf("error parsing pointer (%s) at line %d: %e\n", tokens[i].pointer, tokens[i].lineNo, err))
				}
				pointer = pointer[:offsetIdx]
			}
			address, found := labelMap[pointer[1:]]
			if !found {
				die(fmt.Sprintf("error parsing pointer (%s) at line %d: %e\n", pointer[1:], tokens[i].lineNo, err))
			}
			address += uint16(offset)
			if isMSB {
				tokens[i].value = byte(address >> 8)
			} else {
				tokens[i].value = byte(address)
			}
		}
	}

	if mode == 'x' {
		chars := ""
		for i, token := range tokens {
			if i%16 == 0 {
				fmt.Printf("%08x ", i)
			}
			fmt.Printf(" %02x", token.value)
			if unicode.IsPrint(rune(token.value)) {
				chars = chars + string(token.value)
			} else {
				chars = chars + "."
			}
			if (i+1)%8 == 0 {
				fmt.Print(" ")
			}
			if (i+1)%16 == 0 || i == len(tokens)-1 {
				fmt.Printf(" \033[61G|%s|\n", chars)
				chars = ""
			}
		}
		fmt.Printf("%08x\n", len(tokens))
	}

	if mode == 'b' {
		var bytes []byte
		for _, token := range tokens {
			bytes = append(bytes, token.value)
		}
		if len(bytes) < paddedSize {
			for i := len(bytes); i < paddedSize; i++ {
				bytes = append(bytes, byte(paddingByte))
			}
		}
		os.WriteFile(strings.Join([]string{filename, ".bin"}, ""), bytes, fs.ModePerm)
	}
}

func isData(field string) bool {
	matched, err := regexp.MatchString("^(('.')|([+-]?(0|[1-9][0-9]*))|(0[0-7]*)|(0x[0-9a-fA-F]*))$", field)
	if err != nil {
		die(fmt.Sprintf("regular expression to match numerics is in error: %e", err))
	}
	return matched
}

func die(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
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
