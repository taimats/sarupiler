package vm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taimats/sarupiler/compiler"
	"github.com/taimats/sarupiler/monkey/ast"
	"github.com/taimats/sarupiler/monkey/lexer"
	"github.com/taimats/sarupiler/monkey/object"
	"github.com/taimats/sarupiler/monkey/parser"
	"github.com/taimats/sarupiler/vm"
)

type vmTestCase struct {
	input string
	want  any
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", &object.Integer{Value: 1}},
		{"2", &object.Integer{Value: 2}},
		{"1 + 2", &object.Integer{Value: 3}},
		{"1 - 2", &object.Integer{Value: -1}},
		{"1 * 2", &object.Integer{Value: 2}},
		{"4 / 2", &object.Integer{Value: 2}},
		{"50 / 2 * 2 + 10 - 5", &object.Integer{Value: 55}},
		{"5 * (2 + 10)", &object.Integer{Value: 60}},
		{"5 + 5 + 5 + 5 - 10", &object.Integer{Value: 10}},
		{"2 * 2 * 2 * 2 * 2", &object.Integer{Value: 32}},
		{"5 * 2 + 10", &object.Integer{Value: 20}},
		{"5 + 2 * 10", &object.Integer{Value: 25}},
		{"5 * (2 + 10)", &object.Integer{Value: 60}},
		// {"-5", -5},
		// {"-10", -10},
		// {"-50 + 100 + -50", 0},
		// {"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()
	a := assert.New(t)
	for _, tt := range tests {
		p := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(p)
		if err != nil {
			t.Fatalf("compiler failed to compile: (error: %s)", err)
		}
		vm := vm.New(comp.Bytecode())

		err = vm.Run()
		if err != nil {
			t.Fatalf("vm failed to run: (error: %s)", err)
		}
		got := vm.LastPoppedStackElem()

		a.Equal(tt.want, got)
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
