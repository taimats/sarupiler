package compiler

import (
	"github.com/taimats/sarupiler/code"
)

// Test helper methods for unexported compiler method
func Emit(c *Compiler, op code.Opcode, operand ...int) {
	c.emit(op, operand...)
}
func CurrentIns(c *Compiler) code.Instructions {
	return c.currentInstructions()
}
func LastIns(c *Compiler) EmittedInstruction {
	return c.scopes[c.scopeIndex].lastInstruction
}
func PrevIns(c *Compiler) EmittedInstruction {
	return c.scopes[c.scopeIndex].previousInstruction
}
func ScopeIndex(c *Compiler) int {
	return c.scopeIndex
}
func EnterScope(c *Compiler) {
	c.enterScope()
}
func LeaveScope(c *Compiler) {
	c.leaveScope()
}
