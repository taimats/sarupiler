package compiler_test

import (
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
	input            string
	wantConstants    []object.Object
	wantInstructions code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
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
			input:         "1; 2",
			wantConstants: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "1 - 2",
			wantConstants: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "1 * 2",
			wantConstants: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "2 / 1",
			wantConstants: []object.Object{&object.Integer{Value: 2}, &object.Integer{Value: 1}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "-1",
			wantConstants: []object.Object{&object.Integer{Value: 1}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
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
		program := parse(tt.input)

		compiler := compiler.New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler failed to Compile: (error: %s)", err)
		}

		bytecode := compiler.Bytecode()
		a.Equal(tt.wantInstructions, bytecode.Instructions, bytecode.Instructions.String())
		a.Equal(tt.wantConstants, bytecode.Constants)
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

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:         "true",
			wantConstants: []object.Object{},
			wantInstructions: concatInstructions(
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "false",
			wantConstants: []object.Object{},
			wantInstructions: concatInstructions(
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "1 > 2",
			wantConstants: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "1 < 2",
			wantConstants: []object.Object{&object.Integer{Value: 2}, &object.Integer{Value: 1}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "1 == 2",
			wantConstants: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "1 != 2",
			wantConstants: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "true == false",
			wantConstants: []object.Object{},
			wantInstructions: concatInstructions(
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "true != false",
			wantConstants: []object.Object{},
			wantInstructions: concatInstructions(
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			),
		},
		{
			input:         "!true",
			wantConstants: []object.Object{},
			wantInstructions: concatInstructions(
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			),
		},
	}
	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			if (true) { 10 }; 3333;
			`,
			wantConstants: []object.Object{&object.Integer{Value: 10}, &object.Integer{Value: 3333}},
			wantInstructions: concatInstructions(
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 10),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJump, 11),
				code.Make(code.OpNull),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			),
		},
		{
			input:         `if (true) { 10 } else { 20 }; 3333;`,
			wantConstants: []object.Object{&object.Integer{Value: 10}, &object.Integer{Value: 20}, &object.Integer{Value: 3333}},
			wantInstructions: concatInstructions(
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 10),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJump, 13),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			),
		},
	}
	runCompilerTests(t, tests)
}
