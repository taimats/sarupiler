package compiler_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taimats/sarupiler/code"
	"github.com/taimats/sarupiler/compiler"
	"github.com/taimats/sarupiler/monkey/ast"
	"github.com/taimats/sarupiler/monkey/lexer"
	"github.com/taimats/sarupiler/monkey/object"
	"github.com/taimats/sarupiler/monkey/parser"
)

type compilerTestCase struct {
	name             string
	input            string
	wantConstants    []object.Object
	wantInstructions code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:          "OpAdd",
			input:         "1 + 2",
			wantConstants: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			),
		},
		{
			name:          "OpPop",
			input:         "1; 2",
			wantConstants: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			),
		},
	}
	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()
	a := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			program := parse(tt.input)

			compiler := compiler.New()
			err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("compiler failed to Compile: (error: %s)", err)
			}

			bytecode := compiler.Bytecode()
			a.Equal(tt.wantInstructions, bytecode.Instructions, fmt.Sprintf("instructions not equal in %s:", tt.name))
			a.Equal(tt.wantConstants, bytecode.Constants, fmt.Sprintf("constants not equal in %s:", tt.name))
		})
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func concatInstructions(ins ...code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, item := range ins {
		out = append(out, item...)
	}
	return out
}
