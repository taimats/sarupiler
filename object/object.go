package object

import (
	"fmt"

	"github.com/taimats/sarupiler/code"
	"github.com/taimats/sarupiler/monkey/object"
)

const (
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION_OBJ"
)

type CompiledFunction struct {
	Instructions code.Instructions
}

func (cf *CompiledFunction) Type() object.ObjectType {
	return COMPILED_FUNCTION_OBJ
}

func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}
