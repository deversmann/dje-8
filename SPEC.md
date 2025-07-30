# Specification for the DJE-8

![Image of the Architecture Diagram for the DJE-8](architecture.drawio.svg)

## Overview
### Goals
* Minicomputer designed to be built with TTL/CMOS 74xx series chips
* 8-bit data width
* 16-bit address bus = 64k addressable memory
* Accessible via a Serial or UART connection
* Rudimentary bootloader and OS
* File system and storage

### Stretch Goals
* Memory banking to extend usable memory
* Integrated video output and keyboard input

## Instruction Set Architecture **DRAFT**
The instruction set will be composed of instructions from the following groups.

### Data Movement
1. Load - Load data into the Accumulator
   - `LODI 0x42` or `LOD #0x42`
   - `LODA 0x4242` or `LOD 0x4242`
   - `LODZ 0x42` or `LOD 0x42`
   - `LODM 0x42` or `LOD (0x42)`
2. Store - Store the contents of the Accumulator into memory
   - `STOA 0x4242` or `STO 0x4242`
   - `STOZ 0x42` or `STO 0x42`
   - `STOM 0x42` or `STO (0x42)`

### Arithmetic
1. Add - Add to the Accumulator and put the results into the Accumulator
   - `ADDI 0x42` or `ADD #0x42`
   - `ADDA 0x4242` or `ADD 0x4242`
   - `ADDZ 0x42` or `ADD 0x42`
   - `ADDM 0x42` or `ADD (0x42)`
2. Subtract - Subtract from the Accumulator and put the results into the Accumulator
   - `SUBI 0x42` or `SUB #0x42`
   - `SUBA 0x4242` or `SUB 0x4242`
   - `SUBZ 0x42` or `SUB 0x42`
   - `SUBM 0x42` or `SUB (0x42)`
3. Add with carry - Add to the Accumulator considering the status of the carry bit and put the results into the Accumulator
   - `ADCI 0x42` or `ADC #0x42`
   - `ADCA 0x4242` or `ADC 0x4242`
   - `ADCZ 0x42` or `ADC 0x42`
   - `ADCM 0x42` or `ADC (0x42)`
4. Subtract with carry - Subtract from the Accumulator considering the status of the carry bit and put the results into the Accumulator
   - `SBCI 0x42` or `SBC #0x42`
   - `SBCA 0x4242` or `SBC 0x4242`
   - `SBCZ 0x42` or `SBC 0x42`
   - `SBCM 0x42` or `SBC (0x42)`
5. Negate - Arithmetically negate the contents of the Accumulator and put the results into the Accumulator
   - `NEG`
6. Compare - Compare with the Accumulator and update the status flags accordingly. The contents of the Accumulator is not changed.
   - `CMPI 0x42` or `CMP #0x42`
   - `CMPA 0x4242` or `CMP 0x4242`
   - `CMPZ 0x42` or `CMP 0x42`
   - `CMPM 0x42` or `CMP (0x42)`
7. Arithmetic Shift
   - `ASL`
   - `ASR`

### Logic
1. And - Logical And with the Accumulator and put the results into the Accumulator
   - `ANDI 0x42` or `AND #0x42`
   - `ANDA 0x4242` or `AND 0x4242`
   - `ANDZ 0x42` or `AND 0x42`
   - `ANDM 0x42` or `AND (0x42)`
2. Or - Logical Or with the Accumulator and put the results into the Accumulator
   - `ORI 0x42` or `OR #0x42`
   - `ORA 0x4242` or `OR 0x4242`
   - `ORZ 0x42` or `OR 0x42`
   - `ORM 0x42` or `OR (0x42)`
3. Not - Logical Not of the Accumulator and put the results into the Accumulator
   - `NOT`
4. Xor - Logical Xor with the Accumulator and put the results into the Accumulator
   - `XORI 0x42` or `XOR #0x42`
   - `XORA 0x4242` or `XOR 0x4242`
   - `XORZ 0x42` or `XOR 0x42`
   - `XORM 0x42` or `XOR (0x42)`
5. Logical Shift
   - `LSL`
   - `LSR`
6. Rotate
   - `ROL`
   - `ROR`

### Flow
1. Branch - operand is a one byte signed offset (-128 to +127) from the address of the byte containing the offset
   - `BCS 0x42` - Branch on Carry flag set
   - `BCC 0x42` - Branch on Carry flag clear
   - `BEQ 0x42` - Branch on Zero flag set (Branch if equal)
   - `BNE 0x42` - Branch on Zero flag clear (Branch if not equal)
   - `BMI 0x42` - Branch on Sign flag set (Branch if negative)
   - `BPL 0x42` - Branch on Sign flag clear (Branch if not negative)
   - `BVS 0x42` - Branch on Overflow flag set
   - `BVC 0x42` - Branch on Overflow flag clear
2. Jump
   - `JMP 0x4242` or `JMP 0x4242`
   - `JMPZ 0x42` or `JMP 0x42`
3. Jump to Sub Routine
   - `JSR 0x4242` or `JSR 0x4242`
   - `JSRZ 0x42` or `JSR 0x42`
4. Return from Sub Routine
   - `RTS`
5. Interrupt
   - `INT`
6. Return from Interrupt
   - `RTI`
7. Set flags
   - `SEI`
8. Clear flags
   - `CLI`
   - `CLZ`
   - `CLC`
   - `CLN`
   - `CLV`
9. Push Accumulator to Stack
    - `PUSH`
10. Pop from Stack into Accumulator
    - `POP`
11. Halt
    - `HALT`
12. No Operation
    - `NOP` - Waits the maximum number of clock cycles a single instruction can take (currently 16)

