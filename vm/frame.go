package vm

import (
	"github.com/taimats/sarupiler/code"
	obj "github.com/taimats/sarupiler/object"
)

const (
	MaxFrames = 1024
)

// a unit of executable function
type Frame struct {
	fn *obj.CompiledFunction
	ip int //instruction pointer
	bp int //base pointer pointing to the bottom of the stack of the current call frame.
}

func NewFrame(fn *obj.CompiledFunction, basePointer int) *Frame {
	return &Frame{fn: fn, ip: -1, bp: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
