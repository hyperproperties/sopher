package examples

// guarantee: forall e.
//
//	e.ret == (e.attempt <= 3 &&
//		digits[0] == 0 && digits[1] == 1 && digits[2] == 2 && digits[3] == 3)
//
// guarantee: forall e0 e1. math.Abs(e0.time - e1.time) <= 0.1 * time.Second
// guarantee: forall e0 e1. e0.counter >= e1.counter; -> e0.attempt >= e1.attempt
func CheckPIN4(counter, attempt int, digits [4]int) bool {
	if attempt > 3 {
		return false
	}

	return digits[0] == 0 &&
		digits[1] == 1 &&
		digits[2] == 2 &&
		digits[3] == 3
}
