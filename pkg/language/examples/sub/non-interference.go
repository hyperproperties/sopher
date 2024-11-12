package examples

// guarantee: forall e0 e1. !(e0.highIn == e1.highIn) || (e0.lowOut == e1.lowOut)
func Retain(lowIn, highIn int) (lowOut, highOut int) {
	lowOut = lowIn
	highOut = highIn + lowIn
	return
}

// guarantee: forall e0. e0.ret >= 0
// guarantee: forall e0 e1. !(e0.input >= e1.input) || (e0.ret >= e1.ret)
func Abs(input int) int {
	if input < 0 {
		return -input
	}
	return input
}
