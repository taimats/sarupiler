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
	cl *obj.Closure
	ip int //instruction pointer
	bp int //base pointer pointing to the bottom of the stack of the current call frame.
}

func NewFrame(cl *obj.Closure, basePointer int) *Frame {
	return &Frame{cl: cl, ip: -1, bp: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
