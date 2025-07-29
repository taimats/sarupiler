package compiler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taimats/sarupiler/compiler"
)

func TestDefine(t *testing.T) {
	want := map[string]compiler.Symbol{
		"a": {Name: "a", Scope: compiler.GlobalScope, Index: 0},
		"b": {Name: "b", Scope: compiler.GlobalScope, Index: 1},
		"c": {Name: "c", Scope: compiler.LocalScope, Index: 0},
		"d": {Name: "d", Scope: compiler.LocalScope, Index: 1},
		"e": {Name: "e", Scope: compiler.LocalScope, Index: 0},
		"f": {Name: "f", Scope: compiler.LocalScope, Index: 1},
	}
	asrt := assert.New(t)

	g := compiler.NewSymbolTable()
	a := g.Define("a")
	asrt.Equal(want["a"], a)
	b := g.Define("b")
	asrt.Equal(want["b"], b)

	first := compiler.NewEnclosedSymbolTable(g)
	c := first.Define("c")
	asrt.Equal(want["c"], c)
	d := first.Define("d")
	asrt.Equal(want["d"], d)

	second := compiler.NewEnclosedSymbolTable(first)
	e := second.Define("e")
	asrt.Equal(want["e"], e)
	f := second.Define("f")
	asrt.Equal(want["f"], f)
}

func TestResolve(t *testing.T) {
	sut := compiler.NewSymbolTable()
	sut.Define("a")
	sut.Define("b")
	names := []string{"a", "b"}
	want := []compiler.Symbol{
		{Name: "a", Scope: compiler.GlobalScope, Index: 0},
		{Name: "b", Scope: compiler.GlobalScope, Index: 1},
	}
	a := assert.New(t)

	for i, n := range names {
		sym, ok := sut.Resolve(n)
		a.Equal(want[i], sym)
		a.True(ok)
	}
}

func TestResolveLocal(t *testing.T) {
	tests := []struct {
		input string
		want  compiler.Symbol
	}{
		{input: "a", want: compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0}},
		{input: "b", want: compiler.Symbol{Name: "b", Scope: compiler.GlobalScope, Index: 1}},
		{input: "c", want: compiler.Symbol{Name: "c", Scope: compiler.LocalScope, Index: 0}},
		{input: "d", want: compiler.Symbol{Name: "d", Scope: compiler.LocalScope, Index: 1}},
	}
	g := compiler.NewSymbolTable()
	g.Define("a")
	g.Define("b")

	sut := compiler.NewEnclosedSymbolTable(g)
	sut.Define("c")
	sut.Define("d")
	a := assert.New(t)

	for _, tt := range tests {
		got, ok := sut.Resolve(tt.input)
		a.Equal(tt.want, got)
		a.True(ok)
	}
}
func TestDefineResolveBuiltins(t *testing.T) {
	a := assert.New(t)
	tests := []struct {
		input string
		want  compiler.Symbol
	}{
		{"a", compiler.Symbol{"a", compiler.BuiltinScope, 0}},
		{"c", compiler.Symbol{"c", compiler.BuiltinScope, 1}},
		{"e", compiler.Symbol{"e", compiler.BuiltinScope, 2}},
		{"f", compiler.Symbol{"f", compiler.BuiltinScope, 3}},
	}
	global := compiler.NewSymbolTable()
	firstlocal := compiler.NewEnclosedSymbolTable(global)
	secondlocal := compiler.NewEnclosedSymbolTable(firstlocal)

	for i, tt := range tests {
		gSymbol := global.DefineBuiltin(i, tt.input)
		fSymbol := firstlocal.DefineBuiltin(i, tt.input)
		sSymbol := secondlocal.DefineBuiltin(i, tt.input)
		a.Equal(tt.want, gSymbol, "defineBuiltin error: global")
		a.Equal(tt.want, fSymbol, "defineBuiltin error: firstlocal")
		a.Equal(tt.want, sSymbol, "defineBuiltin error: secondlocal")

		gResult, ok := global.Resolve(tt.input)
		a.Equal(tt.want, gResult, "resolve error: global")
		a.True(ok)
		fResult, ok := firstlocal.Resolve(tt.input)
		a.Equal(tt.want, fResult, "resolve error: firstlocal")
		a.True(ok)
		sResult, ok := secondlocal.Resolve(tt.input)
		a.Equal(tt.want, sResult, "resolve error: secondlocal")
		a.True(ok)
	}
}
