package main

type Asd struct {
	value int
}

// assume: forall e. !(e.low > 0 && e.high > 0;)
// guarantee: forall e0 e1. (e0.low == e1.low; -> e0.ret0 == e1.ret0;)
func (asd Asd) Foo(low, high int) int {
	if low < 0 {
		return low + asd.value
	}
	return high
}
