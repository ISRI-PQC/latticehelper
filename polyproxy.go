package devkit

import (
	"strconv"
	"strings"

	"github.com/tuneinsight/lattigo/v5/ring"
)

type PolyProxy struct {
	Poly  *ring.Poly
	IsNTT bool
}

func NewConstantPolyProxy(constant uint64) PolyProxy {
	ret := MainRing.NewPoly()
	ret.Coeffs[0][0] = constant

	return PolyProxy{&ret, false}
}

func NewRandomPoly() PolyProxy {
	ret := MainUniformSampler.ReadNew()
	return PolyProxy{&ret, false}
}

func (poly PolyProxy) String() string {
	coeffs := poly.Poly.Coeffs[0]
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

	if poly.IsNTT {
		ret = append(ret, " (NTT form)")
	}

	return strings.Join(ret, " + ")
}

func (poly PolyProxy) ToNTT() {
	if poly.IsNTT {
		return
	}
	MainRing.NTT(*poly.Poly, *poly.Poly)
	poly.IsNTT = true
}

func (poly PolyProxy) FromNTT() {
	if !poly.IsNTT {
		return
	}
	MainRing.INTT(*poly.Poly, *poly.Poly)
	poly.IsNTT = false
}

func (poly PolyProxy) InfiniteNorm() uint64 {
	max := int64(0)
	for _, coeff := range poly.Listize() {
		centered_coeff := centeredModulo(int64(coeff), MainRing.Modulus().Int64())

		// We need absolute value
		if centered_coeff < 0 {
			centered_coeff = -centered_coeff
		}

		if centered_coeff > max {
			max = centered_coeff
		}
	}
	return uint64(max)
}

func (poly PolyProxy) Listize() []uint64 {
	return poly.Poly.Coeffs[0]
}

func (poly PolyProxy) Neg() PolyProxy {
	ret_poly := MainRing.NewPoly()
	MainRing.Neg(*poly.Poly, ret_poly)
	return PolyProxy{&ret_poly, poly.IsNTT}
}

func (poly PolyProxy) Add(input_poly PolyProxy) PolyProxy {
	if poly.IsNTT != input_poly.IsNTT {
		panic("Add: two polynomials don't have the same form.")
	}

	ret_poly := MainRing.NewPoly()
	MainRing.Add(*poly.Poly, *input_poly.Poly, ret_poly)

	return PolyProxy{&ret_poly, poly.IsNTT}
}

func (poly PolyProxy) Sub(input_poly PolyProxy) PolyProxy {
	if poly.IsNTT != input_poly.IsNTT {
		panic("Sub: two polynomials don't have the same form.")
	}

	ret_poly := MainRing.NewPoly()
	MainRing.Sub(*poly.Poly, *input_poly.Poly, ret_poly)

	return PolyProxy{&ret_poly, poly.IsNTT}
}

func (poly PolyProxy) Mul(input_poly PolyProxy) PolyProxy {
	if poly.IsNTT != input_poly.IsNTT {
		panic("Mul: two polynomials don't have the same form.")
	}

	was_ntt := poly.IsNTT

	if !was_ntt {
		poly.ToNTT()
		input_poly.ToNTT()
	}

	ret_poly := MainRing.NewPoly()
	MainRing.MulCoeffsBarrett(*poly.Poly, *input_poly.Poly, ret_poly)

	if !was_ntt {
		MainRing.INTT(ret_poly, ret_poly)
	}

	return PolyProxy{&ret_poly, poly.IsNTT}
}

func (poly PolyProxy) MulScalar(scalar uint64) PolyProxy {
	ret_poly := MainRing.NewPoly()

	MainRing.MulScalar(*poly.Poly, scalar, ret_poly)

	return PolyProxy{&ret_poly, poly.IsNTT}
}
