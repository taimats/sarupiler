package compiler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taimats/sarupiler/code"
	"github.com/taimats/sarupiler/compiler"
	"github.com/taimats/sarupiler/monkey/ast"
	"github.com/taimats/sarupiler/monkey/lexer"
	"github.com/taimats/sarupiler/monkey/parser"
)

type compilerTestCase struct {
	input            string
	wantConstants    []any
	wantInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:         "1 + 2",
			wantConstants: []any{1, 2},
			wantInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}
	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()
	a := assert.New(t)

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := compiler.New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler failed to Compile: (error: %s)", err)
		}

		bytecode := compiler.Bytecode()
		a.Equal(tt.wantInstructions, bytecode.Instructions)
		a.Equal(tt.wantConstants, bytecode.Constants)
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
