package language

type Assignments[T any] []T

func NewAssignments[T any](valuations ...T) Assignments[T] {
	return valuations
}

func (assignments *Assignments[T]) Assign(offset int, valuations ...T) {
	copy((*assignments)[offset:], valuations)
}

func (assignments *Assignments[T]) Expand(expansion int) (from, to int) {
	length := len(*assignments)
	*assignments = append(*assignments, make([]T, expansion)...)
	return length, length+expansion
}

func (assignments *Assignments[T]) Shrink(shrinkage int) {
	length := len(*assignments)
	(*assignments) = (*assignments)[0 : length-shrinkage]
}