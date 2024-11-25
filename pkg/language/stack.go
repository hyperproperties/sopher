package language

// Represents a first in first out data structure with push and pop.
type Stack[T any] []T

func NewStack[T any](stack ...T) Stack[T] {
	return stack
}

// Pushes the elements to the stack.
func (stack *Stack[T]) Push(elements ...T) {
	*stack = append(*stack, elements...)
}

// Pops a single element to the stack and returns it.
// If the stack is empty "len(stack) == 0" then it panics.
func (stack *Stack[T]) Pop() T {
	length := len(*stack)
	top := (*stack)[length-1]
	(*stack) = (*stack)[0 : length-1]
	return top
}

// Returns the top-most element of the stack without popping it.
// If the stack is empty "len(stack) == 0" then it panics.
func (stack *Stack[T]) Peek() T {
	length := len(*stack)
	top := (*stack)[length-1]
	return top
}

// Returns the index of the top (or head) of the stack.
// If the stack is empty then -1 is returned.
func (stack Stack[T]) Top() int {
	return len(stack) - 1
}

// Returns true if there are no elements in the stack.
func (stack Stack[T]) IsEmpty() bool {
	return len(stack) == 0
}
