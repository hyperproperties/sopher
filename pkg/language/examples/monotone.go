package examples

// assume: forall e. e.value >= 0 && e.value <= 100_000
// guarantee: forall e0 e1. e0.value >= e1.value; <-> e0.ret0 >= e1.ret0
func Monotone(value float64) float64 {
	return value * 2
}
