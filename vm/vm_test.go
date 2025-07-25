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
		{"-5", &object.Integer{Value: -5}},
		{"-10", &object.Integer{Value: -10}},
		{"-50 + 100 + -50", &object.Integer{Value: 0}},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", &object.Integer{Value: 50}},
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

func TestBooleanExpressions(t *testing.T) {
	True := &object.Boolean{Value: true}
	False := &object.Boolean{Value: false}
	tests := []vmTestCase{
		{"true", True},
		{"false", False},
		{"1 < 2", True},
		{"1 > 2", False},
		{"1 < 1", False},
		{"1 > 1", False},
		{"1 == 1", True},
		{"1 != 1", False},
		{"1 == 2", False},
		{"1 != 2", True},
		{"true == true", True},
		{"false == false", True},
		{"true == false", False},
		{"true != false", True},
		{"false != true", True},
		{"(1 < 2) == true", True},
		{"(1 < 2) == false", False},
		{"(1 > 2) == true", False},
		{"(1 > 2) == false", True},
		{"!true", False},
		{"!false", True},
		{"!5", False},
		{"!!true", True},
		{"!!false", False},
		{"!!5", True},
		{"!(if (false) { 5; })", True},
	}
	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", &object.Integer{Value: 10}},
		{"if (true) { 10 } else { 20 }", &object.Integer{Value: 10}},
		{"if (false) { 10 } else { 20 } ", &object.Integer{Value: 20}},
		{"if (1) { 10 }", &object.Integer{Value: 10}},
		{"if (1 < 2) { 10 }", &object.Integer{Value: 10}},
		{"if (1 < 2) { 10 } else { 20 }", &object.Integer{Value: 10}},
		{"if (1 > 2) { 10 } else { 20 }", &object.Integer{Value: 20}},
		{"if (1 > 2) { 10 }", vm.Null},
		{"if (false) { 10 }", vm.Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", &object.Integer{Value: 20}},
	}
	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", &object.Integer{Value: 1}},
		{"let one = 1; let two = 2; one + two", &object.Integer{Value: 3}},
		{"let one = 1; let two = one + one; one + two", &object.Integer{Value: 3}},
	}
	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, &object.String{Value: "monkey"}},
		{`"mon" + "key"`, &object.String{Value: "monkey"}},
		{`"mon" + "key" + "banana"`, &object.String{Value: "monkeybanana"}},
	}
	runVmTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{`[]`, &object.Array{Elements: []object.Object{}}},
		{`[1, 2, 3]`, &object.Array{
			Elements: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
			}},
		},
		{`[1 + 2, 3 - 4, 5 * 6]`, &object.Array{
			Elements: []object.Object{
				&object.Integer{Value: 3},
				&object.Integer{Value: -1},
				&object.Integer{Value: 30},
			}},
		},
	}
	runVmTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}", &object.Hash{Pairs: map[object.HashKey]object.HashPair{}},
		},
		{
			"{1: 2, 2: 3}",
			&object.Hash{Pairs: map[object.HashKey]object.HashPair{
				newHashKey(object.INTEGER_OBJ, uint64(1)): newHashPair(&object.Integer{Value: 1}, &object.Integer{Value: 2}),
				newHashKey(object.INTEGER_OBJ, uint64(2)): newHashPair(&object.Integer{Value: 2}, &object.Integer{Value: 3}),
			}},
		},
		{
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			&object.Hash{Pairs: map[object.HashKey]object.HashPair{
				newHashKey(object.INTEGER_OBJ, uint64(2)): newHashPair(&object.Integer{Value: 2}, &object.Integer{Value: 4}),
				newHashKey(object.INTEGER_OBJ, uint64(6)): newHashPair(&object.Integer{Value: 6}, &object.Integer{Value: 16}),
			}},
		},
	}
	runVmTests(t, tests)
}

func newHashPair(key, value object.Object) object.HashPair {
	return object.HashPair{
		Key:   key,
		Value: value,
	}
}
func newHashKey(objType object.ObjectType, value uint64) object.HashKey {
	return object.HashKey{
		Type:  objType,
		Value: value,
	}
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", &object.Integer{Value: 2}},
		{"[1, 2, 3][0 + 2]", &object.Integer{Value: 3}},
		{"[[1, 1, 1]][0][0]", &object.Integer{Value: 1}},
		{"[][0]", vm.Null},
		{"[1, 2, 3][99]", vm.Null},
		{"[1][-1]", vm.Null},
		{"{1: 1, 2: 2}[1]", &object.Integer{Value: 1}},
		{"{1: 1, 2: 2}[2]", &object.Integer{Value: 2}},
		{"{1: 1}[0]", vm.Null},
		{"{}[0]", vm.Null},
	}
	runVmTests(t, tests)
}
