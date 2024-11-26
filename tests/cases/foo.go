package main

import sopher "github.com/hyperproperties/sopher/pkg/language"

type Foo_ExecutionModel struct {
	low, high int
	ret0      int
}

var Foo_Contract sopher.AGHyperContract[Foo_ExecutionModel] = sopher.NewAGHyperContract(sopher.NewAllAssertion[Foo_ExecutionModel](sopher.NewUniversalHyperAssertion(1, sopher.NewPredicateHyperAssertion(func(assignments []Foo_ExecutionModel) bool {
	e := assignments[0]
	_ = e
	return e.low > 0 && e.high > 0
}))), sopher.NewAllAssertion[Foo_ExecutionModel](sopher.NewUniversalHyperAssertion(2, sopher.NewBinaryHyperAssertion[Foo_ExecutionModel](sopher.NewPredicateHyperAssertion(func(assignments []Foo_ExecutionModel) bool {
	e0, e1 := assignments[0], assignments[1]
	_, _ = e0, e1
	return e0.low == e1.low
}), sopher.LogicalImplication, sopher.NewPredicateHyperAssertion(func(assignments []Foo_ExecutionModel) bool {
	e0, e1 := assignments[0], assignments[1]
	_, _ = e0, e1
	return e0.ret0 == e1.ret0
})))))

// assume: forall e. e.low > 0 && e.high > 0
// guarantee: forall e0 e1. (e0.low == e1.low; -> e0.ret0 == e1.ret0;)
func Foo(low, high int) int {
	caller := sopher.Caller()
	execution := Foo_ExecutionModel{low: low, high: high}
	call := func(execution Foo_ExecutionModel) Foo_ExecutionModel {
		wrap := func(low, high int) int {
			if low < 0 {
				return low
			}
			return high
		}
		ret0 := wrap(execution.low, execution.high)
		execution.ret0 = ret0
		return execution
	}
	assumption, execution, guarantee := Foo_Contract.Call(caller, execution, call)
	if assumption.IsFalse() {
		panic("")
	}
	if guarantee.IsFalse() {
		panic("")
	}
	return execution.ret0
}
