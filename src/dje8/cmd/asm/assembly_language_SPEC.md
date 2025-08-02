# Assembly for the DJE-8 Language Specification

1. All language tokens must be separated by whitespace (space, tab, newline).
2. Comments run from the first instance of a semicolon (`;`) in a line to the end of the line and are ignored by the assembler
3. All language tokens are single strings of characters with no whitespace embedded
4. Language tokens fall into the following categories:
   1. Instruction
   2. Label
   3. Operand, including label references
   4. Label reference
   5. Assembler directive
5. Instructions appear in all capital characters and may not be abbreviated
   1. Instruction are reserved words and may not be used as labels
6. Labels end in a colon (`:`)
7. Labels appear immediately before (separated only by whitespace) the token containing the byte they reference
8. Labels may be made up of the following characters: `[A-Za-z0-9_]`
9. Operands are limited to one or two bytes and may be one of the following:
   1.  An integer literal in decimal made up only of numeric characters with no leading zeros (`42`)
   2.  An integer literal in octal made up of only numeric characters with a leading zero (`042`)
   3.  An integer literal in hexadecimal made up of numeric characters and the letters A,B,C,D,E,F in upper or lower case preceeded by a zero and a lower case 'x' (`0x42de`)
   4.  A single character literal surrounded by single quotes (`'Z'`)
   5.  A reference to a labeled memory location which defaults to two bytes. 
       1.  A reference takes 1 bytes if it is prepended with a less than (`<`) or greater than (`>`) symbol.  These indicate only loading the low byte or high byte of the address respectively.
       2.  Label references may also have a plus (`+`) or minus (`-`) appended with a number signifying an number of offset bytes after or before the memory location respectively. The symbol and the number are appended without spaces (`label+4`) and the offset is applied before applying a byte selection.
10. `#org` followed an address is a directive that sets the memory location of the following bytes. 
    1.  If the use of the directive creates a gap after the preceeding code, that space will be filled with 0x00. 
    2.  `#org` cannot move backward.  
    3.  If the directive is the first non-comment, non-whitespace token, it is understood to be the starting point of the assembly and all following bytes and addresses will be numbered from that point. 
    4.  `#org` and its address must appear on their own line with the exception of a label that immediately follows an `#org` directive.  This label will refer to the memory location named by the directive. 

### TODOs
- [x] Implement octal and character literals
- [ ] Implementation of directives, e.g. `#org`
- [ ] Definition and implementation of strings in ASM, e.g. `'this is a zero terminated string', 0`
- [ ] Implement operand length checking, i.e. are enough bytes included after each instruction?
