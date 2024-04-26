package polyvector

import (
	"cyber.ee/muzosh/pqdevkit"
	"cyber.ee/muzosh/pqdevkit/polyproxy"
)

type PolyVector struct {
	PolyProxies []polyproxy.PolyProxy
	IsNTT       bool
}

var mainRing = pqdevkit.MainRing

func NewRandomPolyVector(length int) PolyVector {
	vec := make([]polyproxy.PolyProxy, length)
	for i := 0; i < len(vec); i++ {
		vec[i] = polyproxy.NewRandomPoly()
	}
	return PolyVector{vec, false}
}

func (vec *PolyVector) ToNTT() {
	if vec.IsNTT {
		return
	}
	for _, poly := range vec.PolyProxies {
		poly.ToNTT()
	}
}

func (vec *PolyVector) FromNTT() {
	if !vec.IsNTT {
		return
	}
	for _, poly := range vec.PolyProxies {
		poly.FromNTT()
	}
}

func (vec PolyVector) Listize() []uint64 {
	if vec.IsNTT {
		vec.FromNTT()
	}

	listizedVec := make([]uint64, 0)
	for _, poly := range vec.PolyProxies {
		listizedPoly := poly.Listize()
		listizedVec = append(listizedVec, listizedPoly...)
	}
	return listizedVec
}

func (vec PolyVector) InfiniteNorm() uint64 {
	if vec.IsNTT {
		vec.FromNTT()
	}

	max := uint64(0)
	for _, poly := range vec.PolyProxies {
		maxPoly := poly.InfiniteNorm()
		if maxPoly > max {
			max = maxPoly
		}
	}
	return max
}

func (vec PolyVector) Scale(input_poly polyproxy.PolyProxy) PolyVector {
	newVec := make([]polyproxy.PolyProxy, len(vec.PolyProxies))
	for i, poly := range vec.PolyProxies {
		newVec[i] = poly.Mul(input_poly)
	}

	return PolyVector{newVec, vec.IsNTT}
}

func (vec PolyVector) ScaleScalar(input uint64) PolyVector {
	newVec := make([]polyproxy.PolyProxy, len(vec.PolyProxies))
	for i, poly := range vec.PolyProxies {
		newVec[i] = poly.MulScalar(input)
	}

	return PolyVector{newVec, vec.IsNTT}
}

func (vec PolyVector) Add(input_vector PolyVector) PolyVector {
	if len(vec.PolyProxies) != len(input_vector.PolyProxies) {
		panic("Add: length of input vector is not equal to length of vector")
	}

	newVec := make([]polyproxy.PolyProxy, len(vec.PolyProxies))
	for i, poly := range vec.PolyProxies {
		newVec[i] = poly.Add(input_vector.PolyProxies[i])
	}
	return PolyVector{newVec, vec.IsNTT && input_vector.IsNTT}
}

func (vec PolyVector) Sub(input_vector PolyVector) PolyVector {
	if len(vec.PolyProxies) != len(input_vector.PolyProxies) {
		panic("Add: length of input vector is not equal to length of vector")
	}

	newVec := make([]polyproxy.PolyProxy, len(vec.PolyProxies))
	for i, poly := range vec.PolyProxies {
		newVec[i] = poly.Sub(input_vector.PolyProxies[i])
	}
	return PolyVector{newVec, vec.IsNTT && input_vector.IsNTT}
}

func (vec PolyVector) Concat(input_vector PolyVector) PolyVector {
	if vec.IsNTT != input_vector.IsNTT {
		panic("Concat: vectors don't have the same (non)NTT form")
	}

	newVec := make([]polyproxy.PolyProxy, len(vec.PolyProxies)+len(input_vector.PolyProxies))
	for _, poly := range vec.PolyProxies {
		newVec = append(newVec, poly)
	}
	for _, poly := range input_vector.PolyProxies {
		newVec = append(newVec, poly)
	}
	return PolyVector{newVec, vec.IsNTT}
}

func (vec PolyVector) DotProduct(input_vector PolyVector) polyproxy.PolyProxy {
	if len(vec.PolyProxies) != len(input_vector.PolyProxies) {
		panic("DotProduct: length of input vector is not equal to length of vector")
	}

	if vec.IsNTT != input_vector.IsNTT {
		panic("Concat: vectors don't have the same (non)NTT form")
	}

	res := mainRing.NewPoly()

	for i := 0; i < len(vec.PolyProxies); i++ {
		mainRing.MulCoeffsBarrettThenAdd(res, vec.PolyProxies[i].Poly, input_vector.PolyProxies[i].Poly)
	}

	return polyproxy.PolyProxy{Poly: res, IsNTT: vec.IsNTT}
}
