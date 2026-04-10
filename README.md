# DJE-8

An 8-bit minicomputer design and toolchain built around TTL/CMOS 74xx series logic chips, complete with custom instruction set architecture, assembler, emulator, and control ROM builder.

## Overview

The DJE-8 is a complete 8-bit computer architecture project that includes both hardware specifications and software tools. The design is inspired by classic 8-bit computers and can be implemented using discrete TTL/CMOS integrated circuits, making it suitable for educational purposes and hobbyist hardware construction.

## Project Description

This project provides a complete ecosystem for an 8-bit computer system:

- **Processor Specification**: Full instruction set architecture (ISA) with 75+ instructions
- **Assembly Language**: Custom assembly syntax with assembler implementation
- **Emulator**: Software emulation of the DJE-8 processor
- **Control ROM Builder**: Tool to generate microcode for hardware implementation
- **Development Tools**: Written in Go for cross-platform compatibility

### Key Features

**Hardware Architecture**
- 8-bit data bus (supports values 0-255 or -128 to 127)
- 16-bit address bus (64KB addressable memory space)
- Accumulator-based processor design
- Hardware stack with 8-bit stack pointer
- 5 status flags: Interrupt Disable, Zero, Carry, Sign (Negative), Overflow
- Memory-mapped I/O for peripherals
- Serial/UART communication interface

**Instruction Set**
- **Data Movement**: Load, store, push, pop operations
- **Arithmetic**: Add, subtract, compare, negate (with and without carry)
- **Logic**: AND, OR, XOR, NOT operations
- **Bit Manipulation**: Logical/arithmetic shifts, rotates
- **Flow Control**: Conditional branches, jumps, subroutines, interrupts
- **Multiple Addressing Modes**: Immediate, Absolute, Zero Page, Memory Indirect

**Processor Registers**
- Accumulator (A): 8-bit primary data register
- Internal Register (B): 8-bit temporary storage
- Instruction Register (IR): 8-bit current instruction
- Program Counter (PC): 16-bit instruction pointer
- Memory Address Register (MAR): 16-bit address latch
- Stack Pointer (SP): 8-bit stack position
- Flags Register: 8-bit status flags

**Control System**
- 28 discrete control signals
- 32-bit control words stored in Control ROM
- Microcode-driven instruction execution
- Variable instruction timing (1-16 clock cycles)

## Technical Stack

- **Language**: Go 1.24.4
- **Module**: damien.live/dje8
- **Dependencies**: golang.org/x/exp (for enhanced utilities)

## File Structure

```
dje-8/
├── architecture.drawio          # Architecture diagram (editable)
├── architecture.drawio.svg      # Architecture diagram (rendered)
├── SPEC.md                      # Complete ISA and hardware specification
├── src/dje8/
│   ├── go.mod                   # Go module definition
│   ├── go.sum                   # Dependency checksums
│   ├── cmd/
│   │   ├── asm/                 # Assembler implementation
│   │   │   ├── main.go
│   │   │   ├── assembly_language_SPEC.md
│   │   │   ├── test.asm
│   │   │   └── test2.asm
│   │   ├── emu/                 # Emulator implementation
│   │   │   └── main.go
│   │   ├── controlrombuilder/   # Microcode ROM generator
│   │   │   └── main.go
│   │   └── test/                # Testing utilities
│   │       └── main.go
│   └── pkg/
│       ├── common/              # Shared types and definitions
│       │   ├── types.go
│       │   ├── opcode_string.go
│       │   ├── flag_string.go
│       │   ├── alumode_string.go
│       │   └── control_string.go
│       └── ucodebuilder/        # Microcode generation library
│           └── ucodebuilder.go
├── LICENSE                      # MIT License
└── README.md                    # This file
```

## Assembly Language

The DJE-8 assembly language features a clean, readable syntax:

