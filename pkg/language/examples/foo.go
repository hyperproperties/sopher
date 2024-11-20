package examples

import sopher "github.com/hyperproperties/sopher/pkg/language"

type Foo_ExecutionModel struct {
	a, b int
	x, y int
}

var Foo_Contract sopher.AGHyperContract[Foo_ExecutionModel] = sopher.NewAGHyperContract(sopher.NewAllAssertion[Foo_ExecutionModel](sopher.NewPredicateHyperAssertion(func(assignments []Foo_ExecutionModel) bool { return true })), sopher.NewAllAssertion[Foo_ExecutionModel](sopher.NewPredicateHyperAssertion(func(assignments []Foo_ExecutionModel) bool { return false })))

// assume: true
// guarantee: false
func Foo(a, b int) (x, y int) {
	execution := Foo_ExecutionModel{a: a, b: b}
	call := func(execution Foo_ExecutionModel) Foo_ExecutionModel {
		wrap := func(a, b int) (x, y int) {
			x = a
			y = b + a
			return
		}
		x, y = wrap(execution.a, execution.b)
		execution.x = x
		execution.y = y
		return execution
	}
	assumption, execution, guarantee := Foo_Contract.Call(execution, call)
	if assumption.IsFalse() {
		panic("")
	}
	if guarantee.IsFalse() {
		panic("")
	}
	return execution.x, execution.y
}

type Bar_ExecutionModel struct {
	x    int
	ret0 int
}

var Bar_Contract sopher.AGHyperContract[Bar_ExecutionModel] = sopher.NewAGHyperContract(sopher.NewAllAssertion[Bar_ExecutionModel](sopher.NewUniversalHyperAssertion(1, sopher.NewPredicateHyperAssertion(func(assignments []Bar_ExecutionModel) bool { e := assignments[0]; _ = e; return e.x >= 0 }))), sopher.NewAllAssertion[Bar_ExecutionModel](sopher.NewUniversalHyperAssertion(1, sopher.NewPredicateHyperAssertion(func(assignments []Bar_ExecutionModel) bool { e := assignments[0]; _ = e; return e.ret0 >= 2 }))))

// assume: forall e. e.x >= 0
// guarantee: forall e. e.ret0 >= 2
func Bar(x int) int {
	execution := Bar_ExecutionModel{x: x}
	call := func(execution Bar_ExecutionModel) Bar_ExecutionModel {
		wrap := func(x int) int {
			return x + 2
		}
		ret0 := wrap(execution.x)
		execution.ret0 = ret0
		return execution
	}
	assumption, execution, guarantee := Bar_Contract.Call(execution, call)
	if assumption.IsFalse() {
		panic("")
	}
	if guarantee.IsFalse() {
		panic("")
	}
	return execution.ret0
}
