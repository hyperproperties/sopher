package language

import "iter"

type Explorer[T any] interface {
	UniversalHyperAssertion(assertion UniversalHyperAssertion[T])
	ExistentialHyperAssertion(assertion ExistentialHyperAssertion[T])
	Explore(predicate PredicateHyperAssertion[T]) iter.Seq[LiftedBoolean] 
}