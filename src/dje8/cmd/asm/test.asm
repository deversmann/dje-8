start:
        LODI 0x01       ; LODI 1
loop:   NOP             ; OUT(NOP)
        ADDA 0x0022     ; ADDA 34   (var2)
        BCS 0x0020      ; BCS 32??  (done)
        NOP             ; OUT(NOP)
        STOA 0x0021     ; STOA 33   (var1)
        LODA 0x0022     ; LODA 34   (var2)
        ADDA 0x0021     ; ADDA 33   (var1)
        BCS 0x0020      ; BCS 32??  (done)
        NOP             ; OUT(NOP)
        STOA 0x0022     ; STOA 34   (var2)
        LODA 0x0021     ; LODA 33   (var1)
        JMP 0x0003      ; JMP 3     (loop)
done:   HALT            ; HALT

data:
var1:   0x01            ; 1
var2:   0x01            ; 1