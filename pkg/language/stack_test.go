package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	var stack Stack[int] = make(Stack[int], 0)
	stack.Push(1)
	assert.Equal(t, 1, len(stack))
	stack.Push(2)
	assert.Equal(t, 2, len(stack))
	top := stack.Pop()
	assert.Equal(t, 2, top)
	assert.Equal(t, 1, len(stack))
	top = stack.Pop()
	assert.Equal(t, 1, top)
	assert.Equal(t, 0, len(stack))
}
