package language

type Explorer[T any] interface {
	Explore(scope Scope, predicate PredicateHyperAssertion[T]) LiftedBoolean
}
