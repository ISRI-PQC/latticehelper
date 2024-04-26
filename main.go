// package main

// import (
// 	"fmt"
// 	"reflect"
// 	"time"

// 	"github.com/tuneinsight/lattigo/v5/ring"
// )

// func schoolbookMultiplication(p1, p2 ring.Poly) []uint64 {
// 	n := MainRing.N()
// 	a := make([]int64, n)
// 	b := make([]int64, n)

// 	// fmt.Println("p1: ", p1)
// 	// fmt.Println("p2: ", p2)

// 	for i, coeff := range p1.Coeffs[0] {
// 		a[i] = int64(coeff)
// 	}

// 	for i, coeff := range p2.Coeffs[0] {
// 		b[i] = int64(coeff)
// 	}

// 	new_coeffs := make([]int64, n)
// 	ret := make([]uint64, n)

// 	for i := 0; i < n; i++ {
// 		for j := 0; j < (n - i); j++ {
// 			new_coeffs[i+j] += (a[i] * b[j])
// 		}
// 	}

// 	for j := 1; j < n; j++ {
// 		for i := (n - j); i < n; i++ {
// 			new_coeffs[i+j-n] -= (a[i] * b[j])
// 		}
// 	}

// 	for i := range new_coeffs {
// 		new_coeffs[i] %= MainRing.Modulus().Int64()

// 		if new_coeffs[i] < 0 {
// 			new_coeffs[i] += MainRing.Modulus().Int64()
// 		}

// 		ret[i] = uint64(new_coeffs[i])
// 	}

// 	// fmt.Println("p1 * p2: ", new_coeffs)
// 	return ret
// }

// func addPolynomials(p1, p2 ring.Poly) []uint64 {
// 	n := MainRing.N()
// 	a := make([]uint64, n)
// 	b := make([]uint64, n)

// 	// fmt.Println("p1: ", p1)
// 	// fmt.Println("p2: ", p2)

// 	for i, coeff := range p1.Coeffs[0] {
// 		a[i] = uint64(coeff)
// 	}

// 	for i, coeff := range p2.Coeffs[0] {
// 		b[i] = uint64(coeff)
// 	}

// 	ret := make([]uint64, n)
// 	for i := range a {
// 		ret[i] = (a[i] + b[i]) % MainRing.Modulus().Uint64()
// 	}
// 	return ret
// }

// func main() {
// 	fmt.Println("Hello, World!")

// 	err := Init(256, 8380417)
// 	if err!= nil {
// 		panic(err)
// 	}

// 	a := MainUniformSampler.ReadNew()
// 	b := MainUniformSampler.ReadNew()
// 	c := MainRing.NewPoly()
// 	// d := MainRing.NewPoly()

// 	start := time.Now()
// 	sb := schoolbookMultiplication(a, b)
// 	elapsed := time.Since(start)
// 	fmt.Println("schoolbookMultiplication: ", elapsed)

// 	MainRing.NTTLazy(a, c)
// 	// MainRing.NTTLazy(b, d)


// 	start = time.Now()
// 	p6 := MainRing.NewPoly()
// 	MainRing.MulCoeffsBarrett(a, b, p6)
// 	elapsed = time.Since(start)
// 	fmt.Println("non-NTT Barrett: ", elapsed)

// 	MainRing.NTT(a, a)
// 	MainRing.NTT(b, b)

// 	start = time.Now()
// 	p3 := MainRing.NewPoly()
// 	MainRing.MulCoeffsBarrett(a, b, p3)
// 	MainRing.INTT(p3, p3)
// 	elapsed = time.Since(start)
// 	fmt.Println("NTT Barrett: ", elapsed)

// 	start = time.Now()
// 	p4 := MainRing.NewPoly()
// 	MainRing.MulCoeffsMontgomeryLazy(c, b, p4)
// 	MainRing.INTT(p4, p4)
// 	elapsed = time.Since(start)
// 	fmt.Println("NTT Montogmery: ", elapsed)


// 	same := reflect.DeepEqual(sb, p3.Coeffs[0])
// 	fmt.Println(same)

// 	same = reflect.DeepEqual(sb, p4.Coeffs[0])
// 	fmt.Println(same)

// 	// ADDING
// 	a = MainUniformSampler.ReadNew()
// 	b = MainUniformSampler.ReadNew()

// 	start = time.Now()
// 	MainRing.NTT(a, a)
// 	fmt.Println("NTT: ", time.Since(start))

// 	start = time.Now()
// 	MainRing.INTT(a, a)
// 	fmt.Println("INTT: ", time.Since(start))

// 	start = time.Now()
// 	normal := addPolynomials(a, b)
// 	fmt.Println("normal add: ", time.Since(start))


// 	l := MainRing.NewPoly()
// 	start = time.Now()
// 	MainRing.Add(a, b, l)
// 	fmt.Println("lattigo normal add: ", time.Since(start))

// 	same = reflect.DeepEqual(normal, l.Coeffs[0])
// 	fmt.Println(same)

// 	MainRing.NTT(a, a)
// 	MainRing.NTT(b, b)
// 	l2 := MainRing.NewPoly()
// 	start = time.Now()
// 	MainRing.Add(a, b, l2)
// 	MainRing.INTT(l2, l2)
// 	fmt.Println("lattigo NTT add: ", time.Since(start))

// 	same = reflect.DeepEqual(normal, l2.Coeffs[0])
// 	fmt.Println(same)
	
// 	l3 := MainRing.NewPoly()
// 	start = time.Now()
// 	MainRing.INTT(a, a)
// 	MainRing.INTT(b, b)
// 	MainRing.Add(a, b, l3)
// 	fmt.Println("lattigo INTT add: ", time.Since(start))

// 	same = reflect.DeepEqual(normal, l3.Coeffs[0])
// 	fmt.Println(same)
// }
