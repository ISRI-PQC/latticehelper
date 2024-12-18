package poly

import (
	cr "crypto/rand"
	"fmt"
	"log"
	"math/big"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/isri-pqc/latticehelper"
	"github.com/tuneinsight/lattigo/v5/ring"
)

type PolyQ struct {
	ring.Poly
}

func NewPolyQFromCoeffs(coeffs ...int64) PolyQ {
	ret := latticehelper.MainRing.AtLevel(latticehelper.MainRing.Level()).NewPoly()

	newCoeffs := make([]*big.Int, latticehelper.MainRing.N())

	for i, coeff := range coeffs {
		newCoeffs[i] = new(big.Int).SetInt64(coeff)
	}

	for i := len(coeffs); i < latticehelper.MainRing.N(); i++ {
		newCoeffs[i] = big.NewInt(0)
	}

	latticehelper.MainRing.SetCoefficientsBigint(newCoeffs, ret)

	return PolyQ{ret}
}

func NewPolyQ() PolyQ {
	ret := latticehelper.MainRing.AtLevel(latticehelper.MainRing.Level()).NewPoly()
	return PolyQ{ret}
}

func NewConstantPolyQ(constant int64) PolyQ {
	ret := NewPolyQ()

	constant = latticehelper.PositiveMod(constant, latticehelper.MainRing.Modulus().Int64())

	ret.Coeffs[latticehelper.MainRing.Level()][0] = uint64(constant)

	return ret
}

// Make sure sampler is not used concurrently. If needed, created new with latticehelper.GetSampler()
// If sampler is nil, default one will be used
func NewRandomPolyQ(sampler *ring.UniformSampler) PolyQ {
	if sampler == nil {
		sampler = latticehelper.DefaultUniformSampler
	}

	ret := sampler.ReadNew()
	return PolyQ{ret}
}

// Input nil seed to use random seed, otherwise, only first 32 bytes from seed will be used!
func NewRandomPolyQWithMaxInfNorm(seed []byte, maxInfNorm int64) PolyQ {
	ret := latticehelper.MainRing.NewPoly()
	newCoeffs := make([]*big.Int, latticehelper.MainRing.N())

	var r *rand.Rand

	if seed == nil {
		seed32 := make([]byte, 32)
		_, err := cr.Read(seed32)
		if err != nil {
			panic(err)
		}
		seed = seed32
	}

	r = rand.New(rand.NewChaCha8([32]byte(seed)))

	for i := range newCoeffs {
		c := r.Int64N(maxInfNorm + 1)
		if r.Float64() > 0.5 {
			c = -c
		}

		newCoeffs[i] = big.NewInt(c)
	}

	latticehelper.MainRing.SetCoefficientsBigint(newCoeffs, ret)

	return PolyQ{ret}
}

func (poly PolyQ) Serialize() []byte {
	b, err := poly.Poly.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}

func DeserializePolyQ(data []byte) PolyQ {
	p := NewPolyQ()
	err := p.Poly.UnmarshalBinary(data)
	if err != nil {
		panic(err)
	}
	return p
}

func (poly PolyQ) CoeffString() string {
	return strings.Replace(fmt.Sprint(poly.Listize()), " ", ",", -1)
}

func (poly PolyQ) String() string {
	coeffs := poly.Poly.Coeffs[latticehelper.MainRing.Level()]
	ret := make([]string, 0, len(coeffs)+1)

	if containsOnlyZeroes(coeffs) {
		return "0"
	}

	for i, coeff := range coeffs {
		if coeff != 0 {
			if i == 0 {
				ret = append(ret, strconv.FormatUint(coeff, 10))
			} else if i == 1 {
				if coeff == 1 {
					ret = append(ret, "x")
				} else {
					ret = append(ret, strconv.FormatUint(coeff, 10)+"*x")
				}
			} else {
				if coeff == 1 {
					ret = append(ret, "x^"+strconv.Itoa(i))
				} else {
					ret = append(ret, strconv.FormatUint(coeff, 10)+"*x^"+strconv.Itoa(i))
				}
			}
		}
	}

	return strings.Join(ret, " + ")
}

func (coeffs PolyQ) NonQ() Poly {
	ret := NewPoly()

	for i, coeff := range coeffs.Coeffs[latticehelper.MainRing.Level()] {
		ret[i] = int64(coeff)
	}

	return ret
}

func (poly PolyQ) InfiniteNorm() int64 {
	max := int64(0)
	for _, coeff := range poly.Listize() {
		centeredCoeff := CenteredModulo(int64(coeff), latticehelper.MainRing.Modulus().Int64())

		// We need absolute value
		if centeredCoeff < 0 {
			centeredCoeff = -centeredCoeff
		}

		if centeredCoeff > max {
			max = centeredCoeff
		}
	}
	return max
}

