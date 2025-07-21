package vm

import (
	"fmt"

	"github.com/taimats/sarupiler/code"
	"github.com/taimats/sarupiler/compiler"
	"github.com/taimats/sarupiler/monkey/object"
)

const StackSize = 2048 //(2KB)

type VM struct {
	constants    []object.Object
	instrcutions code.Instructions

	stack []object.Object
	sp    int //stack pointer always points to the free slot in the stack. Top of stack is stack[sp - 1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instrcutions: bytecode.Instructions,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func (vm *VM) Run() error {
	for i := 0; i < len(vm.instrcutions); i++ {
		op := code.Opcode(vm.instrcutions[i])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instrcutions[i+1:])
			i += 2 //advancing by 2Bytes

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}
