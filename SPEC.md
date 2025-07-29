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
   - `LDI #0x42`
   - `LDA 0x4242`
   - `LDZ 0x42`
   - `LDM (0x42)`
2. Store - Store the contents of the Accumulator into memory
   - `STA 0x4242`
   - `STZ 0x42`
   - `STM (0x42)`

### Arithmetic
1. Add - Add to the Accumulator and put the results into the Accumulator
   - `ADI`
   - `ADA`
   - `ADZ`
   - `ADM`
2. Subtract - Subtract from the Accumulator and put the results into the Accumulator
   - `SBI`
   - `SBA`
   - `SBZ`
   - `SBM`
3. Add with carry - Add to the Accumulator considering the status of the carry bit and put the results into the Accumulator
   - `ACI`
   - `ACA`
   - `ACZ`
   - `ACM`
4. Subtract with carry - Subtract from the Accumulator considering the status of the carry bit and put the results into the Accumulator
   - `SCI`
   - `SCA`
   - `SCZ`
   - `SCM`
5. Negate - Arithmetically negate the contents of the Accumulator and put the results into the Accumulator
   - `NEG`
6. Compare - Compare with the Accumulator and update the status flags accordingly.  The contents of the Accumulator is not changed.
   - `CPI`
   - `CPA`
   - `CPZ`
   - `CPM`
7. Arithmetic Shift
   - `ASL`
   - `ASR`

### Logic
1. And - Logical And with the Accumulator and put the results into the Accumulator
   - `ANI`
   - `ANA`
   - `ANZ`
   - `ANM`
2. Or - Logical Or with the Accumulator and put the results into the Accumulator
   - `ORI`
   - `ORA`
   - `ORZ`
   - `ORM`
3. Not - Logical Not of the Accumulator and put the results into the Accumulator
   - `NOT`
4. (Xor) - Logical Xor with the Accumulator and put the results into the Accumulator
5. Logical Shift
   - `LSL`
   - `LSR`
6. Rotate
   - `ROL`
   - `ROR`

### Flow
1. Branch - operand is a one byte (-128 to +127) offset from the current address
   - `BCS` - Branch on Carry flag set
   - `BCC` - Branch on Carry flag clear
   - `BEQ` - Branch on Zero flag set (Branch if equal)
   - `BNE` - Branch on Zero flag clear (Branch if not equal)
   - `BMI` - Branch on Sign flag set (Branch if negative)
   - `BPL` - Branch on Sign flag clear (Branch if not negative)
   - `BVS` - Branch on Overflow flag set
   - `BVC` - Branch on Overflow flag clear
2. Jump
   - `JMP`
3. Jump to Sub Routine
   - `JSR`
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
9.  Push to Stack
    - `PSH`
10. Pop from Stack
    - `POP`
11. Halt
    - `HLT`
12. No Operation
    - `NOP`

### Addressing Modes
1. **Immediate `I`** - Argument is the value of the operand
2. **Absolute `A`** - Argument is the 16-bit memory address of the operand
3. **ZeroPage `Z`** - Argument is the 8-bit address  (referencing the ZeroPage (`0x0000-0x00ff`)) of the operand
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
The control signals make up the core of the logic of the processor.  Combined together, they form the microcode that makes up each one of the processor instructions.  There are 28 separate control signals and they are stored in Control ROM as 32-bit Control Words.

| Signal | Name	| Purpose |
|---|---|---|
| `0x80000000` | `HLT` | HALT |
| `0x40000000` | `AI` | Accumulator Register In |
| `0x20000000` | `AO` | Accumulator Register Out |
| `0x10000000` | `BI` | Internal Register In |
| `0x08000000` | `II` | Instruction Register In |
| `0x04000000` | `CI2` | Program Counter In from Address Bus |
| `0x02000000` | `CIL` | Program Counter LSB In from Data Bus |
| `0x01000000` | `CIH` | Program Counter MSB In from Data Bus |
| `0x00800000` | `CO2` | Program Counter Out |
| `0x00400000` | `CU` | Program Counter Increment |
| `0x00200000` | `CU2` | Program Counter Increment by 2 |
| `0x00100000` | `MI2` | Memory Address Register In From Address Bus (removed MIL & MIH) |
| `0x00080000` | `MU` | Memory Address Register Increment |
| `0x00040000` | `MU2` | Memory Address Register Increment by 2 |
| `0x00020000` | `RI` | Write to Memory |
| `0x00010000` | `RO` | Read from Memory to Data Bus |
| `0x00008000` | `RO2` | Read 2 bytes from memory to Address Bus |
| `0x00004000` | `PO2` | Stack Pointer Out |
| `0x00002000` | `PU` | Stack Pointer Increment |
| `0x00001000` | `PU2` | Stack Pointer Increment |
| `0x00000800` | `PD` | Stack Pointer Decrement |
| `0x00000400` | `PD2` | Stack Pointer Decrement |
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

`0xB000 - 0xBFFF`   4 KB    Video Character Buffer (80x25 @ 8-bit color)

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







