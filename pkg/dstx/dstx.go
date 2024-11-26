package dstx

func Construct[A ~[]B, B any, S ~[]E, E any](s S, factory func(i int, e E) B) A {
	expressions := make(A, len(s))
	for i, e := range s {
		expressions[i] = factory(i, e)
	}
	return expressions
}

func RepeatS(s string, n int) []string {
	result := make([]string, n)
	for idx := range result {
		result[idx] = s
	}
	return result
}