func (poly PolyQ) Length() int {
	return latticehelper.MainRing.N()
}

func (poly PolyQ) Listize() []int64 {
	ret := make([]int64, len(poly.Poly.Coeffs[latticehelper.MainRing.Level()]))
	for i := 0; i < len(ret); i++ {
		ret[i] = int64(poly.Poly.Coeffs[latticehelper.MainRing.Level()][i])
	}
	return ret
}

func (poly *PolyQ) ApplyToEveryCoeff(f func(int64) any) {
	newCoeffs := make([]*big.Int, poly.Length())

	for i := 0; i < poly.Length(); i++ {
		c := f(int64(poly.Coeffs[latticehelper.MainRing.Level()][i]))
		switch t := c.(type) {
		case int64:
			newCoeffs[i] = new(big.Int).SetInt64(t)
		case uint64:
			newCoeffs[i] = new(big.Int).SetUint64(t)
		default:
			panic("unexpected type")
		}
	}

	latticehelper.MainRing.SetCoefficientsBigint(newCoeffs, poly.Poly)
}

func (poly PolyQ) Power2Round(d int64) (PolyQ, PolyQ) {
	r1coeffs := make([]int64, poly.Length())
	r0coeffs := make([]int64, poly.Length())

	for i, coeff := range poly.Coeffs[latticehelper.MainRing.Level()] {
		centered := CenteredModulo(int64(coeff), latticehelper.Pow(2, d))

		r1coeffs[i] = latticehelper.FloorDivision(int64(coeff)-centered, latticehelper.Pow(2, d))
		r0coeffs[i] = centered
	}

	ret1 := NewPolyQFromCoeffs(r1coeffs...)
	ret0 := NewPolyQFromCoeffs(r0coeffs...)

	return ret1, ret0
}

func (poly PolyQ) HighBits(alpha int64) PolyQ {
	ret := poly.CopyNew()

	for i, coeff := range poly.Coeffs[latticehelper.MainRing.Level()] {
		ret.Coeffs[latticehelper.MainRing.Level()][i] = uint64(highBits(int64(coeff), alpha, latticehelper.MainRing.Modulus().Int64()))
	}

	return PolyQ{*ret}
}

func (poly PolyQ) Neg() PolyQ {
	retPoly := NewPolyQ()
	latticehelper.MainRing.Neg(poly.Poly, retPoly.Poly)
	return retPoly
}

func (poly PolyQ) Add(inputPolyQ PolyQ) PolyQ {
	retPoly := NewPolyQ()
	latticehelper.MainRing.Add(poly.Poly, inputPolyQ.Poly, retPoly.Poly)
	return retPoly

}

func (poly PolyQ) Sub(inputPolyQ PolyQ) PolyQ {
	retPoly := NewPolyQ()
	latticehelper.MainRing.Sub(poly.Poly, inputPolyQ.Poly, retPoly.Poly)
	return retPoly
}

func (poly PolyQ) Mul(inputPolyQ PolyQ) PolyQ {
	r := latticehelper.MainRing.AtLevel(latticehelper.MainRing.Level())

	polyNTT := r.NewPoly()
	inputNTT := r.NewPoly()

	r.NTT(poly.Poly, polyNTT)
	r.NTT(inputPolyQ.Poly, inputNTT)

	retPoly := NewPolyQ()
	r.MulCoeffsBarrett(polyNTT, inputNTT, retPoly.Poly)

	r.INTT(retPoly.Poly, retPoly.Poly)

	return retPoly
}

func (poly PolyQ) Pow(exp int64) PolyQ {
	if exp < 0 {
		log.Panic("Pow: Negative powers are not supported for elements of a PolyQ")
	}

	g := NewConstantPolyQ(1)

	for exp > 0 {
		if exp%2 == 1 {
			g = g.Mul(poly)
		}

		poly = poly.Mul(PolyQ{*poly.CopyNew()})
		exp = latticehelper.FloorDivision(exp, 2)
	}

	return g
}

func (poly PolyQ) ScaledByInt(scalar int64) PolyQ {
	retPoly := NewPolyQ()

	sc := latticehelper.PositiveMod(scalar, latticehelper.MainRing.Modulus().Int64())

	latticehelper.MainRing.MulScalar(poly.Poly, uint64(sc), retPoly.Poly)

	return retPoly
}

func (poly PolyQ) AddedToFirstCoeff(input int64) PolyQ {
	retPoly := *poly.CopyNew()

	inputQ := latticehelper.PositiveMod(input, latticehelper.MainRing.Modulus().Int64())

	addPoly := NewConstantPolyQ(inputQ)

	latticehelper.MainRing.Add(retPoly, addPoly.Poly, retPoly)

	return PolyQ{retPoly}
}

func (poly PolyQ) Equals(other PolyQ) bool {
	return latticehelper.MainRing.Equal(poly.Poly, other.Poly)
}
