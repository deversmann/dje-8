LODI 0x01       ; LODI 1
NOP             ; OUT(NOP)
ADDA 0x0022     ; ADDA 34
BCS 0x0020      ; BCS 32??
NOP             ; OUT(NOP)
STOA 0x0021     ; STOA 33
LODA 0x0022     ; LODA 34
ADDA 0x0021     ; ADDA 33
BCS 0x0020      ; BCS 32??
NOP             ; OUT(NOP)
STOA 0x0022     ; STOA 34
LODA 0x0021     ; LODA 33
JMP 0x0003      ; JMP 3
HALT            ; HALT 0
0x01            ; 1
0x01            ; 1