package examples

import (
	"log"

	sopher "github.com/hyperproperties/sopher/pkg/language"
)

var Retain_Contract sopher.AGHyperContract[Retain_ExecutionModel] = sopher.NewAGHyperContract(
	sopher.AllAssertion[Retain_ExecutionModel]{},
	sopher.AllAssertion[Retain_ExecutionModel]{},
)

// assume: ...
// guarantee: ...
func Foo(a, b int) (x, y int) {
	caller := sopher.Caller()

	// Construct the execution model without return values.
	execution := Retain_ExecutionModel{
		a: a,
		b: b,
	}

	call := func(execution Retain_ExecutionModel) Retain_ExecutionModel {
		// Wrap the body of the function.
		wrap := func(a, b int) (x, y int) {
			x = a
			y = b + a
			return
		}

		// Call the wrapped function and store the outputs.
		x, y := wrap(execution.a, execution.b)
		execution.x = x
		execution.y = y

		// Return the execution with outputs.
		return execution
	}

	// Execute the function under the contract's obligations.
	assumption, execution, guarantee := Retain_Contract.Call(caller, execution, call)

	// Check the assumption against the assumed model.
	if assumption.IsFalse() {
		log.Println("Assumption failed")
	}

	// Check the guarantee against the assumed model.
	if guarantee.IsFalse() {
		log.Println("Guarantee failed")
	}

	// Return the outputs stored in the execution from calling the triple.
	return execution.x, execution.y
}

type Retain_ExecutionModel struct {
	a, b int
	x, y int
}
