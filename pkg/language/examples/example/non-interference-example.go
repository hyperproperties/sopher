package examples

import (
	"log"

	sopher "github.com/hyperproperties/sopher/pkg/language"
)

var Retain_Contract sopher.AGContract[Retain_ExecutionModel] = sopher.NewAGContract(
	[]sopher.IncrementalMonitor[Retain_ExecutionModel]{},
	[]sopher.IncrementalMonitor[Retain_ExecutionModel]{},
)

// guarantee: forall e0 e1. !(e0.highIn == e1.highIn) || (e0.lowOut == e1.lowOut)
func Retain(lowIn, highIn int) (lowOut, highOut int) {
	// Wrap the body of the function.
	wrap := func(lowIn, highIn int) (lowOut, highOut int) {
		lowOut = lowIn
		highOut = highIn + lowIn
		return
	}

	// Construct the execution model without return values.
	execution := Retain_ExecutionModel{
		lowIn:  lowIn,
		highIn: highIn,
	}

	// Check the assumption against all stored executions.
	if Retain_Contract.Assume(execution).IsFalse() {
		log.Println("Assumption failed")
	}

	// Store the return values.
	lowOut, highOut = wrap(lowIn, highIn)

	// Save the return values in the execution model.
	execution.lowOut = lowOut
	execution.highOut = highOut

	// Check the guarantee against all stored executions.
	if Retain_Contract.Guarantee(execution).IsFalse() {
		log.Println("Guarantee failed")
	}

	// Return the function's return values.
	return lowOut, highOut
}

// guarantee: forall e0. e0.ret >= 0
// guarantee: forall e0 e1. !(e0.input >= e1.input) || (e0.ret >= e1.ret)
func Abs(input int) int {
	if input < 0 {
		return -input
	}
	return input
}

type Retain_ExecutionModel struct {
	lowIn, highIn   int
	lowOut, highOut int
}
type Abs_ExecutionModel struct {
	input int
	ret0  int
}
