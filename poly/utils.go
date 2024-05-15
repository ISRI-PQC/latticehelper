package poly

import "cyber.ee/muzosh/pq/devkit"

/*
	Returns x mod q, but centered around 0

Args:

	x (int): number to be modded
	q (int): modulus

Returns:

	int: x mod q, centered around 0
*/
func centeredModulo(x int64, q uint64) int64 {
	ret := int64(devkit.PositiveMod(x, q))
	if ret > (int64(q) >> 1) {
		ret -= int64(q)
	}
	return ret
}

func containsOnlyZeroes[V uint64 | int64](a []V) bool {
	for _, v := range a {
		if v != 0 {
			return false
		}
	}
	return true
}

func schoolbookMultiplication(p1, p2 []int64) []int64 {
	if len(p1) != len(p2) {
		panic("schoolbookMultiplication: p1 and p2 must be of the same length")
	}
	n := devkit.MainRing.N()

	newCoeffs := make([]int64, n)

	for i := 0; i < n; i++ {
		for j := 0; j < (n - i); j++ {
			newCoeffs[i+j] += (p1[i] * p2[j])
		}
	}

	for j := 1; j < n; j++ {
		for i := (n - j); i < n; i++ {
			newCoeffs[i+j-n] -= (p1[i] * p2[j])
		}
	}

	return newCoeffs
}
