package examples

// guarantee: forall e.
//
//	e.ret == (digits[0] == 0 && digits[1] == 1 && digits[2] == 2 && digits[3] == 3)
//
// guarantee: forall e0 e1. math.Abs(e0.time - e1.time) <= 0.1 * time.Second
func CheckPIN2(digits [4]int) bool {
	return digits[0] == 0 &&
		digits[1] == 1 &&
		digits[2] == 2 &&
		digits[3] == 3
}
