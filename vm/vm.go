package vm

import (
	"fmt"

	"github.com/taimats/sarupiler/code"
	"github.com/taimats/sarupiler/compiler"
	"github.com/taimats/sarupiler/monkey/object"
)

const StackSize = 2048 //(2KB)
const GlobalSize = 65536

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants    []object.Object
	instrcutions code.Instructions
	globals      []object.Object

	stack []object.Object
	sp    int //stack pointer always points to the free slot in the stack. Top of stack is stack[sp - 1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instrcutions: bytecode.Instructions,
		globals:      make([]object.Object, GlobalSize),
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func NewWithGlobalStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

func (vm *VM) Run() error {
	//ip is instruction pointer.
	for ip := 0; ip < len(vm.instrcutions); ip++ {
		op := code.Opcode(vm.instrcutions[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instrcutions[ip+1:])
			ip += 2 //advancing by 2Bytes

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}
		case code.OpBang:
			err := vm.executeBangOperation()
			if err != nil {
				return err
			}
		case code.OpMinus:
			err := vm.executeMinusOperation()
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		case code.OpJump:
			pos := int(code.ReadUint16(vm.instrcutions[ip+1:]))
			ip = pos - 1
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.instrcutions[ip+1:]))
			ip += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globIndex := code.ReadUint16(vm.instrcutions[ip+1:])
			ip += 2
			vm.globals[globIndex] = vm.pop()
		case code.OpGetGlobal:
			globIndex := code.ReadUint16(vm.instrcutions[ip+1:])
			ip += 2
			err := vm.push(vm.globals[globIndex])
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	rType := right.Type()
	lType := left.Type()

	if rType == object.INTEGER_OBJ && lType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}
	if rType == object.STRING_OBJ && lType == object.STRING_OBJ {
		return vm.executeBinaryStringOperation(op, left, right)
	}
	return fmt.Errorf("invalid operand type")
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	lv := left.(*object.Integer).Value
	rv := right.(*object.Integer).Value

	var result int64
	switch op {
	case code.OpAdd:
		result = lv + rv
	case code.OpSub:
		result = lv - rv
	case code.OpMul:
		result = lv * rv
	case code.OpDiv:
		result = lv / rv
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}
	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}
	lv := left.(*object.String).Value
	rv := right.(*object.String).Value
	return vm.push(&object.String{Value: lv + rv})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	if right.Type() == object.INTEGER_OBJ && left.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}
	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	lv := left.(*object.Integer).Value
	rv := right.(*object.Integer).Value
	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(lv == rv))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(lv != rv))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(lv > rv))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func (vm *VM) executeBangOperation() error {
	operand := vm.pop()
	switch operand {
	case True:
		vm.push(False)
	case False:
		vm.push(True)
	case Null:
		vm.push(True)
	default:
		vm.push(False)
	}
	return nil
}

func (vm *VM) executeMinusOperation() error {
	operand := vm.pop()
	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}
	v := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -v})
}
