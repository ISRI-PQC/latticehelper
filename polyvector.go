package devkit

import "strings"

type PolyVector struct {
	PolyProxies []PolyProxy
	IsNTT       bool
}

func NewZeroPolyVector(length int) PolyVector {
	vec := make([]PolyProxy, length)
	for i := 0; i < len(vec); i++ {
		vec[i] = NewConstantPolyProxy(0)
	}
	return PolyVector{vec, false}
}

func NewRandomPolyVector(length int) PolyVector {
	vec := make([]PolyProxy, length)
	for i := 0; i < len(vec); i++ {
		vec[i] = NewRandomPoly()
	}
	return PolyVector{vec, false}
}

func (vec PolyVector) String() string {
	var sb strings.Builder
	sb.WriteString("PolyVector{")
	for i, poly := range vec.PolyProxies {
		sb.WriteString(poly.String())
		if i != len(vec.PolyProxies)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func (vec PolyVector) Length() int {
	return len(vec.PolyProxies)
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

	listizedVec := make([]uint64, vec.Length()*MainRing.N())
	for _, poly := range vec.PolyProxies {
		listizedVec = append(listizedVec, poly.Listize()...)
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

func (vec PolyVector) Scale(input_poly PolyProxy) PolyVector {
	newVec := make([]PolyProxy, vec.Length())
	for i, poly := range vec.PolyProxies {
		newVec[i] = poly.Mul(input_poly)
	}

	return PolyVector{newVec, vec.IsNTT}
}

func (vec PolyVector) ScaleScalar(input uint64) PolyVector {
	newVec := make([]PolyProxy, vec.Length())
	for i, poly := range vec.PolyProxies {
		newVec[i] = poly.MulScalar(input)
	}

	return PolyVector{newVec, vec.IsNTT}
}

func (vec PolyVector) Add(input_vector PolyVector) PolyVector {
	if vec.Length() != input_vector.Length() {
		panic("Add: length of input vector is not equal to length of vector")
	}

	newVec := make([]PolyProxy, vec.Length())
	for i, poly := range vec.PolyProxies {
		newVec[i] = poly.Add(input_vector.PolyProxies[i])
	}
	return PolyVector{newVec, vec.IsNTT}
}

func (vec PolyVector) Sub(input_vector PolyVector) PolyVector {
	if vec.Length() != input_vector.Length() {
		panic("Sub: length of input vector is not equal to length of vector")
	}

	newVec := make([]PolyProxy, vec.Length())
	for i, poly := range vec.PolyProxies {
		newVec[i] = poly.Sub(input_vector.PolyProxies[i])
	}
	return PolyVector{newVec, vec.IsNTT}
}

func (vec PolyVector) Concat(input_vector PolyVector) PolyVector {
	if vec.IsNTT != input_vector.IsNTT {
		panic("Concat: vectors don't have the same form")
	}

	newVec := make([]PolyProxy, vec.Length()+input_vector.Length())

	newVec = append(newVec, vec.PolyProxies...)
	newVec = append(newVec, input_vector.PolyProxies...)

	return PolyVector{newVec, vec.IsNTT}
}

func (vec PolyVector) DotProduct(input_vector PolyVector) PolyProxy {
	if vec.Length() != input_vector.Length() {
		panic("DotProduct: length of input vector is not equal to length of vector")
	}

	if vec.IsNTT != input_vector.IsNTT {
		panic("DotProduct: vectors don't have the same (non)NTT form")
	}

	was_ntt := vec.IsNTT

	if !was_ntt {
		vec.ToNTT()
		input_vector.ToNTT()
	}

	newPoly := MainRing.NewPoly()

	for i := 0; i < vec.Length(); i++ {
		MainRing.MulCoeffsBarrettThenAdd(*vec.PolyProxies[i].Poly, *input_vector.PolyProxies[i].Poly, newPoly)
	}

	if !was_ntt {
		MainRing.INTT(newPoly, newPoly)
	}

	return PolyProxy{Poly: &newPoly, IsNTT: was_ntt}
}

func (vec PolyVector) MatMul(input_mat PolyMatrix) PolyVector {
	if vec.Length() != input_mat.Cols() {
		panic("MatMul: vectors don't have the same length")
	}

	if vec.IsNTT != input_mat.IsNTT {
		panic("Concat: vector and matrix don't have the same form")
	}

	was_ntt := vec.IsNTT

	if !was_ntt {
		vec.ToNTT()
		input_mat.ToNTT()
	}

	newVec := make([]PolyProxy, input_mat.Rows())

	for i := 0; i < input_mat.Rows(); i++ {
		currentPoly := MainRing.NewPoly()

		for j := 0; j < vec.Length(); j++ {
			MainRing.MulCoeffsBarrettThenAdd(*vec.PolyProxies[j].Poly, *input_mat.PolyVectors[i].PolyProxies[j].Poly, currentPoly)
		}

		newVec[i] = PolyProxy{&currentPoly, true}
	}

	if !was_ntt {
		vec.FromNTT()
		input_mat.FromNTT()
	}

	return PolyVector{newVec, was_ntt}
}
