package language

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	var stack Stack[int] = make(Stack[int], 0)

	for counter := 0; counter < 10000; counter++ {
		switch rand.IntN(4) {
		case 0: // Push an element to the stack.
			oldLength := len(stack)
			element := rand.Int()
			stack.Push(element)
			assert.Equal(t, oldLength+1, len(stack))
			assert.Equal(t, element, stack.Peek())
		case 1: // Pop an element to the stack.
			if len(stack) > 0 {
				top := stack.Peek()
				oldLength := len(stack)
				element := stack.Pop()
				assert.Len(t, stack, oldLength-1)
				assert.Equal(t, top, element)
			} else {
				assert.Panics(t, func() {
					stack.Pop()
				})
			}
		case 2:
			amount := rand.IntN(5)
			if amount > 0 && len(stack) >= amount {
				oldLength := len(stack)
				stack.PopN(amount)
				assert.Len(t, stack, oldLength-amount)
			} else {
				assert.Panics(t, func() {
					stack.PopN(amount)
				})
			}
		case 3: // Peek the top of the stack.
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				element := stack.Peek()
				assert.Equal(t, top, element)
			} else {
				assert.Panics(t, func() {
					stack.Peek()
				})
			}
		}
	}
}
