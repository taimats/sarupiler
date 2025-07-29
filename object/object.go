package object

import (
	"fmt"

	"github.com/taimats/sarupiler/code"
	"github.com/taimats/sarupiler/monkey/object"
)

const (
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION_OBJ"
	CLOSURE_OBJ           = "CLOSURE"
)

type CompiledFunction struct {
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
}

func (cf *CompiledFunction) Type() object.ObjectType {
	return COMPILED_FUNCTION_OBJ
}

func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}

type Closure struct {
	Fn   *CompiledFunction
	Free []object.Object //Free is a place where Fn keeps the free variables until runtime.
}

func (c *Closure) Type() object.ObjectType {
	return CLOSURE_OBJ
}

func (c *Closure) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", c)
}
