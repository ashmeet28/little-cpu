package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	VM_STATUS_UNKNOWN int = iota
	VM_STATUS_READY
	VM_STATUS_RUNNING
	VM_STATUS_HALT
	VM_STATUS_ERROR
)

type VMState struct {
	pc       uint32
	byteCode []byte

	s  []uint32
	sp uint32
	fp uint32

	rs  []uint32
	rsp uint32

	fas  []uint32
	fasp uint32

	frv uint32

	g []uint32

	status int
}

var (
	OP_NOP   byte = 1
	OP_ECALL byte = 2

	OP_ADD byte = 4
	OP_SUB byte = 5
	OP_XOR byte = 6
	OP_OR  byte = 7
	OP_AND byte = 8
	OP_SR  byte = 9
	OP_SL  byte = 10

	OP_PUSH_LITERAL      byte = 12
	OP_PUSH_LOCAL        byte = 13
	OP_PUSH_GLOBAL       byte = 14
	OP_PUSH_FUNC_ARG     byte = 15
	OP_PUSH_FUNC_RET_VAL byte = 16

	OP_POP_LITERAL      byte = 20
	OP_POP_LOCAL        byte = 21
	OP_POP_GLOBAL       byte = 22
	OP_POP_FUNC_ARG     byte = 23
	OP_POP_FUNC_RET_VAL byte = 24

	OP_EQ byte = 28
	OP_NE byte = 29
	OP_LT byte = 30
	OP_GE byte = 31

	OP_JUMP   byte = 32
	OP_CALL   byte = 33
	OP_RETURN byte = 34
)

func VMExecInst(vm VMState) VMState {
	var op byte = vm.byteCode[vm.pc]

	switch op {
	case OP_NOP:
		vm.pc++

	case OP_ECALL:
		vm.pc++
		vm.status = VM_STATUS_HALT

	case OP_ADD:
		vm.s[vm.sp-2] = vm.s[vm.sp-2] + vm.s[vm.sp-1]
		vm.sp--
		vm.pc++

	case OP_SUB:
	case OP_XOR:
	case OP_OR:
	case OP_AND:
		vm.s[vm.sp-2] = vm.s[vm.sp-2] & vm.s[vm.sp-1]
		vm.sp--
		vm.pc++

	case OP_SR:
	case OP_SL:

	case OP_PUSH_LITERAL:
		vm.s[vm.sp] = uint32(vm.byteCode[vm.pc+1]) | (uint32(vm.byteCode[vm.pc+2]) << 8) | (uint32(vm.byteCode[vm.pc+3]) << 16) | (uint32(vm.byteCode[vm.pc+4]) << 24)
		vm.sp++
		vm.pc += 5

	case OP_PUSH_LOCAL:
		vm.s[vm.sp-1] = vm.s[vm.s[vm.sp-1]+vm.fp]
		vm.pc++

	case OP_PUSH_GLOBAL:
		vm.s[vm.sp-1] = vm.g[vm.s[vm.sp-1]]
		vm.pc++

	case OP_PUSH_FUNC_ARG:
		vm.s[vm.sp] = vm.fas[vm.fasp-1]
		vm.sp++
		vm.fasp--
		vm.pc++

	case OP_PUSH_FUNC_RET_VAL:
		vm.s[vm.sp] = vm.frv
		vm.sp++
		vm.pc++

	case OP_POP_LITERAL:
	case OP_POP_LOCAL:
		vm.s[vm.s[vm.sp-2]+vm.fp] = vm.s[vm.sp-1]
		vm.sp -= 2
		vm.pc++

	case OP_POP_GLOBAL:
		vm.g[vm.s[vm.sp-2]] = vm.s[vm.sp-1]
		vm.sp -= 2
		vm.pc++

	case OP_POP_FUNC_ARG:
		vm.fas[vm.fasp] = vm.s[vm.sp-1]
		vm.sp--
		vm.fasp++
		vm.pc++

	case OP_POP_FUNC_RET_VAL:
		vm.frv = vm.s[vm.sp-1]
		vm.sp--
		vm.pc++

	case OP_EQ:
	case OP_NE:
	case OP_LT:
	case OP_GE:

	case OP_JUMP:
		vm.pc = vm.s[vm.sp-1]
		vm.sp--

	case OP_CALL:
		vm.rs[vm.rsp] = vm.pc + 1
		vm.rs[vm.rsp+1] = vm.fp

		vm.rsp += 2

		vm.pc = vm.s[vm.sp-1]
		vm.sp--
		vm.fp = vm.sp

	case OP_RETURN:
		vm.sp = vm.fp

		vm.fp = vm.rs[vm.rsp-1]
		vm.pc = vm.rs[vm.rsp-2]
		vm.rsp -= 2

	default:
		vm.status = VM_STATUS_ERROR
	}

	fmt.Println(op)

	fmt.Println("pc", vm.pc)

	fmt.Println("s", vm.s)
	fmt.Println("sp", vm.sp)
	fmt.Println("fp", vm.fp)

	fmt.Println("rs", vm.rs)
	fmt.Println("rsp", vm.rsp)

	fmt.Println("fas", vm.fas)
	fmt.Println("fasp", vm.fasp)

	fmt.Println("frv", vm.frv)

	fmt.Println("g", vm.g)

	fmt.Println("------------------------------------")

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

	vm.pc = 0

	vm.byteCode = append(vm.byteCode, byteCode...)

	vm.sp = 0
	vm.fp = 0
	vm.s = make([]uint32, 16)

	vm.rsp = 0
	vm.rs = make([]uint32, 16)

	vm.fasp = 0
	vm.fas = make([]uint32, 8)

	vm.frv = 0

	vm.g = make([]uint32, 16)

	vm.status = VM_STATUS_READY

	return vm
}

func GetByteCode(p string) []byte {
	data, err := os.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func main() {
	vm := VMCreate(GetByteCode(os.Args[1]))
	VMRun(vm)
}
