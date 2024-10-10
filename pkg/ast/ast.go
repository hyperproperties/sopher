package ast

type Node interface{}

type UniversalQuantifier struct {
	offset, size int
	body         Node
}

func NewUniversalQuantifier(offset, size int, body Node) UniversalQuantifier {
	return UniversalQuantifier{
		offset: offset,
		size:   size,
		body:   body,
	}
}

type ExistentialQuantifier struct {
	offset, size int
	body         Node
}

func NewExistentialQuantifier(offset, size int, body Node) ExistentialQuantifier {
	return ExistentialQuantifier{
		offset: offset,
		size:   size,
		body:   body,
	}
}

// P(A) </<= probability
type ProbabilisticQuantifier struct {
	offset, size int
	body         Node
}

func NewProbabilisticQuantifier(offset, size int, body Node) ProbabilisticQuantifier {
	return ProbabilisticQuantifier{
		offset:     offset,
		size:       size,
		body:       body,
	}
}

// P(A|B) >/>= probability
type ConditionalProbabilityQuantifier struct {
	offset, size int
	event, given Node // A|B
}

func NewConditionalProbabilityQuantifier(offset, size int, event, given Node) ConditionalProbabilityQuantifier {
	return ConditionalProbabilityQuantifier{
		offset:      offset,
		size:        size,
		event:       event,
		given:       given,
	}
}

type BinaryBooleanOperator uint8

const (
	BinaryBooleanConjunction = BinaryBooleanOperator(iota)
	BinaryBooleanDisjunction
	BinaryBooleanImplication
	BinaryBooleanBiimplication
)

type BinaryBooleanExpression struct {
	lhs, rhs Node
	operator BinaryBooleanOperator
}

func And(lhs, rhs Node) BinaryBooleanExpression {
	return BinaryBooleanExpression{
		lhs: lhs,
		operator: BinaryBooleanConjunction,
		rhs: rhs,
	}
}

func Or(lhs, rhs Node) BinaryBooleanExpression {
	return BinaryBooleanExpression{
		lhs: lhs,
		operator: BinaryBooleanDisjunction,
		rhs: rhs,
	}
}

func Implication(lhs, rhs Node) BinaryBooleanExpression {
	return BinaryBooleanExpression{
		lhs: lhs,
		operator: BinaryBooleanImplication,
		rhs: rhs,
	}
}

func Biimplication(lhs, rhs Node) BinaryBooleanExpression {
	return BinaryBooleanExpression{
		lhs: lhs,
		operator: BinaryBooleanBiimplication,
		rhs: rhs,
	}
}

type BinaryInequalityOperator uint8

const (
	BinaryInequalityLessThan = BinaryInequalityOperator(iota)
	BinaryInequalityLessThanOrEqual
	BinaryInequalityGreaterThan
	BinaryInequalityGreaterThanOrEqual
)

type BinaryInequality struct {
	lhs, rhs Node
	operator BinaryInequalityOperator
}

func LessThan(lhs, rhs Node) BinaryInequality {
	return BinaryInequality{
		lhs: lhs,
		operator: BinaryInequalityLessThan,
		rhs: rhs,
	}
}

func LessThanOrEqual(lhs, rhs Node) BinaryInequality {
	return BinaryInequality{
		lhs: lhs,
		operator: BinaryInequalityLessThanOrEqual,
		rhs: rhs,
	}
}

func GreaterThan(lhs, rhs Node) BinaryInequality {
	return BinaryInequality{
		lhs: lhs,
		operator: BinaryInequalityGreaterThan,
		rhs: rhs,
	}
}

func GreaterThanOrEqual(lhs, rhs Node) BinaryInequality {
	return BinaryInequality{
		lhs: lhs,
		operator: BinaryInequalityGreaterThanOrEqual,
		rhs: rhs,
	}
}

type ConstantNumber struct {
	value float32
}

func Number(value float32) ConstantNumber {
	return ConstantNumber{ value }
}

type PredicateExpression[T any] struct {
	predicate func(assignments []T) bool
}

func NewPredicateExpression[T any](predicate func(assignments []T) bool) PredicateExpression[T] {
	return PredicateExpression[T]{
		predicate: predicate,
	}
}
