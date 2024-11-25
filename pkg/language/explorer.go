package language

type Explorer[T any] interface {
	Explore(scope Scope, interpreter *Interpreter[T], assertion HyperAssertion[T]) LiftedBoolean
}
