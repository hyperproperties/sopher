package ast

import "github.com/hyperproperties/sopher/pkg/iterx"

// TODO: Sequential hypothesis testing (SPRT)
// https://en.wikipedia.org/wiki/Sequential_probability_ratio_test
type Interpreter[T any] struct {
	assignments []T
	elements    []T
}

func NewInterpreter[T any](size int, elements []T) Interpreter[T] {
	return Interpreter[T]{
		assignments: make([]T, size),
		elements:    elements,
	}
}

func (interpreter *Interpreter[T]) Model(node Node) (bool, []T) {
	satisfied := interpreter.Check(node)
	if satisfied {
		return satisfied, nil
	}
	
	counter := make([]T, len(interpreter.assignments))
	copy(counter, interpreter.assignments)
	return satisfied, counter
}

func (interpreter *Interpreter[T]) Check(node Node) bool {
	switch cast := node.(type) {
	case UniversalQuantifier:
		return interpreter.UniversalQuantifier(cast)
	case ExistentialQuantifier:
		return interpreter.ExistentialQuantifier(cast)
	case PredicateExpression[T]:
		return interpreter.PredicateExpression(cast)
	case BinaryBooleanExpression:
		return interpreter.BinaryBooleanExpression(cast)
	case BinaryInequality:
		return interpreter.BinaryInequality(cast)
	}
	panic("unknown or unsupported AST")
}

func (interpreter *Interpreter[T]) Number(node Node) float32 {
	switch cast := node.(type) {
	case ProbabilisticQuantifier:
		return interpreter.ProbabilisticQuantifier(cast)
	case ConditionalProbabilityQuantifier:
		return interpreter.ConditionalProbabilityQuantifier(cast)
	case ConstantNumber:
		return cast.value
	}

	panic("unknown or unsupported AST")
}

func (interpreter *Interpreter[T]) UniversalQuantifier(universal UniversalQuantifier) bool {
	for permutation := range iterx.Permutations(universal.size, len(interpreter.elements)) {
		for idx, element := range permutation {
			interpreter.assignments[universal.offset+idx] = interpreter.elements[element]
		}

		if !interpreter.Check(universal.body) {
			return false
		}
	}

	return true
}

func (interpreter *Interpreter[T]) ExistentialQuantifier(existential ExistentialQuantifier) bool {
	for permutation := range iterx.Permutations(existential.size, len(interpreter.elements)) {
		for idx, element := range permutation {
			interpreter.assignments[existential.offset+idx] = interpreter.elements[element]
		}

		if interpreter.Check(existential.body) {
			return true
		}
	}

	return false
}

func (interpreter *Interpreter[T]) ProbabilisticQuantifier(probabilistic ProbabilisticQuantifier) float32 {
	var total, satisifed uint = 0, 0
	for permutation := range iterx.Permutations(probabilistic.size, len(interpreter.elements)) {
		for idx, element := range permutation {
			interpreter.assignments[probabilistic.offset+idx] = interpreter.elements[element]
		}

		if interpreter.Check(probabilistic.body) {
			satisifed += 1
		}
		total += 1
	}

	return float32(satisifed) / float32(total)
}

func (interpreter *Interpreter[T]) ConditionalProbabilityQuantifier(conditional ConditionalProbabilityQuantifier) float32 {
	var total, joint, marginal uint = 0, 0, 0
	for permutation := range iterx.Permutations(conditional.size, len(interpreter.elements)) {
		for idx, element := range permutation {
			interpreter.assignments[conditional.offset+idx] = interpreter.elements[element]
		}

		event := interpreter.Check(conditional.event)
		given := interpreter.Check(conditional.given)

		if given {
			marginal += 1
			if event {
				joint += 1
			}
		}
		total += 1
	}

	if marginal == 0 {
		return 0
	}

	return float32(joint) / float32(marginal)

}

func (interpreter *Interpreter[T]) BinaryBooleanExpression(binary BinaryBooleanExpression) bool {
	switch binary.operator {
	case BinaryBooleanConjunction:
		return interpreter.Check(binary.lhs) && interpreter.Check(binary.rhs)
	case BinaryBooleanDisjunction:
		return interpreter.Check(binary.lhs) || interpreter.Check(binary.rhs)
	case BinaryBooleanImplication:
		return !interpreter.Check(binary.lhs) || interpreter.Check(binary.rhs)
	case BinaryBooleanBiimplication:
		return interpreter.Check(binary.lhs) == interpreter.Check(binary.rhs)
	}
	panic("unkown binary operator")
}

func (interpreter *Interpreter[T]) BinaryInequality(binary BinaryInequality) bool {
	lhs := interpreter.Number(binary.lhs)
	rhs := interpreter.Number(binary.rhs)
	switch binary.operator {
	case BinaryInequalityLessThan:
		return lhs < rhs
	case BinaryInequalityLessThanOrEqual:
		return lhs <= rhs
	case BinaryInequalityGreaterThan:
		return lhs > rhs
	case BinaryInequalityGreaterThanOrEqual:
		return lhs >= rhs
	}

	panic("unknown or unspported binary numeric comparison operator")
}

func (interpreter *Interpreter[T]) PredicateExpression(predicate PredicateExpression[T]) bool {
	return predicate.predicate(interpreter.assignments)
}