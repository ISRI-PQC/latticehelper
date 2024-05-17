package poly

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"cyber.ee/muzosh/pq/devkit"
	"github.com/raszia/gotiny"
)

type Poly []int64

func NewPolyFromCoeffs(coeffs ...int64) Poly {
	l := len(coeffs)

	ret := make(Poly, l)
	copy(ret, coeffs)

	return ret
}

func NewPoly() Poly {
	ret := make(Poly, devkit.MainRing.N())
	return ret
}

func NewConstantPoly(constant int64) Poly {
	ret := make(Poly, devkit.MainRing.N())
	ret[0] = constant
	return ret
}

func NewRandomPoly() Poly {
	ret := make(Poly, devkit.MainRing.N())
	for i := 0; i < len(ret); i++ {
		ret[i] = rand.Int63() >> 8
		if chance := rand.Float32(); chance < 0.5 {
			ret[i] *= -1
		}
	}
	return ret
}

func (coeffs Poly) Serialize() []byte {
	return gotiny.MarshalCompress(&coeffs)
}

func DeserializePoly(data []byte) Poly {
	var p Poly
	n := gotiny.UnmarshalCompress(data, &p)
	if n == 0 {
		panic("failed to deserialize")
	}
	return p
}

func (coeffs Poly) CoeffString() string {
	return strings.Replace(fmt.Sprint(coeffs.Listize()), " ", ",", -1)
}

func (coeffs Poly) String() string {
	ret := make([]string, 0, len(coeffs))

	if containsOnlyZeroes[int64](coeffs) {
		return "0"
	}

	for i, coeff := range coeffs {
		if coeff != 0 {
			if i == 0 {
				ret = append(ret, strconv.FormatInt(coeff, 10))
			} else if i == 1 {
				if coeff == 1 {
					ret = append(ret, "x")
				} else {
					ret = append(ret, strconv.FormatInt(coeff, 10)+"*x")
				}
			} else {
				if coeff == 1 {
					ret = append(ret, "x^"+strconv.Itoa(i))
				} else {
					ret = append(ret, strconv.FormatInt(coeff, 10)+"*x^"+strconv.Itoa(i))
				}
			}
		}
	}

	return strings.Join(ret, " + ")
}

func (coeffs Poly) TransformedToPolyQ() PolyQ {
	ret := NewPolyQ()
	newCoeffs := make([]*big.Int, devkit.MainRing.N())

	for i, coeff := range coeffs {
		newCoeffs[i] = big.NewInt(coeff)
	}

	devkit.MainRing.SetCoefficientsBigint(newCoeffs, *ret.Poly)

	return ret
}

func (coeffs Poly) Length() int {
	return len(coeffs)
}

func (coeffs Poly) Listize() []int64 {
	return coeffs
}

func (coeffs Poly) Neg() PolyProxy {
	ret := make(Poly, len(coeffs))
	for i, coeff := range coeffs {
		ret[i] = -coeff
	}
	return ret
}

func (coeffs Poly) Add(inputPolyProxy PolyProxy) PolyProxy {
	switch input := inputPolyProxy.(type) {
	case Poly:
		ret := make(Poly, devkit.MainRing.N())
		for i, coeff := range coeffs {
			ret[i] = coeff + input[i]
		}
		return ret
	case PolyQ:
		return coeffs.TransformedToPolyQ().Add(input)
	default:
		log.Panic("Invalid PolyProxy")
		return nil
	}
}

func (coeffs Poly) Sub(inputPolyProxy PolyProxy) PolyProxy {
	switch input := inputPolyProxy.(type) {
	case Poly:
		ret := make(Poly, devkit.MainRing.N())
		for i, coeff := range coeffs {
			ret[i] = coeff - input[i]
		}
		return ret
	case PolyQ:
		return coeffs.TransformedToPolyQ().Sub(input)
	default:
		log.Panic("Invalid PolyProxy")
		return nil
	}
}

func (coeffs Poly) Mul(inputPolyProxy PolyProxy) PolyProxy {
	switch input := inputPolyProxy.(type) {
	case Poly:
		return Poly(schoolbookMultiplication(coeffs, input))
	case PolyQ:
		return coeffs.TransformedToPolyQ().Mul(input)
	default:
		log.Panic("Invalid PolyProxy")
		return nil
	}
}

func (poly Poly) Pow(exp int) PolyProxy {
	if exp < 0 {
		log.Panic("Pow: Negative powers are not supported for elements of a PolyQ")
	}

	g := NewConstantPoly(1)

	for exp > 0 {
		if exp%2 == 1 {
			g = g.Mul(poly).(Poly)
		}

		poly = poly.Mul(poly).(Poly)
		exp = int(devkit.FloorDivision(exp, 2))
	}

	return g
}

func (coeffs Poly) ScaleByInt(scalar int64) PolyProxy {
	ret := make(Poly, devkit.MainRing.N())
	for i, coeff := range coeffs {
		ret[i] = coeff * scalar
	}
	return ret
}

func (coeffs Poly) AddToFirstCoeff(input int64) PolyProxy {
	ret := coeffs
	ret[0] += input
	return ret
}

func (coeffs Poly) Equals(other PolyProxy) bool {
	switch other.(type) {
	case Poly:
		return reflect.DeepEqual(coeffs, other)
	default:
		return false
	}
}