### Language Features
- **Instructions**: All uppercase, no abbreviations (e.g., `LODI`, `ADDA`, `JMP`)
- **Labels**: Alphanumeric with underscores, terminated with colon (e.g., `loop:`, `main_entry:`)
- **Comments**: Semicolon to end of line (e.g., `; This is a comment`)
- **Number Formats**:
  - Decimal: `42`
  - Octal: `042` (leading zero)
  - Hexadecimal: `0x42de`
  - Character: `'Z'`
- **Label References**: Support for high/low byte selection (`>label`, `<label`) and offsets (`label+4`, `label-2`)
- **Directives**: `#org` for controlling memory layout

### Example Program
```asm
#org 0x0000
start:
    LODI 42         ; Load immediate value
    STOA counter    ; Store to memory
loop:
    LODA counter    ; Load from memory
    ADDI 1          ; Increment
    STOA counter    ; Store back
    CMPI 100        ; Compare with limit
    BNE loop        ; Branch if not equal
    HALT            ; Stop execution

counter: 0x00
```

## Memory Map

The DJE-8 provides a comprehensive memory-mapped I/O system:

| Address Range | Size | Purpose |
|---|---|---|
| `0x0000-0x00FF` | 256 bytes | Zero Page (fast access) |
| `0xB000-0xBFFF` | 4 KB | Video character buffer (80x25 @ 8-bit color) |
| `0xF000-0xF00F` | 16 bytes | Serial Port 1 (UART) |
| `0xF010-0xF01F` | 16 bytes | Serial Port 2 (UART) |
| `0xF020-0xF02F` | 16 bytes | Keyboard interface |
| `0xF030-0xF3FF` | ~1 KB | Reserved for peripherals |
| `0xF400-0xF7FF` | 1 KB | Expansion slot 1 |
| `0xF800-0xFBFF` | 1 KB | Expansion slot 2 |
| `0xFC00-0xFFFD` | 1 KB | Video controller registers/buffer |
| `0xFFFE-0xFFFF` | 2 bytes | Interrupt vector |

## Tools

### Assembler (`cmd/asm`)
Converts DJE-8 assembly language to machine code.

### Emulator (`cmd/emu`)
Software simulation of the DJE-8 processor for testing and development.

### Control ROM Builder (`cmd/controlrombuilder`)
Generates microcode ROM images for hardware implementation.

## Current State

**Status**: Active development / Design phase  
**Last Updated**: April 2026  
**Architecture**: Fully specified with detailed ISA documentation  
**Software Tools**: Assembler and emulator in active development  
**Hardware**: Design complete, ready for physical implementation

The project has undergone multiple revisions with significant work on the assembler. The instruction set architecture is well-defined with comprehensive documentation. The Go-based toolchain provides cross-platform support for development and testing.

### Implementation Progress

**Completed:**
- ✅ Complete instruction set architecture
- ✅ Control signal definitions and microcode structure
- ✅ Memory map and I/O specification
- ✅ Assembly language syntax specification
- ✅ Assembler implementation (with octal and character literal support)
- ✅ Directive support (`#org`)
- ✅ Architecture diagrams

**In Progress:**
- 🔄 Emulator development
- 🔄 Control ROM builder refinement

**Planned:**
- ⏳ String literal support in assembler
- ⏳ Operand length validation
- ⏳ Rudimentary bootloader and OS
- ⏳ File system implementation

### Stretch Goals
- Integrated video output and keyboard input (text mode: 80x25 characters, 640x480 pixels)
- Memory banking to extend beyond 64KB
- DMA controller for high-speed peripheral access

## Educational Value

The DJE-8 project is ideal for:
- **Computer Architecture Study**: Understanding CPU design from first principles
- **Digital Logic**: Practical application of TTL/CMOS circuits
- **Low-Level Programming**: Assembly language and machine code
- **Retrocomputing**: Building computers with discrete logic chips
- **Embedded Systems**: Microcontroller-style programming patterns

## Design Philosophy

The DJE-8 balances simplicity with capability:
- Simple enough to build with 74xx series chips
- Rich enough instruction set for practical programming
- Inspired by proven 8-bit architectures (6502, 8080 influences)
- Emphasis on educational transparency over performance optimization

---

*This README was auto-generated by Claude on April 10, 2026.*
