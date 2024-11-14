package examples

import (
	"log"

	sopher "github.com/hyperproperties/sopher/pkg/language"
)

var Retain_Contract sopher.AGHyperContract[Retain_ExecutionModel] = sopher.NewAGHyperContract(
	[]sopher.HyperAssertion[Retain_ExecutionModel]{},
	[]sopher.HyperAssertion[Retain_ExecutionModel]{},
)

// assume: ...
// guarantee: ...
func Foo(a, b int) (x, y int) {
	// Wrap the body of the function.
	wrap := func(a, b int) (x, y int) {
		x = a
		y = b + a
		return
	}

	// Create a helper function to update the execution with the returned values.
	call := func(execution Retain_ExecutionModel) Retain_ExecutionModel {
		// Store the return values.
		x, y = wrap(a, b)

		// Assign the return values to the execution.
		execution.x = x
		execution.y = y

		return execution
	}

	// It is up to the model how that set looks like but we need a set
	// we can use as a form of ground truth about the implementation.
	// We need executions to base current execution's judgement on.
	Retain_Contract.Model(call)

	// Construct the execution model without return values.
	execution := Retain_ExecutionModel{
		a: a,
		b: b,
	}

	// Check the assumption against the assumed model.
	if Retain_Contract.Assume(execution).IsFalse() {
		log.Println("Assumption failed")
	}

	execution = call(execution)

	// Check the guarantee against the assumed model.
	if Retain_Contract.Guarantee(execution).IsFalse() {
		log.Println("Guarantee failed")
	}

	return x, y
}

type Retain_ExecutionModel struct {
	a, b int
	x, y int
}
