package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type VMState struct {
	pc     uint32
	status int
	reg    []uint32
	mem    []uint32
}

const (
	VM_STATUS_UNKNOWN int = iota
	VM_STATUS_READY
	VM_STATUS_RUNNING
	VM_STATUS_HALT
	VM_STATUS_ERROR
)

func GetByteCode(f string) []byte {
	data, err := os.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func VMExecInst(vm VMState) VMState {
	var (
		ECALL uint32 = 0x38

		ADD uint32 = 0x08
		SUB uint32 = 0x09
		// XOR uint32 = 0x0c
		// OR  uint32 = 0x0e
		// AND uint32 = 0x0f
		// SRA uint32 = 0x0a
		// SRL uint32 = 0x0b
		// SLL uint32 = 0x0d

		// LB  uint32 = 0x10
		// LH  uint32 = 0x11
		LW uint32 = 0x12
		// LBU uint32 = 0x14
		// LHU uint32 = 0x15

		// SB uint32 = 0x18
		// SH uint32 = 0x19
		SW uint32 = 0x1a

		LUI uint32 = 0x21
		// LLI  uint32 = 0x22
		LLIU uint32 = 0x26

		// BEQ  uint32 = 0x28
		// BNE  uint32 = 0x29
		// BLT  uint32 = 0x2c
		// BGE  uint32 = 0x2d
		// BLTU uint32 = 0x2e
		// BGEU uint32 = 0x2f

		// JAL  uint32 = 0x30
		JALR uint32 = 0x31
	)

	var inst uint32 = vm.mem[vm.pc>>2]

	var op uint32 = inst & 0xff
	var p1 uint32 = (inst >> 8) & 0xff
	var p2 uint32 = (inst >> 16) & 0xff
	var p3 uint32 = (inst >> 24) & 0xff
	var temp uint32

	fmt.Println(op, p1, p2, p3)
	fmt.Println(vm.pc)
	fmt.Println(vm.reg)
	fmt.Println(vm.mem[0x2000_0000>>2 : 0x2000_0080>>2])
	fmt.Println(vm.mem[0x3000_0000>>2 : 0x3000_0080>>2])
	fmt.Println("----")

	switch op {
	case ECALL:
		vm.status = VM_STATUS_HALT
		vm.pc += 4
	case ADD:
		vm.reg[p1] = vm.reg[p2] + vm.reg[p3]
		vm.pc += 4
	case SUB:
		vm.reg[p1] = vm.reg[p2] - vm.reg[p3]
		vm.pc += 4
	case LW:
		vm.reg[p1] = vm.mem[(vm.reg[p2]+vm.reg[p3])>>2]
		vm.pc += 4
	case SW:
		vm.mem[(vm.reg[p2]+vm.reg[p3])>>2] = vm.reg[p1]
		vm.pc += 4
	case LUI:
		vm.reg[p1] = (p2 << 16) | (p3 << 24)
		vm.pc += 4
	case LLIU:
		vm.reg[p1] = p2 | (p3 << 8)
		vm.pc += 4
	case JALR:
		temp = vm.pc + 4
		vm.pc = vm.reg[p2] + vm.reg[p3]
		vm.reg[p1] = temp
	default:
		vm.status = VM_STATUS_ERROR
	}

	vm.reg[0] = 0

	return vm
}

func VMRun(vm VMState) {
	if vm.status == VM_STATUS_READY {
		vm.status = VM_STATUS_RUNNING
	}

	for vm.status == VM_STATUS_RUNNING {
		vm = VMExecInst(vm)
		if vm.status == VM_STATUS_ERROR {
			fmt.Println("VM STATUS: ERROR")
		} else if vm.status == VM_STATUS_HALT {
			fmt.Println("VM STATUS: HALT")
		}
		if len(os.Args) > 2 && os.Args[2] == "--debug" {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func VMCreate(byteCode []byte) VMState {
	var vm VMState

	vm.reg = make([]uint32, 32)
	vm.mem = make([]uint32, 0x1000_0000)
	vm.pc = 0x1000_0000
	vm.status = VM_STATUS_READY

	for i, b := range byteCode {
		vm.mem[(vm.pc+uint32(i))>>2] = vm.mem[(vm.pc+uint32(i))>>2] | (uint32(b) << ((i % 4) << 3))
	}

	return vm
}

func main() {
	vm := VMCreate(GetByteCode(os.Args[1]))
	VMRun(vm)
}
