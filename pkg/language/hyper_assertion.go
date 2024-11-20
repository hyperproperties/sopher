package language

type HyperAssertionVisitor[T any] interface {
	UniversalHyperAssertion(assertion UniversalHyperAssertion[T])
	ExistentialHyperAssertion(assertion ExistentialHyperAssertion[T])
	UnaryHyperAssertion(assertion UnaryHyperAssertion[T])
	BinaryHyperAssertion(assertion BinaryHyperAssertion[T])
	PredicateHyperAssertion(assertion PredicateHyperAssertion[T])
	TrueHyperAssertion(assertion TrueHyperAssertion[T])
	AllAssertion(assertion AllAssertion[T])
	AnyAssertion(assertion AnyAssertion[T])
}

// HyperAssertion represents an interface for tracking and evaluating the state of
// assignments over a sequence of elements, typically in the context of
// quantifiers in formal verification or monitoring systems.
type HyperAssertion[T any] interface {
	Accept(visitor HyperAssertionVisitor[T])
}

// Compile time checking interface implementations:
var (
	_ HyperAssertion[any] = (*UniversalHyperAssertion[any])(nil)
	_ HyperAssertion[any] = (*ExistentialHyperAssertion[any])(nil)
	_ HyperAssertion[any] = (*PredicateHyperAssertion[any])(nil)
	_ HyperAssertion[any] = (*TrueHyperAssertion[any])(nil)
	_ HyperAssertion[any] = (*BinaryHyperAssertion[any])(nil)
	_ HyperAssertion[any] = (*AllAssertion[any])(nil)
	_ HyperAssertion[any] = (*AnyAssertion[any])(nil)
)

type AllAssertion[T any] []HyperAssertion[T]

func NewAllAssertion[T any](assertions ...HyperAssertion[T]) AllAssertion[T] {
	return AllAssertion[T](assertions)
}

func (assertion AllAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.AllAssertion(assertion)
}

type AnyAssertion[T any] []HyperAssertion[T]

func NewAnyAssertion[T any](assertions ...HyperAssertion[T]) AnyAssertion[T] {
	return AnyAssertion[T](assertions)
}

func (assertion AnyAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.AnyAssertion(assertion)
}

type TrueHyperAssertion[T any] struct{}

func NewTrueHyperAssertion[T any]() *TrueHyperAssertion[T] {
	return &TrueHyperAssertion[T]{}
}

func (assertion TrueHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.TrueHyperAssertion(assertion)
}

type UnaryOperator uint8

const LogicalNegation = UnaryOperator(iota)

type UnaryHyperAssertion[T any] struct {
	operator UnaryOperator
	operand  HyperAssertion[T]
}

func NewUnaryHyperAssertion[T any](operator UnaryOperator, operand HyperAssertion[T]) UnaryHyperAssertion[T] {
	return UnaryHyperAssertion[T]{
		operator: operator,
		operand:  operand,
	}
}

func (assertion UnaryHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.UnaryHyperAssertion(assertion)
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

func (assertion BinaryHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.BinaryHyperAssertion(assertion)
}

func NewBinaryHyperAssertion[T any](
	lhs HyperAssertion[T],
	operator BinaryOperator,
	rhs HyperAssertion[T],
) BinaryHyperAssertion[T] {
	return BinaryHyperAssertion[T]{
		lhs:      lhs,
		operator: operator,
		rhs:      rhs,
	}
}

type PredicateHyperAssertion[T any] struct {
	predicate func(assignments []T) bool
}

func NewPredicateHyperAssertion[T any](predicate func(assignments []T) bool) PredicateHyperAssertion[T] {
	return PredicateHyperAssertion[T]{
		predicate: predicate,
	}
}

func (assertion PredicateHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.PredicateHyperAssertion(assertion)
}

type UniversalHyperAssertion[T any] struct {
	size   int
	body   HyperAssertion[T]
	result LiftedBoolean
}

func NewUniversalHyperAssertion[T any](size int, body HyperAssertion[T]) UniversalHyperAssertion[T] {
	return UniversalHyperAssertion[T]{
		size:   size,
		body:   body,
		result: LiftedTrue,
	}
}

func (assertion UniversalHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.UniversalHyperAssertion(assertion)
}

type ExistentialHyperAssertion[T any] struct {
	size int
	body HyperAssertion[T]
}

func NewExistentialHyperAssertion[T any](size int, body HyperAssertion[T]) ExistentialHyperAssertion[T] {
	return ExistentialHyperAssertion[T]{
		size: size,
		body: body,
	}
}

func (assertion ExistentialHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.ExistentialHyperAssertion(assertion)
}
