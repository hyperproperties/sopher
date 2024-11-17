package examples

import (
	"log"

	sopher "github.com/hyperproperties/sopher/pkg/language"
)

var Retain_Contract sopher.AGHyperContract[Retain_ExecutionModel] = sopher.NewAGHyperContract(
	[]sopher.HyperAssertion[Retain_ExecutionModel]{},
	func(execution Retain_ExecutionModel) Retain_ExecutionModel {
		// Wrap the body of the function.
		wrap := func(a, b int) (x, y int) {
			x = a
			y = b + a
			return
		}

		x, y := wrap(execution.a, execution.b)
		execution.x = x
		execution.y = y

		return execution
	},
	[]sopher.HyperAssertion[Retain_ExecutionModel]{},
)

// assume: ...
// guarantee: ...
func Foo(a, b int) (x, y int) {
	// Construct the execution model without return values.
	execution := Retain_ExecutionModel{
		a: a,
		b: b,
	}

	// Check the assumption against the assumed model.
	if Retain_Contract.Assume(execution).IsFalse() {
		log.Println("Assumption failed")
	}

	execution = Retain_Contract.Call(execution)

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
