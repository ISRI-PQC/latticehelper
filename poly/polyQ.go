package poly

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"cyber.ee/muzosh/pq/devkit"
	"github.com/tuneinsight/lattigo/v5/ring"
)

type PolyQ struct {
	*ring.Poly
}

func NewPolyQFromCoeffs(coeffs ...int64) PolyQ {
	ret := devkit.MainRing.NewPoly()

	newCoeffs := make([]*big.Int, devkit.MainRing.N())

	for i, coeff := range coeffs {
		newCoeffs[i] = new(big.Int).SetInt64(coeff)
	}

	for i := len(coeffs); i < devkit.MainRing.N(); i++ {
		newCoeffs[i] = big.NewInt(0)
	}

	devkit.MainRing.SetCoefficientsBigint(newCoeffs, ret)

	return PolyQ{&ret}
}

func NewPolyQ() PolyQ {
	ret := devkit.MainRing.NewPoly()
	return PolyQ{&ret}
}

func NewConstantPolyQ(constant uint64) PolyQ {
	ret := NewPolyQ()
	ret.Coeffs[0][0] = constant

	return ret
}

func NewRandomPolyQ() PolyQ {
	ret := devkit.MainUniformSampler.ReadNew()
	return PolyQ{&ret}
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
	coeffs := poly.Poly.Coeffs[0]
	ret := make([]string, 0, len(coeffs)+1)

	if containsOnlyZeroes[uint64](coeffs) {
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

func (poly PolyQ) InfiniteNorm() uint64 {
	max := int64(0)
	for _, coeff := range poly.Listize() {
		centeredCoeff := centeredModulo(int64(coeff), devkit.MainRing.Modulus().Uint64())

		// We need absolute value
		if centeredCoeff < 0 {
			centeredCoeff = -centeredCoeff
		}

		if centeredCoeff > max {
			max = centeredCoeff
		}
	}
	return uint64(max)
}

func (poly PolyQ) Length() int {
	return devkit.MainRing.N()
}

func (poly PolyQ) Listize() []int64 {
	ret := make([]int64, len(poly.Poly.Coeffs[0]))
	for i := 0; i < len(ret); i++ {
		ret[i] = int64(poly.Poly.Coeffs[0][i])
	}
	return ret
}

func (poly PolyQ) Neg() PolyProxy {
	retPoly := NewPolyQ()
	devkit.MainRing.Neg(*poly.Poly, *retPoly.Poly)
	return retPoly
}

func (poly PolyQ) Add(inputPolyProxy PolyProxy) PolyProxy {
	var inputPolyQ PolyQ

	switch input := inputPolyProxy.(type) {
	case PolyQ:
		inputPolyQ = input
	case Poly:
		inputPolyQ = input.TransformedToPolyQ()
	}

	retPoly := NewPolyQ()
	devkit.MainRing.Add(*poly.Poly, *inputPolyQ.Poly, *retPoly.Poly)

	return retPoly

}

func (poly PolyQ) Sub(inputPolyProxy PolyProxy) PolyProxy {
	var inputPolyQ PolyQ

	switch input := inputPolyProxy.(type) {
	case PolyQ:
		inputPolyQ = input
	case Poly:
		inputPolyQ = input.TransformedToPolyQ()
	}

	retPoly := NewPolyQ()
	devkit.MainRing.Sub(*poly.Poly, *inputPolyQ.Poly, *retPoly.Poly)

	return retPoly
}

func (poly PolyQ) Mul(inputPolyProxy PolyProxy) PolyProxy {
	var inputPolyQ PolyQ

	switch input := inputPolyProxy.(type) {
	case PolyQ:
		inputPolyQ = input
	case Poly:
		inputPolyQ = input.TransformedToPolyQ()
	}

	devkit.MainRing.NTT(*poly.Poly, *poly.Poly)
	devkit.MainRing.NTT(*inputPolyQ.Poly, *inputPolyQ.Poly)

	retPoly := NewPolyQ()
	devkit.MainRing.MulCoeffsBarrett(*poly.Poly, *inputPolyQ.Poly, *retPoly.Poly)

	devkit.MainRing.INTT(*poly.Poly, *poly.Poly)
	devkit.MainRing.INTT(*inputPolyQ.Poly, *inputPolyQ.Poly)
	devkit.MainRing.INTT(*retPoly.Poly, *retPoly.Poly)

	return retPoly
}

func (poly PolyQ) Pow(exp int) PolyProxy {
	if exp < 0 {
		log.Panic("Pow: Negative powers are not supported for elements of a PolyQ")
	}

	g := NewConstantPolyQ(1)

	for exp > 0 {
		if exp%2 == 1 {
			g = g.Mul(poly).(PolyQ)
		}

		poly = poly.Mul(PolyQ{poly.CopyNew()}).(PolyQ)
		exp = int(devkit.FloorDivision(exp, 2))
	}

	return g
}

func (poly PolyQ) ScaleByInt(scalar int64) PolyProxy {
	retPoly := NewPolyQ()

	sc := devkit.PositiveMod(scalar, devkit.MainRing.Modulus().Uint64())

	devkit.MainRing.MulScalar(*poly.Poly, sc, *retPoly.Poly)

	return retPoly
}

func (poly PolyQ) AddToFirstCoeff(input int64) PolyProxy {
	retPoly := poly.CopyNew()

	inputQ := devkit.PositiveMod(input, devkit.MainRing.Modulus().Uint64())

	addPoly := NewConstantPolyQ(inputQ)

	devkit.MainRing.Add(*retPoly, *addPoly.Poly, *retPoly)

	return PolyQ{retPoly}
}

func (poly PolyQ) Equals(other PolyProxy) bool {
	switch input := other.(type) {
	case PolyQ:
		return devkit.MainRing.Equal(*poly.Poly, *input.Poly)
	default:
		return false
	}
}
