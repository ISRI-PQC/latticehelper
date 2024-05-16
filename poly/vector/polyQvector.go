package vector

import (
	"strings"

	"cyber.ee/muzosh/pq/devkit"
	"cyber.ee/muzosh/pq/devkit/poly"
)

type PolyQVector []poly.PolyQ

func (vec PolyQVector) Serialize() ([]byte, error) {
	poly.GobBuffer.Reset()
	err := poly.GobEncoder.Encode(vec)
	ret := poly.GobBuffer.Bytes()
	poly.GobBuffer.Reset()
	return ret, err
}

func DeserializePolyQVector(data []byte) (PolyQVector, error) {
	poly.GobBuffer.Reset()
	var p PolyQVector
	err := poly.GobDecoder.Decode(&p)
	poly.GobBuffer.Reset()
	return p, err
}

func NewPolyQVectorFromCoeffs(coeffs [][]int64) PolyQVector {
	vec := make(PolyQVector, len(coeffs))
	for i, coeffsI := range coeffs {
		vec[i] = poly.NewPolyQFromCoeffs(coeffsI...)
	}
	return vec
}

func NewZeroPolyQVector(length int) PolyQVector {
	vec := make(PolyQVector, length)
	for i := 0; i < len(vec); i++ {
		vec[i] = poly.NewConstantPolyQ(0)
	}
	return vec
}

func NewRandomPolyQVector(length int) PolyQVector {
	vec := make(PolyQVector, length)
	for i := 0; i < len(vec); i++ {
		vec[i] = poly.NewRandomPolyQ()
	}
	return vec
}

func (vec PolyQVector) CoeffString() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, currentPoly := range vec {
		sb.WriteString(currentPoly.CoeffString())
		if i != len(vec)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (vec PolyQVector) String() string {
	var sb strings.Builder
	sb.WriteString("PolyQVector{")
	for i, currentPoly := range vec {
		sb.WriteString(currentPoly.String())
		if i != len(vec)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func (vec PolyQVector) Length() int {
	return len(vec)
}

func (vec PolyQVector) Listize() []int64 {
	listizedVec := make([]int64, 0, vec.Length()*devkit.MainRing.N())
	for _, currentPoly := range vec {
		listizedVec = append(listizedVec, currentPoly.Listize()...)
	}
	return listizedVec
}

func (vec PolyQVector) InfiniteNorm() uint64 {
	max := uint64(0)
	for _, currentPoly := range vec {
		maxPoly := currentPoly.InfiniteNorm()
		if maxPoly > max {
			max = maxPoly
		}
	}
	return max
}

func (vec PolyQVector) ScaleByPolyProxy(inputPolyProxy poly.PolyProxy) PolyProxyVector {
	newVec := make(PolyQVector, vec.Length())
	for i, currentPoly := range vec {
		newVec[i] = currentPoly.Mul(inputPolyProxy).(poly.PolyQ)
	}

	return newVec
}

func (vec PolyQVector) ScaleByInt(input int64) PolyProxyVector {
	newVec := make(PolyQVector, vec.Length())
	for i, currentPoly := range vec {
		newVec[i] = currentPoly.ScaleByInt(input).(poly.PolyQ)
	}

	return newVec
}

func (vec PolyQVector) Add(inputPolyProxyVector PolyProxyVector) PolyProxyVector {
	var inputPolyQVector PolyQVector

	switch input := inputPolyProxyVector.(type) {
	case PolyQVector:
		inputPolyQVector = input
	case PolyVector:
		inputPolyQVector = input.TransformedToPolyQVector()
	}

	if vec.Length() != inputPolyQVector.Length() {
		panic("Add: length of input vector is not equal to length of vector")
	}

	newVec := make(PolyQVector, vec.Length())
	for i, currentPoly := range vec {
		newVec[i] = currentPoly.Add(inputPolyQVector[i]).(poly.PolyQ)
	}
	return newVec
}

func (vec PolyQVector) Sub(inputPolyProxyVector PolyProxyVector) PolyProxyVector {
	var inputPolyQVector PolyQVector

	switch input := inputPolyProxyVector.(type) {
	case PolyQVector:
		inputPolyQVector = input
	case PolyVector:
		inputPolyQVector = input.TransformedToPolyQVector()
	}

	if vec.Length() != inputPolyQVector.Length() {
		panic("Sub: length of input vector is not equal to length of vector")
	}

	newVec := make(PolyQVector, vec.Length())
	for i, currentPoly := range vec {
		newVec[i] = currentPoly.Sub(inputPolyQVector[i]).(poly.PolyQ)
	}
	return newVec
}

func (vec PolyQVector) Concat(inputPolyProxyVector PolyProxyVector) PolyProxyVector {
	var inputPolyQVector PolyQVector

	switch input := inputPolyProxyVector.(type) {
	case PolyQVector:
		inputPolyQVector = input
	case PolyVector:
		inputPolyQVector = input.TransformedToPolyQVector()
	}

	newVec := make(PolyQVector, 0, vec.Length()+inputPolyQVector.Length())

	newVec = append(newVec, vec...)
	newVec = append(newVec, inputPolyQVector...)

	return newVec
}

func (vec PolyQVector) DotProduct(inputPolyProxyVector PolyProxyVector) poly.PolyProxy {
	var inputPolyQVector PolyQVector

	switch input := inputPolyProxyVector.(type) {
	case PolyQVector:
		inputPolyQVector = input
	case PolyVector:
		inputPolyQVector = input.TransformedToPolyQVector()
	}

	if inputPolyQVector.Length() != vec.Length() {
		panic("DotProduct: two vectors don't have the same length.")
	}

	newPoly := poly.NewPolyQ()

	for i := 0; i < vec.Length(); i++ {
		devkit.MainRing.NTT(*vec[i].Poly, *vec[i].Poly)
		devkit.MainRing.NTT(*inputPolyQVector[i].Poly, *inputPolyQVector[i].Poly)

		devkit.MainRing.MulCoeffsBarrettThenAdd(*vec[i].Poly, *inputPolyQVector[i].Poly, *newPoly.Poly)

		devkit.MainRing.INTT(*vec[i].Poly, *vec[i].Poly)
		devkit.MainRing.INTT(*inputPolyQVector[i].Poly, *inputPolyQVector[i].Poly)
	}
	devkit.MainRing.INTT(*newPoly.Poly, *newPoly.Poly)

	return newPoly
}

func (vec PolyQVector) Equals(other PolyProxyVector) bool {
	switch input := other.(type) {
	case PolyQVector:
		for i := 0; i < vec.Length(); i++ {
			if !vec[i].Equals(input[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
