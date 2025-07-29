package object

import (
	"fmt"

	"github.com/taimats/sarupiler/monkey/object"
)

var Builtins = []struct {
	Name    string
	Builtin *object.Builtin
}{
	{
		"len",
		&object.Builtin{Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of args: (got=%d want=1)", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("unsupported type for len(): (type=%s)", arg.Type())
			}
		}},
	},
	{
		"puts",
		&object.Builtin{Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return nil
		}},
	},
	{
		"first",
		&object.Builtin{Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of args: (got=%d want=1)", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("an arg must be ARRAY_OBJ: (got=%s)", args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return nil
		}},
	},
	{
		"last",
		&object.Builtin{Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of args: (got=%d want=1)", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("an arg must be ARRAY_OBJ: (got=%s)", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return nil
		}},
	},
	{
		"rest",
		&object.Builtin{Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of args: (got=%d want=1)", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("an arg must be ARRAY_OBJ: (got=%s)", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElems := make([]object.Object, length-1)
				copy(newElems, arr.Elements[1:])
				return &object.Array{Elements: newElems}
			}
			return nil
		}},
	},
	{
		"push",
		&object.Builtin{Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of args: (got=%d want=2)", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("an arg must be ARRAY_OBJ: (got=%s)", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length >= 0 {
				newElems := make([]object.Object, length+1)
				copy(newElems, arr.Elements)
				newElems[length] = args[1]
				return &object.Array{Elements: newElems}
			}
			return nil
		}},
	},
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *object.Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
