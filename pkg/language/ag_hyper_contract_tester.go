package language

import "iter"

type AGHyperContractTester[T any] interface {
	Test(inputs iter.Seq[T], call func(input T) (output T), contract AGHyperContract[T], model ...T)
}
