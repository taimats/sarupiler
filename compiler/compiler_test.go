package compiler_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taimats/sarupiler/code"
	"github.com/taimats/sarupiler/compiler"
	"github.com/taimats/sarupiler/monkey/ast"
	"github.com/taimats/sarupiler/monkey/lexer"
	"github.com/taimats/sarupiler/monkey/object"
	"github.com/taimats/sarupiler/monkey/parser"
	obj "github.com/taimats/sarupiler/object"
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
		a.Equal(tt.wantConstants, bytecode.Constants, printConsts(tt.wantConstants, bytecode.Constants))
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func printConsts(want, got []object.Object) string {
	var buf strings.Builder
	buf.WriteString("{ want }\n")
	for _, o := range want {
		buf.WriteString(o.Inspect())
		buf.WriteString("\n")
	}
	buf.WriteString("{ got }\n")
	for _, o := range got {
		buf.WriteString(o.Inspect())
		buf.WriteString("\n")
	}
	return buf.String()
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

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let one = 1;
			let two = 2;
			`,
			wantConstants: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			),
		},
		{
			input: `
			let one = 1;
			one;
			`,
			wantConstants: []object.Object{&object.Integer{Value: 1}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			),
		},
		{
			input: `
			let one = 1;
			let two = one;
			two;
			`,
			wantConstants: []object.Object{&object.Integer{Value: 1}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			),
		},
	}
	runCompilerTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:         `"monkey"`,
			wantConstants: []object.Object{&object.String{Value: "monkey"}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			),
		},
		{
			input:         `"mon" + "key"`,
			wantConstants: []object.Object{&object.String{Value: "mon"}, &object.String{Value: "key"}},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			),
		},
	}
	runCompilerTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:         `[]`,
			wantConstants: []object.Object{},
			wantInstructions: concatInstructions(
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			),
		},
		{
			input: `[1, 2, 3]`,
			wantConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
			},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			),
		},
		{
			input: `[1 + 2, 3 - 4, 5 * 6]`,
			wantConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
				&object.Integer{Value: 5},
				&object.Integer{Value: 6},
			},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			),
		},
	}
	runCompilerTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:         "{}",
			wantConstants: []object.Object{},
			wantInstructions: concatInstructions(
				code.Make(code.OpHash, 0),
				code.Make(code.OpPop),
			),
		},
		{
			input: "{1: 2, 3: 4, 5: 6}",
			wantConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
				&object.Integer{Value: 5},
				&object.Integer{Value: 6},
			},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpHash, 6),
				code.Make(code.OpPop),
			),
		},
		{
			input: "{1: 2 + 3, 4: 5 * 6}",
			wantConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
				&object.Integer{Value: 5},
				&object.Integer{Value: 6},
			},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpHash, 4),
				code.Make(code.OpPop),
			),
		},
	}
	runCompilerTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "[1, 2, 3][1 + 1]",
			wantConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 1},
				&object.Integer{Value: 1},
			},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpAdd),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			),
		},
		{
			input: "{1: 2}[2 - 1]",
			wantConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 2},
				&object.Integer{Value: 1},
			},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpHash, 2),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			),
		},
	}
	runCompilerTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { return 5 + 10 }`,
			wantConstants: []object.Object{
				&object.Integer{Value: 5},
				&object.Integer{Value: 10},
				&obj.CompiledFunction{
					Instructions: concatInstructions(
						code.Make(code.OpConstant, 0),
						code.Make(code.OpConstant, 1),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					),
				},
			},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			),
		},
		{
			input: `fn() { 5 + 10 }`,
			wantConstants: []object.Object{
				&object.Integer{Value: 5},
				&object.Integer{Value: 10},
				&obj.CompiledFunction{
					Instructions: concatInstructions(
						code.Make(code.OpConstant, 0),
						code.Make(code.OpConstant, 1),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					),
				},
			},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			),
		},
		{
			input: `fn() { 1; 2 }`,
			wantConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&obj.CompiledFunction{
					Instructions: concatInstructions(
						code.Make(code.OpConstant, 0),
						code.Make(code.OpPop),
						code.Make(code.OpConstant, 1),
						code.Make(code.OpReturnValue),
					),
				},
			},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			),
		},
	}
	runCompilerTests(t, tests)
}

