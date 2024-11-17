package language

// TODO: Remove "offset" from assertions.

type HyperAssertionVisitor[T any] interface {
	UniversalHyperAssertion(assertion UniversalHyperAssertion[T])
	ExistentialHyperAssertion(assertion ExistentialHyperAssertion[T])
	BinaryHyperAssertion(assertion BinaryHyperAssertion[T])
	PredicateHyperAssertion(assertion PredicateHyperAssertion[T])
	TrueHyperAssertion(assertion TrueHyperAssertion[T])
}

// HyperAssertion represents an interface for tracking and evaluating the state of
// assignments over a sequence of elements, typically in the context of
// quantifiers in formal verification or monitoring systems.
type HyperAssertion[T any] interface {
	Accept(visitor HyperAssertionVisitor[T])
	Size() int
}

// Compile time checking interface implementations:
var (
	_ HyperAssertion[any] = (*UniversalHyperAssertion[any])(nil)
	_ HyperAssertion[any] = (*ExistentialHyperAssertion[any])(nil)
	_ HyperAssertion[any] = (*PredicateHyperAssertion[any])(nil)
	_ HyperAssertion[any] = (*TrueHyperAssertion[any])(nil)
	_ HyperAssertion[any] = (*BinaryHyperAssertion[any])(nil)
)

func HyperAssertionFromAST[T any](node Node) HyperAssertion[T] {
	var recurse func(node Node, variables int) HyperAssertion[T]
	recurse = func(node Node, variables int) HyperAssertion[T] {
		switch cast := node.(type) {
		case Universal:
			body := recurse(cast.assertion, variables+len(cast.variables))
			monitor := NewUniversalHyperAssertion[T](variables, len(cast.variables), body)
			return monitor
		case Existential:
			body := recurse(cast.assertion, variables+len(cast.variables))
			monitor := NewExistentialHyperAssertion[T](variables, len(cast.variables), body)
			return monitor
		case PredicateExpression[T]:
			monitor := NewPredicateHyperAssertion(cast.predicate)
			return monitor
		}
		panic("unknown or unsupported AST node for the incremental monitor")
	}
	return recurse(node, 0)
}

type TrueHyperAssertion[T any] struct{}

func NewTrueHyperAssertion[T any]() *TrueHyperAssertion[T] {
	return &TrueHyperAssertion[T]{}
}

func (assertion TrueHyperAssertion[T]) Size() int {
	return 0
}

func (assertion TrueHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.TrueHyperAssertion(assertion)
}

type BinaryOperator uint8
const (
	LogicalConjunction = BinaryOperator(iota)
	LogicalDisjunction
	LogicalImplication
	LogicalBiimplication
)

type BinaryHyperAssertion[T any] struct {
	lhs, rhs HyperAssertion[T]
	operator BinaryOperator
}

func (assertion BinaryHyperAssertion[T]) Size() int {
	return assertion.lhs.Size() + assertion.rhs.Size()
}

func (assertion BinaryHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.BinaryHyperAssertion(assertion)
}

func NewBinaryHyperAssertion[T any](
	lhs HyperAssertion[T],
	operator BinaryOperator,
	rhs HyperAssertion[T],
) BinaryHyperAssertion[T] {
	return BinaryHyperAssertion[T]{
		lhs: lhs,
		operator: operator,
		rhs: rhs,
	}
}

type PredicateHyperAssertion[T any] struct {
	predicate func(assignments []T) bool
}

func NewPredicateHyperAssertion[T any](predicate func(assignments []T) bool) *PredicateHyperAssertion[T] {
	return &PredicateHyperAssertion[T]{
		predicate: predicate,
	}
}

func (assertion PredicateHyperAssertion[T]) Size() int {
	return 0
}

func (assertion PredicateHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.PredicateHyperAssertion(assertion)
}

type UniversalHyperAssertion[T any] struct {
	offset, size int
	body         HyperAssertion[T]
	result       LiftedBoolean
}

func NewUniversalHyperAssertion[T any](offset, size int, body HyperAssertion[T]) *UniversalHyperAssertion[T] {
	return &UniversalHyperAssertion[T]{
		offset: offset,
		size:   size,
		body:   body,
		result: LiftedTrue,
	}
}

func (assertion UniversalHyperAssertion[T]) Size() int {
	return assertion.size + assertion.body.Size()
}

func (assertion UniversalHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.UniversalHyperAssertion(assertion)
}

type ExistentialHyperAssertion[T any] struct {
	offset, size int
	body         HyperAssertion[T]
}

func NewExistentialHyperAssertion[T any](offset, size int, body HyperAssertion[T]) *ExistentialHyperAssertion[T] {
	return &ExistentialHyperAssertion[T]{
		offset: offset,
		size:   size,
		body:   body,
	}
}

func (assertion ExistentialHyperAssertion[T]) Size() int {
	return assertion.size + assertion.body.Size()
}

func (assertion ExistentialHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.ExistentialHyperAssertion(assertion)
}
