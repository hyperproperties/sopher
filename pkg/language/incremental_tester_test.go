package language

import (
	"iter"
	"testing"

	"github.com/hyperproperties/sopher/pkg/quick"
)

func TestIncrementalTester(t *testing.T) {
	tester := NewIncrementalTester[int]()
	contract := NewAGHyperContract(
		NewTrueHyperAssertion[int](),
		NewTrueHyperAssertion[int](),
	)

	next, _ := iter.Pull2(
		tester.Test(
			quick.Iterator[int](),
			func(input int) (output int) {
				return input
			},
			contract,
		),
	)

	for idx := 0; idx < 10; idx++ {
		_, result, _ := next()
		if !result.IsTrue() {
			t.Fail()
		}
	}
}
