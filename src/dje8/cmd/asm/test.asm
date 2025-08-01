start:
        LODI 0          ; testing decimal number
        NOP             ; OUT?
loop:   ADDA data       ; testing label
        BCS done        ; 
        NOP             ; OUT
        STOA var2-1     ; testing minus label offset
        LODA data       ; 
        ADDA data+1     ; testing plus label offset
        BCS done        ; 
        NOP             ; OUT
        STOA data       ; 
        LODA data+1     ; 
        JMP loop        ; 
done:   HALT            ; 

data:
var1:   01              ; testing octal
var2:   0x01            ; testing hex