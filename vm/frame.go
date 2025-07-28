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
}

func NewFrame(fn *obj.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
