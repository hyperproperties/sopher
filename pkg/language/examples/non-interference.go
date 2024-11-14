package examples

import (
	_ "math"

	sopher "github.com/hyperproperties/sopher/pkg/language"
)

type Retain_ExecutionModel struct {
	lowIn, highIn   int
	lowOut, highOut int
}

var Retain_Contract sopher.AGHyperContract[Retain_ExecutionModel] = sopher.NewAGHyperContract(
	[]sopher.HyperAssertion[Retain_ExecutionModel]{},
	[]sopher.HyperAssertion[Retain_ExecutionModel]{sopher.NewUniversalHyperAssertion[Retain_ExecutionModel](0, 2, sopher.NewPredicateHyperAssertion(func(assignments []Retain_ExecutionModel) bool {
		e0, e1 := assignments[0], assignments[1]
		_, _ = e0, e1
		return !(e0.highIn == e1.highIn) || (e0.lowOut == e1.lowOut)
	}))})

// guarantee: forall e0 e1. !(e0.highIn == e1.highIn) || (e0.lowOut == e1.lowOut)
func Retain(lowIn, highIn int) (lowOut, highOut int) {
	wrap := func(lowIn, highIn int) (lowOut, highOut int) {
		lowOut = lowIn
		highOut = highIn + lowIn
		return
	}
	execution := Retain_ExecutionModel{lowIn: lowIn, highIn: highIn}
	if Retain_Contract.Assume(execution).IsFalse() {
		panic("")
	}
	lowOut, highOut = wrap(lowIn, highIn)
	execution.lowOut = lowOut
	execution.highOut = highOut
	if Retain_Contract.Guarantee(execution).IsFalse() {
		panic("")
	}
	return lowOut, highOut
}

type Abs_ExecutionModel struct {
	input int
	ret0  int
}

var Abs_Contract sopher.AGHyperContract[Abs_ExecutionModel] = sopher.NewAGHyperContract([]sopher.HyperAssertion[Abs_ExecutionModel]{}, []sopher.HyperAssertion[Abs_ExecutionModel]{
	sopher.NewUniversalHyperAssertion[Abs_ExecutionModel](0, 1, sopher.NewPredicateHyperAssertion(func(assignments []Abs_ExecutionModel) bool { e0 := assignments[0]; _ = e0; return e0.ret0 >= 0 })),
	sopher.NewUniversalHyperAssertion[Abs_ExecutionModel](0, 2, sopher.NewPredicateHyperAssertion(func(assignments []Abs_ExecutionModel) bool {
		e0, e1 := assignments[0], assignments[1]
		_, _ = e0, e1
		return !(e0.input >= e1.input) || (e0.ret0 >= e1.ret0)
	}))})

// guarantee: forall e0. e0.ret0 >= 0
// guarantee: forall e0 e1. !(e0.input >= e1.input) || (e0.ret0 >= e1.ret0)
func Abs(input int) int {
	wrap := func(input int) int {
		if input < 0 {
			return -input
		}
		return input
	}
	execution := Abs_ExecutionModel{input: input}
	if Abs_Contract.Assume(execution).IsFalse() {
		panic("")
	}
	ret0 := wrap(input)
	execution.ret0 = ret0
	if Abs_Contract.Guarantee(execution).IsFalse() {
		panic("")
	}
	return ret0
}
