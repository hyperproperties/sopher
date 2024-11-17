package language

import "testing"

func Test(t *testing.T) {
	/*type Execution struct {
		input int
		output int
	}

	contract := NewAGHyperContract(
		[]HyperAssertion[Execution]{
			NewUniversalHyperAssertion(0, 1, NewPredicateHyperAssertion(
				func(assignments []Execution) bool {
					e := assignments[0]
					return e.input >= 0
				},
			)),
		},
		[]HyperAssertion[Execution]{
			NewUniversalHyperAssertion(0, 2, NewPredicateHyperAssertion(
				func(assignments []Execution) bool {
					e0, e1 := assignments[0], assignments[1]
					return (e0.input >= e1.input) == (e0.output >= e1.output)
				},
			)),
		},
	)

	monotone := func(input int) int {
		return input + 1
	}

	contract.Model(func(execution Execution) Execution {
		output := monotone(execution.input)
		execution.output = output
		return execution
	})*/
}
