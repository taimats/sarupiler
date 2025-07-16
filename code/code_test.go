package code_test

import (
	"testing"

	"github.com/taimats/sarupiler/code"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       code.Opcode
		operands []int
		want     []byte
	}{
		{code.OpConstant, []int{65534}, []byte{byte(code.OpConstant), 255, 254}},
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
