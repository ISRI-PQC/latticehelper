package polyproxy

import (
	"cyber.ee/muzosh/pqdevkit"
	"github.com/tuneinsight/lattigo/v5/ring"
)

type PolyProxy ring.Poly

func NewRandomPoly() PolyProxy {
	return PolyProxy(pqdevkit.MainUniformSampler.ReadNew())
}

func (poly PolyProxy) InfiniteNorm() uint64 {
	max := int64(0)
	for _, coeff := range poly.Listize() {
		centered_coeff := centeredModulo(int64(coeff), pqdevkit.MainRing.Modulus().Int64())

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
	return poly.Coeffs[0]
}
