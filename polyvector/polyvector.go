package polyvector

import (
	"cyber.ee/muzosh/pqdevkit"
	"cyber.ee/muzosh/pqdevkit/polyproxy"
	"github.com/tuneinsight/lattigo/v5/ring"
)

type PolyVector []polyproxy.PolyProxy

var mainRing = pqdevkit.MainRing

func NewRandomPolyVector(length int) PolyVector {
	vec := make([]polyproxy.PolyProxy, length)
	for i := 0; i < len(vec); i++ {
		vec[i] = polyproxy.NewRandomPoly()
	}
	return vec
}

func (vec PolyVector) Listize() []uint64 {
	listizedVec := make([]uint64, 0)
	for _, poly := range vec {
		listizedPoly := poly.Listize()
		listizedVec = append(listizedVec, listizedPoly...)
	}
	return listizedVec
}

func (vec PolyVector) InfiniteNorm() uint64 {
	max := uint64(0)
	for _, poly := range vec {
		maxPoly := poly.InfiniteNorm()
		if maxPoly > max {
			max = maxPoly
		}
	}
	return max
}

func (vec PolyVector) Scale(poly polyproxy.PolyProxy) PolyVector {
	newVec := make([]polyproxy.PolyProxy, len(vec))
	for i, v := range vec {
		newVec[i] = polyproxy.PolyProxy(mainRing.NewPoly())
		pqdevkit.MainRing.MulCoeffsMontgomery(ring.Poly(v), ring.Poly(poly), ring.Poly(newVec[i]))
	}
	return newVec
}
