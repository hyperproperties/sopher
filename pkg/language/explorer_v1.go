package language

type ExplorerV1[T any] interface {
	Explore(scope Scope, predicate PredicateHyperAssertion[T]) LiftedBoolean
}
