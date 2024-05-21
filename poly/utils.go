package poly

import (
	"log"

	"cyber.ee/pq/devkit"
)

func CenteredModulo(x, q int64) int64 {
	ret := devkit.PositiveMod(x, q)
	if ret > (q >> 1) {
		ret -= q
	}
	return ret
}

func checkNormBound(n, b, q int64) bool {
	x := devkit.PositiveMod(n, q)
	x = ((q - 1) >> 1) - x
	x = x ^ (x >> 31)
	x = ((q - 1) >> 1) - x
	return x >= b
}

func decompose(r, a, q int64) (int64, int64) {
	r = devkit.PositiveMod(r, q)
	r0 := CenteredModulo(r, a)
	r1 := r - r0
	if r1 == q-1 {
		return 0, r0 - 1
	}

	r1 = devkit.FloorDivision(r1, a)
	if r != r1*a+r0 {
		panic("r!= r1*a+r0")
	}
	return r1, r0
}

func highBits(r, a, q int64) int64 {
	r1, _ := decompose(r, a, q)
	return r1
}

func lowBits(r, a, q int64) int64 {
	_, r0 := decompose(r, a, q)
	return r0
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
		log.Panic("schoolbookMultiplication: p1 and p2 must be of the same length")
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
