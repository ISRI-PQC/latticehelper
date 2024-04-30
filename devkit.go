package devkit

import (
	"fmt"
	"reflect"
	"time"

	"github.com/tuneinsight/lattigo/v5/ring"
	"github.com/tuneinsight/lattigo/v5/utils/sampling"
)

// If you encounter [cyber.ee/muzosh/pq/common.MainRing] or
// [cyber.ee/muzosh/pq/common.MainUniformSampler] being nil,
// you must initialize it first using
// [cyber.ee/muzosh/pq/common.InitSingle] or
// [cyber.ee/muzosh/pq/common.InitMultiple] functions!
var (
	MainRing           *ring.Ring
	MainUniformSampler *ring.UniformSampler
)

func InitSingle(degree int, modulus uint64) error {
	return InitMultiple(degree, []uint64{modulus})
}

func InitMultiple(degree int, moduli []uint64) error {
	r, err := ring.NewRing(degree, moduli)

	if err != nil {
		return err
	}

	MainRing = r.AtLevel(0)

	prng, err := sampling.NewPRNG()

	if err != nil {
		return err
	}

	us := ring.NewUniformSampler(prng, r)

	MainUniformSampler = us

	return nil
}

const (
	iterations = 10000
)

func main() {
	fmt.Println("Hello, World!")

	err := InitSingle(256, 8380417)
	if err != nil {
		panic(err)
	}

	test()
	benchmark()
}

func benchmark() {

	ntts := make([]time.Duration, iterations)
	nttslazy := make([]time.Duration, iterations)
	intts := make([]time.Duration, iterations)
	inttslazy := make([]time.Duration, iterations)
	sbs := make([]time.Duration, iterations)
	barretntts := make([]time.Duration, iterations)
	barretnttlazy := make([]time.Duration, iterations)

	for i := 0; i < iterations; i++ {
		// fmt.Println("Iteration: ", i+1)

		a := MainUniformSampler.ReadNew()
		b := MainUniformSampler.ReadNew()

		// Preparation
		antt := MainRing.NewPoly()
		bntt := MainRing.NewPoly()
		anttlazy := MainRing.NewPoly()
		bnttlazy := MainRing.NewPoly()
		var r ring.Poly
		var same bool

		// NTT
		start := time.Now()
		MainRing.NTT(a, antt)
		ntts[i] = time.Since(start)
		MainRing.NTT(b, bntt)

		// NTTLazy
		start = time.Now()
		MainRing.NTTLazy(a, anttlazy)
		nttslazy[i] = time.Since(start)
		MainRing.NTTLazy(b, bnttlazy)

		start = time.Now()
		sb := schoolbookMultiplication(a, b)
		sbs[i] = time.Since(start)

		// Barret with NTT
		r = MainRing.NewPoly()
		start = time.Now()
		MainRing.MulCoeffsBarrett(antt, bntt, r)
		barretntts[i] = time.Since(start)

		start = time.Now()
		MainRing.INTT(r, r)
		intts[i] = time.Since(start)

		same = reflect.DeepEqual(r.Coeffs[0], sb)
		if !same {
			panic("Barret with NTT not same")
		}

		// Barret with NTTLazy
		r = MainRing.NewPoly()
		start = time.Now()
		MainRing.MulCoeffsBarrettLazy(anttlazy, bnttlazy, r)
		barretnttlazy[i] = time.Since(start)

		start = time.Now()
		MainRing.INTTLazy(r, r)
		inttslazy[i] = time.Since(start)

		same = reflect.DeepEqual(r.Coeffs[0], sb)
		if !same {
			panic("Barret with NTT Lazy not same")
		}
	}

	fmt.Println("NTT: ", average(ntts), "µs")                 // (", ntts, ")")
	fmt.Println("NTTLazy: ", average(nttslazy), "µs")         // (", nttslazy, ")")
	fmt.Println("INTT: ", average(intts), "µs")               // (", intts, ")")
	fmt.Println("INTTLazy: ", average(inttslazy), "µs")       // (", inttslazy, ")")
	fmt.Println("SB: ", average(sbs), "µs")                   // (", sbs, ")")
	fmt.Println("Barret: ", average(barretntts), "µs")        // (", barretntts, ")")
	fmt.Println("BarretLazy: ", average(barretnttlazy), "µs") // (", barretnttlazy, ")")
}

