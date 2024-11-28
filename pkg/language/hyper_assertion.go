package language

type HyperAssertionQuantitativeVisitor[T any] interface {
	UniversalHyperAssertion(assertion UniversalHyperAssertion[T])
	ExistentialHyperAssertion(assertion ExistentialHyperAssertion[T])
}

type HyperAssertionQualitativeVisitor[T any] interface {
	UnaryHyperAssertion(assertion UnaryHyperAssertion[T])
	BinaryHyperAssertion(assertion BinaryHyperAssertion[T])
	PredicateHyperAssertion(assertion PredicateHyperAssertion[T])
	TrueHyperAssertion(assertion TrueHyperAssertion[T])
	AllAssertion(assertion AllAssertion[T])
	AnyAssertion(assertion AnyAssertion[T])
}

type HyperAssertionVisitor[T any] interface {
	HyperAssertionQuantitativeVisitor[T]
	HyperAssertionQualitativeVisitor[T]
}

// HyperAssertion represents an interface for tracking and evaluating the state of
// assignments over a sequence of elements, typically in the context of
// quantifiers in formal verification or monitoring systems.
type HyperAssertion[T any] interface {
	Accept(visitor HyperAssertionVisitor[T])
}

type HyperAssertionQuantitative[T any] interface {
	Quantitatively(visitor HyperAssertionQuantitativeVisitor[T])
}

type HyperAssertionQualitative[T any] interface {
	Qualitatively(visitor HyperAssertionQualitativeVisitor[T])
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

	_ HyperAssertionQuantitative[any] = (*UniversalHyperAssertion[any])(nil)
	_ HyperAssertionQuantitative[any] = (*ExistentialHyperAssertion[any])(nil)

	_ HyperAssertionQualitative[any] = (*PredicateHyperAssertion[any])(nil)
	_ HyperAssertionQualitative[any] = (*TrueHyperAssertion[any])(nil)
	_ HyperAssertionQualitative[any] = (*BinaryHyperAssertion[any])(nil)
	_ HyperAssertionQualitative[any] = (*AllAssertion[any])(nil)
	_ HyperAssertionQualitative[any] = (*AnyAssertion[any])(nil)
)

type AllAssertion[T any] []HyperAssertion[T]

func NewAllAssertion[T any](assertions ...HyperAssertion[T]) AllAssertion[T] {
	return AllAssertion[T](assertions)
}

func (assertion AllAssertion[T]) Qualitatively(visitor HyperAssertionQualitativeVisitor[T]) {
	visitor.AllAssertion(assertion)
}

func (assertion AllAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.AllAssertion(assertion)
}

type AnyAssertion[T any] []HyperAssertion[T]

func NewAnyAssertion[T any](assertions ...HyperAssertion[T]) AnyAssertion[T] {
	return AnyAssertion[T](assertions)
}

func (assertion AnyAssertion[T]) Qualitatively(visitor HyperAssertionQualitativeVisitor[T]) {
	visitor.AnyAssertion(assertion)
}

func (assertion AnyAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.AnyAssertion(assertion)
}

type TrueHyperAssertion[T any] struct{}

func NewTrueHyperAssertion[T any]() *TrueHyperAssertion[T] {
	return &TrueHyperAssertion[T]{}
}

func (assertion TrueHyperAssertion[T]) Qualitatively(visitor HyperAssertionQualitativeVisitor[T]) {
	visitor.TrueHyperAssertion(assertion)
}

func (assertion TrueHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.TrueHyperAssertion(assertion)
}

type UnaryOperator uint8

const LogicalNegation = UnaryOperator(iota)

func (operator UnaryOperator) String() string {
	switch operator {
	case LogicalNegation: return "¬"
	}
	panic("unknown unary operator")
}


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

func (assertion UnaryHyperAssertion[T]) Qualitatively(visitor HyperAssertionQualitativeVisitor[T]) {
	visitor.UnaryHyperAssertion(assertion)
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

func (operator BinaryOperator) String() string {
	switch operator {
	case LogicalConjunction: return "∧"
	case LogicalDisjunction: return "∨"
	case LogicalImplication: return "→"
	case LogicalBiimplication: return "↔"
	}
	panic("unknown binary operator")
}

type BinaryHyperAssertion[T any] struct {
	lhs, rhs HyperAssertion[T]
	operator BinaryOperator
}

func (assertion BinaryHyperAssertion[T]) Qualitatively(visitor HyperAssertionQualitativeVisitor[T]) {
	visitor.BinaryHyperAssertion(assertion)
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

func (assertion PredicateHyperAssertion[T]) Qualitatively(visitor HyperAssertionQualitativeVisitor[T]) {
	visitor.PredicateHyperAssertion(assertion)
}

func (assertion PredicateHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.PredicateHyperAssertion(assertion)
}

type UniversalHyperAssertion[T any] struct {
	size int
	body HyperAssertion[T]
}

func NewUniversalHyperAssertion[T any](size int, body HyperAssertion[T]) UniversalHyperAssertion[T] {
	return UniversalHyperAssertion[T]{
		size: size,
		body: body,
	}
}

func (assertion UniversalHyperAssertion[T]) Quantitatively(visitor HyperAssertionQuantitativeVisitor[T]) {
	visitor.UniversalHyperAssertion(assertion)
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

func (assertion ExistentialHyperAssertion[T]) Quantitatively(visitor HyperAssertionQuantitativeVisitor[T]) {
	visitor.ExistentialHyperAssertion(assertion)
}

func (assertion ExistentialHyperAssertion[T]) Accept(visitor HyperAssertionVisitor[T]) {
	visitor.ExistentialHyperAssertion(assertion)
}