func TestCompilerScopes(t *testing.T) {
	a := assert.New(t)
	c := compiler.New()

	compiler.Emit(c, code.OpMul)
	compiler.EnterScope(c)
	a.Equal(1, compiler.ScopeIndex(c), "wrong scopeIndex after Emit(code.OpMul)")

	compiler.Emit(c, code.OpSub)
	a.Equal(1, len(compiler.CurrentIns(c)), "wrong instructions length after Emit(code.OpSub)")
	a.Equal(code.OpSub, compiler.LastIns(c).Opcode, "wrong Opcode after Emit(code.OpSub)")

	compiler.LeaveScope(c)
	a.Equal(0, compiler.ScopeIndex(c), "wrong scopeIndex after leaveScope")

	compiler.Emit(c, code.OpAdd)
	a.Equal(2, len(compiler.CurrentIns(c)), "wrong instructions length after Emit(code.OpAdd)")
	a.Equal(code.OpAdd, compiler.LastIns(c).Opcode, "wrong LastIns Opcode after Emit(code.OpAdd)")
	a.Equal(code.OpMul, compiler.PrevIns(c).Opcode, "wrong PrevIns Opcode after Emit(code.OpAdd)")
}

func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { 24 }();`,
			wantConstants: []object.Object{
				&object.Integer{Value: 24},
				&obj.CompiledFunction{
					Instructions: concatInstructions(
						code.Make(code.OpConstant, 0),
						code.Make(code.OpReturnValue),
					),
				},
			},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall),
				code.Make(code.OpPop),
			),
		},
		{
			input: `
			let noArg = fn() { 24 };
			noArg();
			`,
			wantConstants: []object.Object{
				&object.Integer{Value: 24},
				&obj.CompiledFunction{
					Instructions: concatInstructions(
						code.Make(code.OpConstant, 0),
						code.Make(code.OpReturnValue),
					),
				},
			},
			wantInstructions: concatInstructions(
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall),
				code.Make(code.OpPop),
			),
		},
		// {
		// 	input: `
		// 	let oneArg = fn(a) { a };
		// 	oneArg(24);
		// 	`,
		// 	wantConstants: []object.Object{
		// 		&obj.CompiledFunction{
		// 			Instructions: concatInstructions(
		// 				code.Make(code.OpGetLocal, 0),
		// 				code.Make(code.OpReturnValue),
		// 			),
		// 		},
		// 		&object.Integer{Value: 24},
		// 	},
		// 	wantInstructions: concatInstructions(
		// 		code.Make(code.OpConstant, 0),
		// 		code.Make(code.OpSetGlobal, 0),
		// 		code.Make(code.OpGetGlobal, 0),
		// 		code.Make(code.OpConstant, 1),
		// 		code.Make(code.OpCall, 1),
		// 		code.Make(code.OpPop),
		// 	),
		// },
		// {
		// 	input: `
		// 	let manyArg = fn(a, b, c) { a; b; c };
		// 	manyArg(24, 25, 26);
		// 	`,
		// 	wantConstants: []object.Object{
		// 		&obj.CompiledFunction{
		// 			Instructions: concatInstructions(
		// 				code.Make(code.OpGetLocal, 0),
		// 				code.Make(code.OpPop),
		// 				code.Make(code.OpGetLocal, 1),
		// 				code.Make(code.OpPop),
		// 				code.Make(code.OpGetLocal, 2),
		// 				code.Make(code.OpReturnValue),
		// 			),
		// 		},
		// 		&object.Integer{Value: 24},
		// 		&object.Integer{Value: 25},
		// 		&object.Integer{Value: 26},
		// 	},
		// 	wantInstructions: concatInstructions(
		// 		code.Make(code.OpConstant, 0),
		// 		code.Make(code.OpSetGlobal, 0),
		// 		code.Make(code.OpGetGlobal, 0),
		// 		code.Make(code.OpConstant, 1),
		// 		code.Make(code.OpConstant, 2),
		// 		code.Make(code.OpConstant, 3),
		// 		code.Make(code.OpCall, 3),
		// 		code.Make(code.OpPop),
		// 	),
		// },
	}
	runCompilerTests(t, tests)
}
