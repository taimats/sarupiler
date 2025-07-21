package code_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taimats/sarupiler/code"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       code.Opcode
		operands []int
		want     []byte
	}{
		{code.OpConstant, []int{65534}, []byte{byte(code.OpConstant), 255, 254}},
		{code.OpAdd, []int{}, []byte{byte(code.OpAdd)}},
	}

	for _, tt := range tests {
		got := code.Make(tt.op, tt.operands...)
		if len(got) != len(tt.want) {
			t.Errorf("wrong length: (got: %d, want: %d)", len(got), len(tt.want))
		}
		for i, b := range got {
			if b != tt.want[i] {
				t.Errorf("wrong byte: (got: %d, want: %d)", b, tt.want[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []code.Instructions{
		code.Make(code.OpAdd),
		code.Make(code.OpConstant, 2),
		code.Make(code.OpConstant, 65535),
	}
	want := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`
	concatted := code.Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	got := concatted.String()

	assert.Equal(t, want, got)
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        code.Opcode
		operands  []int
		bytesRead int
	}{
		{code.OpConstant, []int{65535}, 2},
	}
	a := assert.New(t)

	for _, tt := range tests {
		instruction := code.Make(tt.op, tt.operands...)
		def, err := code.Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("found no definition: (error: %s)", err)
		}

		got, num := code.ReadOperands(def, instruction[1:])

		a.Equal(tt.operands, got)
		a.Equal(tt.bytesRead, num)
	}
}
