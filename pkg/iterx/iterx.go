package iterx

import "iter"

func Permutations(sub, slice int) iter.Seq[[]int] {
	if sub < 0 {
		panic("permutations of negative length subslices")
	}

	if sub == 0 {
		return func(yield func([]int) bool) {}
	}

	return func(yield func([]int) bool) {
		permutation := make([]int, sub)

		for {
			if !yield(permutation) {
				return
			}

			// Find the rightmost element that can still be incremented.
			i := sub - 1
			for i >= 0 && permutation[i] == slice-1 {
				i--
			}

			if i < 0 {
				break
			}

			permutation[i]++

			// Reset all elements to the right of i to zero.
			for j := i + 1; j < sub; j++ {
				permutation[j] = 0
			}
		}
	}
}

func CollectN[T any](iterator iter.Seq[T], n int) (slice []T) {
	for element := range iterator {
		if n == 0 {
			break
		}
		slice = append(slice, element)
		n--
	}
	return
}

func Collect[T any](iterator iter.Seq[T]) (slice []T) {
	for element := range iterator {
		slice = append(slice, element)
	}
	return
}
