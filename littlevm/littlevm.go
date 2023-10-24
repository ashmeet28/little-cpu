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
	sp       int
	fp       int
	s        []uint32
	g        []uint32
	rasp     int
	ras      []uint32
	status   int
}

func VMExecInst(vm VMState) VMState {
	var (
		NOP   byte = 1
		ECALL byte = 2

		ADD byte = 4
		SUB byte = 5
		XOR byte = 6
		OR  byte = 7
		AND byte = 8
		SR  byte = 9
		SL  byte = 10

		PUSH_LITERAL byte = 16
		PUSH_LOCAL   byte = 17
		PUSH_GLOBAL  byte = 18

		POP_LITERAL byte = 20
		POP_LOCAL   byte = 21
		POP_GLOBAL  byte = 22

		EQ byte = 24
		NE byte = 25
		LT byte = 26
		GE byte = 27

		JUMP   byte = 28
		CALL   byte = 29
		RETURN byte = 30
	)

	var op byte = vm.byteCode[vm.pc]

	switch op {
	case NOP:
	case ECALL:

	case ADD:
	case SUB:
	case XOR:
	case OR:
	case AND:
	case SR:
	case SL:

	case PUSH_LITERAL:
	case PUSH_LOCAL:
	case PUSH_GLOBAL:

	case POP_LITERAL:
	case POP_LOCAL:
	case POP_GLOBAL:

	case EQ:
	case NE:
	case LT:
	case GE:

	case JUMP:
	case CALL:
	case RETURN:

	default:
		vm.status = VM_STATUS_ERROR
	}

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
	vm.s = make([]uint32, 16777216)

	vm.g = make([]uint32, 16777216)

	vm.rasp = 0
	vm.ras = make([]uint32, 65536)

	vm.status = VM_STATUS_READY

	return vm
}

func GetByteCode(f string) []byte {
	data, err := os.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func main() {
	vm := VMCreate(GetByteCode(os.Args[1]))
	VMRun(vm)
}
