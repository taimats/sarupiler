package compiler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taimats/sarupiler/compiler"
)

func TestDefine(t *testing.T) {
	sut := compiler.NewSymbolTable()
	want := map[string]compiler.Symbol{
		"a": {Name: "a", Scope: compiler.GlobalScope, Index: 0},
		"b": {Name: "b", Scope: compiler.GlobalScope, Index: 1},
	}
	a := assert.New(t)

	got1 := sut.Define("a")
	got2 := sut.Define("b")

	a.Equal(want["a"], got1)
	a.Equal(want["b"], got2)
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
