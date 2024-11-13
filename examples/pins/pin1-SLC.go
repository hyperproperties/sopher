package examples

// guarantee: forall e. e.ret == (len(digits) == 4 &&
//
//	digits[0] == 0 && digits[1] == 1 && digits[2] == 2 && digits[3] == 3)
//
// guarantee: forall e0 e1. math.Abs(e0.time - e1.time) <= 0.1 * time.Second
// For all pairs if one is correct the other with different digits is incorrect.
// gurantee forall e0 e1. e0.ret && !slices.Equal(e0.digits, e1.digits); -> !e1.ret
func CheckPIN1(digits ...int) bool {
	return len(digits) == 4 &&
		digits[0] == 0 &&
		digits[1] == 1 &&
		digits[2] == 2 &&
		digits[3] == 3
}