### Addressing Modes
1. **Immediate `I`** - Argument is the value of the operand
2. **Absolute `A`** - Argument is the 16-bit memory address of the operand
3. **ZeroPage `Z`** - Argument is the 8-bit address (referencing the ZeroPage (`0x0000-0x00ff`)) of the operand
4. **Memory Indirect `M`**- Argument is the 8-bit memory address (referencing the ZeroPage (`0x0000-0x00ff`)) of a location containing the address of the location of the operand

> [!NOTE]
> Other Modes considered but not implemented at this time:
> 1. **Indexed** - Same as Absolute except that the address is offset by the contents of the Accumulator (the address wraps around at the min and max addresses)
> 2. **Indexed ZeroPage** - Same as ZeroPage except that the address is offset by the value of the Accumulator (the address wraps around and remains on the ZeroPage)
> 3. **Indexed Indirect** - the argument is offset by the contents of the Accumulator (ZeroPage wrap around applies) to acquire the Indirect address
> 4. **Indirect Indexed** - the argument contains an address which is offset by the contents of the accumulator to find the address containing the operand



## Registers
- **Accumulator Register** (A) - 8-bit
- **Internal Register** (B) - 8-bit
- **Instruction Register** (IR) - 8-bit
- **Program Counter** (PC) - 16-but
- **Memory Address Register** (MAR) - 16-bit
- **Stack Pointer** (SP) - 8-bit
- **Flags Register** - 8-bit
    ```
    - - - I Z C N V
          | | | | |
          | | | | +-- Overflow Flag 
          | | | +---- Sign Flag 
          | | +------ Carry Flag 
          | +-------- Zero Flag 
          +---------- Interrupt Disable Flag 
    ```

## Control Logic Signals
The control signals make up the core of the logic of the processor. Combined together, they form the microcode that makes up each one of the processor instructions. There are 28 separate control signals and they are stored in Control ROM as 32-bit Control Words.

| Signal | Name	| Purpose |
|---|---|---|
| `0x80000000` | `HLT` | HALT |
| `0x40000000` | `AI` | Accumulator Register In |
| `0x20000000` | `AO` | Accumulator Register Out |
| `0x10000000` | `BI` | Internal Register In |
| `0x08000000` | `II` | Instruction Register In |
| `0x04000000` | `CIW` | Program Counter In from Address Bus |
| `0x02000000` | `CIL` | Program Counter LSB In from Data Bus |
| `0x01000000` | `CIH` | Program Counter MSB In from Data Bus |
| `0x00800000` | `COW` | Program Counter Out |
| `0x00400000` | `CU` | Program Counter Increment |
| `0x00200000` | `CUW` | Program Counter Increment by 2 |
| `0x00100000` | `MIW` | Memory Address Register In From Address Bus (removed MIL & MIH) |
| `0x00080000` | `MU` | Memory Address Register Increment |
| `0x00040000` | `MUW` | Memory Address Register Increment by 2 |
| `0x00020000` | `RI` | Write to Memory |
| `0x00010000` | `RO` | Read from Memory to Data Bus |
| `0x00008000` | `ROW` | Read 2 bytes from memory to Address Bus |
| `0x00004000` | `POW` | Stack Pointer Out |
| `0x00002000` | `PU` | Stack Pointer Increment |
| `0x00001000` | `PUW` | Stack Pointer Increment |
| `0x00000800` | `PD` | Stack Pointer Decrement |
| `0x00000400` | `PDW` | Stack Pointer Decrement |
| `0x00000200` | `AU3` | ALU Mode Bit 3 |
| `0x00000100` | `AU2` | ALU Mode Bit 2 |
| `0x00000080` | `AU1` | ALU Mode Bit 1 |
| `0x00000040` | `AU0` | ALU Mode Bit 0 |
| `0x00000020` | `FL` | Flags Register In |
| `0x00000010` | `EX3` | Reserved |
| `0x00000008` | `EX2` | Reserved |
| `0x00000004` | `EX1` | Reserved |
| `0x00000002` | `EX0` | Reserved |
| `0x00000001` | `STR` | Step Counter Reset |

## Memory Map


### I/O Memory Map **DRAFT**
| Address Range | Size | Purpose |
|---|---|---|
| `0xF000 - 0xF00F` | 16 Bytes | Serial Port 1 (UART) |
| `0xF010 - 0xF01F` | 16 Bytes | Serial Port 2 (UART) |
| `0xF020 - 0xF02F` | 16 Bytes | Keyboard Interface |
| `0xF030 - 0xF3FF` | ~1 KB | Reserved for other built-in devices |
| `0xF400 - 0xF7FF` | 1 KB | Expansion Slot 1 |
| `0xF800 - 0xFBFF` | 1 KB | Expansion Slot 2 |
| `0xFC00 - 0xFFFD` | 1 KB | Video Controller Registers / Buffer |
| `0xFFFE - 0xFFFF` | 2 bytes | Interrupt Vector |

`0xB000 - 0xBFFF`  4 KB   Video Character Buffer (80x25 @ 8-bit color)

-----

## Notes

### DMA
Need Bus-Arbiter for DMA... x number of prioritized DMA peripheral slots
- CPU needs BUS_REQ and BUS_ACK signal pins
- DMA Devices need address and data buses and DMA_REQ and DMA_GNT lines

### Text Mode Considerations
1. Text mode video RAM: 80x25 characters = 2000 character positions on screen
2. 1 byte for ASCII code + 1 byte for color = 2 bytes per character
3. Total 4KB for 80x25 character mode at 256 colors (w/96 bytes to spare)
4. 640x480 px and 80x25 char = 8x19 font (including perimeter space)







