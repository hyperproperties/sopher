package language

type Stack[T any] []T

func (stack *Stack[T]) Push(elements ...T) {
	*stack = append(*stack, elements...)
}

func (stack *Stack[T]) Pop() T {
	if len(*stack) == 0 {
		panic("empty stack")
	}
	length := len(*stack)
	top := (*stack)[length-1]
	(*stack) = (*stack)[0 : length-1]
	return top
}

func (stack *Stack[T]) Peek() T {
	if len(*stack) == 0 {
		panic("empty stack")
	}
	length := len(*stack)
	top := (*stack)[length-1]
	return top
}
