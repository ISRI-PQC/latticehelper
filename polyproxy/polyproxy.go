package polyproxy

import (
	"cyber.ee/muzosh/pqdevkit"
	"github.com/tuneinsight/lattigo/v5/ring"
)

var mainRing = pqdevkit.MainRing
var mainUniformSampler = pqdevkit.MainUniformSampler

type PolyProxy struct {
	Poly  ring.Poly
	IsNTT bool
}

func NewRandomPoly() PolyProxy {
	return PolyProxy{mainUniformSampler.ReadNew(), false}
}

func (poly *PolyProxy) ToNTT() {
	if poly.IsNTT {
		return
	}
	mainRing.NTT(poly.Poly, poly.Poly)
	poly.IsNTT = true
}

func (poly *PolyProxy) FromNTT() {
	if !poly.IsNTT {
		return
	}
	mainRing.INTT(poly.Poly, poly.Poly)
	poly.IsNTT = false
}

func (poly PolyProxy) InfiniteNorm() uint64 {
	max := int64(0)
	for _, coeff := range poly.Listize() {
		centered_coeff := centeredModulo(int64(coeff), mainRing.Modulus().Int64())

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
	ret_poly := mainRing.NewPoly()
	mainRing.Neg(poly.Poly, ret_poly)
	return PolyProxy{ret_poly, poly.IsNTT}
}

func (poly PolyProxy) Add(input_poly PolyProxy) PolyProxy {
	ret_poly := mainRing.NewPoly()

	// If both are NTT, it's faster to add directly in NTT space
	if poly.IsNTT && input_poly.IsNTT {
		mainRing.Add(poly.Poly, input_poly.Poly, ret_poly)
		return PolyProxy{ret_poly, true}
	}

	// Cheapest is non-NTT addition
	if poly.IsNTT {
		poly.FromNTT()
	}
	if input_poly.IsNTT {
		input_poly.FromNTT()
	}

	mainRing.Add(poly.Poly, input_poly.Poly, ret_poly)
	return PolyProxy{ret_poly, false}
}

func (poly PolyProxy) Sub(input_poly PolyProxy) PolyProxy {
	ret_poly := mainRing.NewPoly()

	// If both are NTT, it's faster to sub directly in NTT space
	if poly.IsNTT && input_poly.IsNTT {
		mainRing.Sub(poly.Poly, input_poly.Poly, ret_poly)
		return PolyProxy{ret_poly, true}
	}

	// Cheapest is non-NTT addition
	if poly.IsNTT {
		poly.FromNTT()
	}
	if input_poly.IsNTT {
		input_poly.FromNTT()
	}

	mainRing.Sub(poly.Poly, input_poly.Poly, ret_poly)
	return PolyProxy{ret_poly, false}
}

func (poly PolyProxy) Mul(input_poly PolyProxy) PolyProxy {
	ret_poly := mainRing.NewPoly()

	// If both are in NTT, multiply and return as NTT
	if poly.IsNTT && input_poly.IsNTT {
		mainRing.MulCoeffsBarrett(poly.Poly, input_poly.Poly, ret_poly)
		return PolyProxy{ret_poly, true}
	}

	// If one is in NTT and the other is not, convert to NTT first then multiply, then return as non-NTT
	if !poly.IsNTT {
		poly.ToNTT()
	}

	if !input_poly.IsNTT {
		input_poly.ToNTT()
	}

	mainRing.MulCoeffsBarrett(poly.Poly, input_poly.Poly, ret_poly)

	ret := PolyProxy{ret_poly, false}
	ret.FromNTT()
	return ret
}

func (poly PolyProxy) MulScalar(scalar uint64) PolyProxy {
	ret_poly := mainRing.NewPoly()

	mainRing.MulScalar(poly.Poly, scalar, ret_poly)

	return PolyProxy{ret_poly, poly.IsNTT}
}
