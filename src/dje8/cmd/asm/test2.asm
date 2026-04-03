a=startidx            LODM start
            initloop:
            STOZ 1
            INC
            bcc initloop
            LDI 0
            STOZ start
            STOZ start+1
            LDI 1
            nextprime:
            inc a
            bcs end
            STOZ current
            LODM current
            CMPI 0
            BEQ nextprime
            LODA current
            sieveloop:
            ADDA current
            BCS cleanup
            STOZ current
            LODI 0
            STOM current
            LODA current
            JMP sieveloop
            cleanup:
            LODA current
            JMP nextprime
            end:
            HALT
            current: 0
            start: 0
            #org $+255
            max: 0


    




data:
current: 1
start: 0
#org $+255
end: 0