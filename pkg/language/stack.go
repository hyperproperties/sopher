package language

// Represents a first in first out data structure with push and pop.
type Stack[T any] []T

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