func test() {
	a := MainUniformSampler.ReadNew()
	b := MainUniformSampler.ReadNew()

	sb := schoolbookMultiplication(a, b)

	var start time.Time
	var elapsed time.Duration
	var same bool

	// r1 := MainRing.NewPoly()
	// start := time.Now()
	// MainRing.MulCoeffsBarrett(a, b, r1)
	// elapsed := time.Since(start)
	// same := reflect.DeepEqual(sb, r1.Coeffs[0])
	// fmt.Println("non-NTT barrett: ", same, " time: ", elapsed)

	MainRing.NTTLazy(a, a)
	MainRing.NTTLazy(b, b)

	r2 := MainRing.NewPoly()
	start = time.Now()
	MainRing.MulCoeffsBarrett(a, b, r2)
	elapsed = time.Since(start)
	MainRing.INTTLazy(r2, r2)
	same = reflect.DeepEqual(sb, r2.Coeffs[0])
	fmt.Println("NTT barrett: ", same, " time: ", elapsed)

	r3 := MainRing.NewPoly()
	t1 := MainRing.NewPoly()
	// t2 := MainRing.NewPoly()

	start = time.Now()
	MainRing.MForm(a, t1)
	// MainRing.MForm(b, t2)
	elapsed = time.Since(start)
	fmt.Println("MForm: ", elapsed)

	start = time.Now()
	MainRing.MulCoeffsMontgomeryLazy(t1, b, r3)
	elapsed = time.Since(start)
	MainRing.INTTLazy(r3, r3)
	same = reflect.DeepEqual(sb, r3.Coeffs[0])
	fmt.Println("NTT montgomery: ", same, " time: ", elapsed)

	// same = reflect.DeepEqual(sb, p4.Coeffs[0])
	// fmt.Println(same)

	// // ADDING
	// a = MainUniformSampler.ReadNew()
	// b = MainUniformSampler.ReadNew()

	// MainRing.NTT(a, a)

	// MainRing.INTT(a, a)

	// normal := addPolynomials(a, b)

	// l := MainRing.NewPoly()
	// MainRing.Add(a, b, l)

	// same = reflect.DeepEqual(normal, l.Coeffs[0])
	// fmt.Println(same)

	// MainRing.NTT(a, a)
	// MainRing.NTT(b, b)
	// l2 := MainRing.NewPoly()
	// MainRing.Add(a, b, l2)
	// MainRing.INTT(l2, l2)

	// same = reflect.DeepEqual(normal, l2.Coeffs[0])
	// fmt.Println(same)

	// l3 := MainRing.NewPoly()
	// MainRing.INTT(a, a)
	// MainRing.INTT(b, b)
	// MainRing.Add(a, b, l3)

	// same = reflect.DeepEqual(normal, l3.Coeffs[0])
	// fmt.Println(same)
}

func average(values []time.Duration) float32 {
	var total time.Duration

	for _, d := range values {
		total += d
	}

	return float32(total.Microseconds()) / float32(len(values))
}

func schoolbookMultiplication(p1, p2 ring.Poly) []uint64 {
	n := MainRing.N()
	a := make([]int64, n)
	b := make([]int64, n)

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

	return ret
}

func addPolynomials(p1, p2 ring.Poly) []uint64 {
	n := MainRing.N()
	a := make([]uint64, n)
	b := make([]uint64, n)

	// fmt.Println("p1: ", p1)
	// fmt.Println("p2: ", p2)

	for i, coeff := range p1.Coeffs[0] {
		a[i] = uint64(coeff)
	}

	for i, coeff := range p2.Coeffs[0] {
		b[i] = uint64(coeff)
	}

	ret := make([]uint64, n)
	for i := range a {
		ret[i] = (a[i] + b[i]) % MainRing.Modulus().Uint64()
	}
	return ret
}
