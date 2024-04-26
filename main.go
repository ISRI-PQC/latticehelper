package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/tuneinsight/lattigo/v5/ring"
)

func schoolbookMultiplication(p1, p2 ring.Poly) []uint64 {
	n := MainRing.N()
	a := make([]int64, n)
	b := make([]int64, n)

	fmt.Println("p1: ", p1)
	fmt.Println("p2: ", p2)

	for i, coeff := range p1.Coeffs[0] {
		a[i] = int64(coeff)
	}

	for i, coeff := range p2.Coeffs[0] {
		b[i] = int64(coeff)
	}

	new_coeffs := make([]int64, n)
	ret := make([]uint64, n)

	for i := 0; i < n; i++ {
		for j := 0; j < (n - i); j++ {
			new_coeffs[i+j] += (a[i] * b[j])
		}
	}

	for j := 1; j < n; j++ {
		for i := (n - j); i < n; i++ {
			new_coeffs[i+j-n] -= (a[i] * b[j])
		}
	}

	for i := range new_coeffs {
		new_coeffs[i] %= MainRing.Modulus().Int64()

		if new_coeffs[i] < 0 {
			new_coeffs[i] += MainRing.Modulus().Int64()
		}

		ret[i] = uint64(new_coeffs[i])
	}

	fmt.Println("p1 * p2: ", new_coeffs)
	return ret
}

func main() {
	fmt.Println("Hello, World!")

	err := Init(128, 4294954753)
	if err!= nil {
		panic(err)
	}

	a := MainUniformSampler.ReadNew()
	b := MainUniformSampler.ReadNew()

	start := time.Now()
	sb := schoolbookMultiplication(a, b)
	elapsed := time.Since(start)
	fmt.Println("schoolbookMultiplication: ", elapsed)

	start = time.Now()
	p3 := MainRing.NewPoly()
	MainRing.NTT(a, a)
	MainRing.NTT(b, b)
	MainRing.MulCoeffsBarrett(a, b, p3)
	MainRing.INTT(p3, p3)
	elapsed = time.Since(start)
	fmt.Println("NTT: ", elapsed)
	fmt.Println("p3: ", p3)

	same := reflect.DeepEqual(sb, p3.Coeffs[0])
	fmt.Println(same)
}
