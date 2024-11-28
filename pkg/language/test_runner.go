package language

import (
	"iter"
	"testing"

	"github.com/hyperproperties/sopher/pkg/quick"
)

type TestRunner[T any] struct {
	contract AGHyperContract[T]
	call func(input T) (output T)
	n int
}

func NewTestRunner[T any](call func(input T) (output T), contract AGHyperContract[T]) TestRunner[T] {
	return TestRunner[T]{
		contract: contract,
		call: call,
	}
}

func (runner *TestRunner[T]) N(n int) *TestRunner[T] {
	runner.n = n
	return runner
}

func (runner TestRunner[T]) Run(t *testing.T) {
	tester := NewIncrementalTester[T]()
	next, _ := iter.Pull2(
		tester.Test(
			quick.Iterator[T](),
			runner.call,
			runner.contract,
		),
	)
	for idx := 0; idx < runner.n; idx++ {
		next()
	}
}