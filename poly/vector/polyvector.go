package vector

import (
	"log"
	"strings"

	"cyber.ee/pq/latticehelper"
	"cyber.ee/pq/latticehelper/poly"
	"github.com/raszia/gotiny"
)

type PolyVector []poly.Poly

func (vec PolyVector) Serialize() []byte {
	return gotiny.MarshalCompress(&vec)
}

func DeserializePolyVector(data []byte) PolyVector {
	var vec PolyVector
	n := gotiny.UnmarshalCompress(data, &vec)
	if n == 0 {
		panic("failed to deserialize")
	}
	return vec
}

func NewPolyVectorFromCoeffs(coeffs [][]int64) PolyVector {
	vec := make(PolyVector, len(coeffs))
	for i, coeffsI := range coeffs {
		vec[i] = poly.NewPolyFromCoeffs(coeffsI...)
	}
	return vec
}

func NewZeroPolyVector(length int) PolyVector {
	vec := make(PolyVector, length)
	for i := 0; i < len(vec); i++ {
		vec[i] = poly.NewConstantPoly(0)
	}
	return vec
}

func NewRandomPolyVector(length int) PolyVector {
	vec := make(PolyVector, length)
	for i := 0; i < len(vec); i++ {
		vec[i] = poly.NewRandomPoly()
	}
	return vec
}

func (vec PolyVector) CoeffString() string {
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

func (vec PolyVector) String() string {
	var sb strings.Builder
	sb.WriteString("PolyVector{")
	for i, currentPoly := range vec {
		sb.WriteString(currentPoly.String())
		if i != len(vec)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func (vec PolyVector) Q() PolyQVector {
	ret := make(PolyQVector, vec.Length())
	for i, currentPoly := range vec {
		ret[i] = currentPoly.Q()
	}
	return ret
}

func (vec PolyVector) WithCenteredModulo() PolyVector {
	ret := make(PolyVector, vec.Length())
	for i, currentPoly := range vec {
		ret[i] = currentPoly.WithCenteredModulo()
	}
	return ret
}

func (vec *PolyVector) ApplyToEveryCoeff(f func(int64) any) {
	for _, poly := range *vec {
		poly.ApplyToEveryCoeff(f)
	}
}

func (vec PolyVector) Length() int {
	return len(vec)
}

func (vec PolyVector) Listize() []int64 {
	listizedVec := make([]int64, 0, vec.Length()*latticehelper.MainRing.N())
	for _, currentPoly := range vec {
		listizedVec = append(listizedVec, currentPoly.Listize()...)
	}
	return listizedVec
}

func (vec PolyVector) CheckNormBound(bound int64) bool {
	for _, poly := range vec {
		if poly.CheckNormBound(bound) {
			return true
		}
	}
	return false
}

func (vec PolyVector) LowBits(alpha int64) PolyVector {
	newVec := make(PolyVector, len(vec))
	for i := 0; i < len(newVec); i++ {
		newVec[i] = vec[i].LowBits(alpha)
	}
	return newVec
}

func (vec PolyVector) ScaledByPoly(inputPoly poly.Poly) PolyVector {
	ret := make(PolyVector, len(vec))
	for i, currentPoly := range vec {
		ret[i] = currentPoly.Mul(inputPoly)
	}
	return ret
}

func (vec PolyVector) ScaledByInt(input int64) PolyVector {
	newVec := make(PolyVector, vec.Length())
	for i, currentPoly := range vec {
		newVec[i] = currentPoly.ScaledByInt(input)
	}

	return newVec
}

func (vec PolyVector) Add(inputPolyVector PolyVector) PolyVector {
	if inputPolyVector.Length() != vec.Length() {
		log.Panic("Add: two vectors don't have the same length.")
	}
	ret := make(PolyVector, len(vec))
	for i, currentPoly := range vec {
		ret[i] = currentPoly.Add(inputPolyVector[i])
	}
	return ret
}

func (vec PolyVector) Sub(inputPolyVector PolyVector) PolyVector {
	if inputPolyVector.Length() != vec.Length() {
		log.Panic("Sub: two vectors don't have the same length.")
	}

	ret := make(PolyVector, len(vec))
	for i, currentPoly := range vec {
		ret[i] = currentPoly.Sub(inputPolyVector[i])
	}
	return ret
}

func (vec PolyVector) Concat(inputPolyVector PolyVector) PolyVector {
	ret := make(PolyVector, 0, vec.Length()+inputPolyVector.Length())
	ret = append(ret, vec...)
	ret = append(ret, inputPolyVector...)
	return ret
}

func (vec PolyVector) DotProduct(inputPolyVector PolyVector) poly.Poly {
	if inputPolyVector.Length() != vec.Length() {
		log.Panic("DotProduct: two vectors don't have the same length.")
	}

	ret := make(poly.Poly, vec[0].Length())

	for i := 0; i < len(vec); i++ {
		ret = ret.Add(vec[i].Mul(inputPolyVector[i]))
	}

	return ret
}

func (vec PolyVector) Equals(other PolyVector) bool {
	for i := 0; i < vec.Length(); i++ {
		if !vec[i].Equals(other[i]) {
			return false
		}
	}
	return true
}
