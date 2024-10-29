package poly

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"cyber.ee/pq/devkit"
	"github.com/raszia/gotiny"
	"github.com/tuneinsight/lattigo/v5/utils/sampling"
)

type Poly []int64

func NewPolyFromCoeffs(coeffs ...int64) Poly {
	l := devkit.MainRing.N()

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
		ret[i] = int64(sampling.RandUint64()) >> 8
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

	if containsOnlyZeroes(coeffs) {
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

	devkit.MainRing.SetCoefficientsBigint(newCoeffs, ret.Poly)

	return ret
}

func (coeffs Poly) WithCenteredModulo() Poly {
	ret := make([]int64, len(coeffs))
	for i, coeff := range coeffs {
		ret[i] = CenteredModulo(coeff, devkit.MainRing.Modulus().Int64())
	}
	return ret
}

func (coeffs *Poly) ApplyToEveryCoeff(f func(int64) any) {
	for i, coeff := range *coeffs {
		c := f(coeff)
		switch t := c.(type) {
		case uint64:
			(*coeffs)[i] = int64(t)
		case int64:
			(*coeffs)[i] = t
		}
	}
}

func (coeffs Poly) CheckNormBound(bound int64) bool {
	for _, coeff := range coeffs {
		if checkNormBound(coeff, bound, devkit.MainRing.Modulus().Int64()) {
			return true
		}
	}
	return false
}

func (coeffs Poly) LowBits(alpha int64) Poly {
	ret := make(Poly, len(coeffs))

	for i, coeff := range coeffs {
		ret[i] = lowBits(coeff, alpha, devkit.MainRing.Modulus().Int64())
	}

	return ret
}

func (coeffs Poly) Length() int {
	return len(coeffs)
}

func (coeffs Poly) Listize() []int64 {
	return coeffs
}

func (coeffs Poly) Neg() Poly {
	ret := make(Poly, len(coeffs))
	for i, coeff := range coeffs {
		ret[i] = -coeff
	}
	return ret
}

func (coeffs Poly) Add(inputPolynomial Polynomial) Polynomial {
	switch input := inputPolynomial.(type) {
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

func (coeffs Poly) Sub(inputPolynomial Polynomial) Polynomial {
	switch input := inputPolynomial.(type) {
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

func (coeffs Poly) Mul(inputPolynomial Polynomial) Polynomial {
	switch input := inputPolynomial.(type) {
	case Poly:
		return Poly(schoolbookMultiplication(coeffs, input))
	case PolyQ:
		return coeffs.TransformedToPolyQ().Mul(input)
	default:
		log.Panic("Invalid PolyProxy")
		return nil
	}
}

func (coeffs Poly) Pow(exp int64) Poly {
	if exp < 0 {
		log.Panic("Pow: Negative powers are not supported for elements of a PolyQ")
	}

	g := NewConstantPoly(1)

	for exp > 0 {
		if exp%2 == 1 {
			g = g.Mul(coeffs).(Poly)
		}

		coeffs = coeffs.Mul(coeffs).(Poly)
		exp = devkit.FloorDivision(exp, 2)
	}

	return g
}

func (coeffs Poly) ScaledByInt(scalar int64) Poly {
	ret := make(Poly, devkit.MainRing.N())
	for i, coeff := range coeffs {
		ret[i] = coeff * scalar
	}
	return ret
}

func (coeffs *Poly) AddedToFirstCoeff(input int64) *Poly {
	ret := coeffs
	(*ret)[0] += input
	return ret
}

func (coeffs Poly) Equals(other Poly) bool {
	return reflect.DeepEqual(coeffs, other)
}
