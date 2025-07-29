package vm

import (
	"fmt"

	"github.com/taimats/sarupiler/code"
	"github.com/taimats/sarupiler/compiler"
	"github.com/taimats/sarupiler/monkey/object"
	obj "github.com/taimats/sarupiler/object"
)

const (
	StackSize  = 2048 //(2KB)
	GlobalSize = 65536
)

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants []object.Object

	globals []object.Object

	stack []object.Object
	sp    int //stack pointer always points to the free slot in the stack. Top of stack is stack[sp - 1]

	frames      []*Frame
	framesIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	cl := &obj.Closure{Fn: &obj.CompiledFunction{Instructions: bytecode.Instructions}}
	frames := make([]*Frame, MaxFrames)
	frames[0] = NewFrame(cl, 0)

	return &VM{
		constants:   bytecode.Constants,
		globals:     make([]object.Object, GlobalSize),
		stack:       make([]object.Object, StackSize),
		sp:          0,
		frames:      frames,
		framesIndex: 1,
	}
}

func NewWithGlobalStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

func (vm *VM) Run() error {
	//ip is instruction pointer.
	var ip int
	var ins code.Instructions
	var op code.Opcode
	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2 //advancing by 2Bytes

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
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			vm.globals[globIndex] = vm.pop()
		case code.OpGetGlobal:
			globIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			err := vm.push(vm.globals[globIndex])
			if err != nil {
				return err
			}
		case code.OpArray:
			numElems := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			array := vm.buildArray(vm.sp-numElems, vm.sp)
			vm.sp = vm.sp - numElems
			err := vm.push(array)
			if err != nil {
				return err
			}
		case code.OpHash:
			numElems := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			hash, err := vm.buildHash(vm.sp-numElems, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = vm.sp - numElems
			err = vm.push(hash)
			if err != nil {
				return err
			}
		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()
			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}
		case code.OpCall:
			numArgs := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			err := vm.executeCall(int(numArgs))
			if err != nil {
				return err
			}
		case code.OpReturnValue:
			returnValue := vm.pop()
			frame := vm.popFrame()
			vm.sp = frame.bp - 1
			err := vm.push(returnValue)
			if err != nil {
				return err
			}
		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.bp - 1
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetLocal:
			localIndex := int(code.ReadUint8(ins[ip+1:]))
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			vm.stack[frame.bp+int(localIndex)] = vm.pop()
		case code.OpGetLocal:
			localIndex := int(code.ReadUint8(ins[ip+1:]))
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			err := vm.push(vm.stack[frame.bp+localIndex])
			if err != nil {
				return err
			}
		case code.OpGetBuiltin:
			builtinIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			def := obj.Builtins[builtinIndex]
			err := vm.push(def.Builtin)
			if err != nil {
				return err
			}
		case code.OpClosure:
			constIndex := code.ReadUint16((ins[ip+1:]))
			numFree := code.ReadUint8(ins[ip+3:])
			vm.currentFrame().ip += 3
			err := vm.pushClosure(int(constIndex), int(numFree))
			if err != nil {
				return err
			}
		case code.OpGetFree:
			freeIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			currenClosure := vm.currentFrame().cl
			err := vm.push(currenClosure.Free[freeIndex])
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

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elems := make([]object.Object, endIndex-startIndex)
	for i := startIndex; i < endIndex; i++ {
		elems[i-startIndex] = vm.stack[i]
	}
	return &object.Array{Elements: elems}
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashPairs := make(map[object.HashKey]object.HashPair)
	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("invalid hash key: %s", key.Type())
		}
		hashPairs[hashKey.HashKey()] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: hashPairs}, nil
}

func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("invalid index operator: %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(left, index object.Object) error {
	array := left.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(array.Elements) - 1)
	if i < 0 || i > max {
		return vm.push(Null)
	}
	return vm.push(array.Elements[i])
}

func (vm *VM) executeHashIndex(left, index object.Object) error {
	hash := left.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("invalid hash key: %s", index.Type())
	}
	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}
	return vm.push(pair.Value)
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func (vm *VM) executeCall(numArgs int) error {
	switch callee := vm.stack[vm.sp-1-numArgs].(type) {
	case *obj.Closure:
		return vm.callFunction(callee, numArgs)
	case *object.Builtin:
		return vm.callBuiltin(callee, numArgs)
	default:
		return fmt.Errorf("calling non-function and non-builtin")
	}
}

func (vm *VM) callFunction(cl *obj.Closure, numArgs int) error {
	if numArgs != cl.Fn.NumParameters {
		return fmt.Errorf("wrong number of args: (got=%d, want=%d)", numArgs, cl.Fn.NumParameters)
	}
	frame := NewFrame(cl, vm.sp-numArgs)
	vm.pushFrame(frame)
	vm.sp = frame.bp + cl.Fn.NumLocals //allocating space on the stack
	return nil
}

func (vm *VM) callBuiltin(builtin *object.Builtin, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]
	result := builtin.Fn(args...)
	if result == nil {
		return vm.push(Null)
	}
	return vm.push(result)
}

func (vm *VM) pushClosure(constIndex int, numFree int) error {
	constant := vm.constants[constIndex]
	cf, ok := constant.(*obj.CompiledFunction)
	if !ok {
		return fmt.Errorf("not a function: %+v", constant)
	}
	free := make([]object.Object, 0, numFree)
	for i := range numFree {
		free = append(free, vm.stack[vm.sp-numFree+i])
	}
	vm.sp = vm.sp - numFree
	cl := &obj.Closure{Fn: cf, Free: free}
	return vm.push(cl)
}
