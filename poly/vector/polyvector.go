package vector

import (
	"log"
	"strings"

	"cyber.ee/muzosh/pq/devkit"
	"cyber.ee/muzosh/pq/devkit/poly"
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

func (vec PolyVector) TransformedToPolyQVector() PolyQVector {
	ret := make(PolyQVector, vec.Length())
	for i, currentPoly := range vec {
		ret[i] = currentPoly.TransformedToPolyQ()
	}
	return ret
}

func (vec PolyVector) Length() int {
	return len(vec)
}

func (vec PolyVector) Listize() []int64 {
	listizedVec := make([]int64, 0, vec.Length()*devkit.MainRing.N())
	for _, currentPoly := range vec {
		listizedVec = append(listizedVec, currentPoly.Listize()...)
	}
	return listizedVec
}

// func (vec PolyVector) InfiniteNorm() uint64 {
// 	max := uint64(0)
// 	for _, poly := range vec {
// 		maxPoly := poly.InfiniteNorm()
// 		if maxPoly > max {
// 			max = maxPoly
// 		}
// 	}
// 	return max
// }

func (vec PolyVector) ScaleByPolyProxy(inputPolyProxy poly.PolyProxy) PolyProxyVector {
	switch input := inputPolyProxy.(type) {
	case poly.Poly:
		ret := make(PolyVector, len(vec))
		for i, currentPoly := range vec {
			ret[i] = currentPoly.Mul(input).(poly.Poly)
		}
		return ret
	case poly.PolyQ:
		currentPolyQVector := vec.TransformedToPolyQVector()

		newVec := make(PolyQVector, vec.Length())

		for i, polyQ := range currentPolyQVector {
			newVec[i] = polyQ.Mul(input).(poly.PolyQ)
		}
		return newVec
	default:
		log.Panic("Invalid PolyProxyVector.")
		return nil
	}
}

func (vec PolyVector) ScaleByInt(input int64) PolyProxyVector {
	newVec := make(PolyVector, vec.Length())
	for i, currentPoly := range vec {
		newVec[i] = currentPoly.ScaleByInt(input).(poly.Poly)
	}

	return newVec
}

func (vec PolyVector) Add(inputPolyProxyVector PolyProxyVector) PolyProxyVector {
	if inputPolyProxyVector.Length() != vec.Length() {
		log.Panic("Add: two vectors don't have the same length.")
	}

	switch input := inputPolyProxyVector.(type) {
	case PolyVector:
		ret := make(PolyVector, len(vec))
		for i, currentPoly := range vec {
			ret[i] = currentPoly.Add(input[i]).(poly.Poly)
		}
		return ret
	case PolyQVector:
		currentPolyQVector := vec.TransformedToPolyQVector()

		newVec := make(PolyQVector, vec.Length())

		for i, polyQ := range currentPolyQVector {
			newVec[i] = polyQ.Add(input[i]).(poly.PolyQ)
		}
		return newVec
	default:
		log.Panic("Invalid PolyProxyVector.")
		return nil
	}
}

func (vec PolyVector) Sub(inputPolyProxyVector PolyProxyVector) PolyProxyVector {
	if inputPolyProxyVector.Length() != vec.Length() {
		log.Panic("Sub: two vectors don't have the same length.")
	}

	switch input := inputPolyProxyVector.(type) {
	case PolyVector:
		ret := make(PolyVector, len(vec))
		for i, currentPoly := range vec {
			ret[i] = currentPoly.Sub(input[i]).(poly.Poly)
		}
		return ret
	case PolyQVector:
		currentPolyQVector := vec.TransformedToPolyQVector()

		newVec := make(PolyQVector, vec.Length())

		for i, polyQ := range currentPolyQVector {
			newVec[i] = polyQ.Sub(input[i]).(poly.PolyQ)
		}
		return newVec
	default:
		log.Panic("Invalid PolyProxyVector.")
		return nil
	}
}

func (vec PolyVector) Concat(inputPolyProxyVector PolyProxyVector) PolyProxyVector {
	switch input := inputPolyProxyVector.(type) {
	case PolyVector:
		ret := make(PolyVector, 0, vec.Length()+input.Length())
		ret = append(ret, vec...)
		ret = append(ret, input...)
		return ret
	case PolyQVector:
		currentPolyQVector := vec.TransformedToPolyQVector()

		newVec := make(PolyQVector, 0, vec.Length()+input.Length())
		newVec = append(newVec, currentPolyQVector...)
		newVec = append(newVec, input...)
		return newVec
	default:
		log.Panic("Invalid PolyProxyVector.")
		return nil
	}
}

func (vec PolyVector) DotProduct(inputPolyProxyVector PolyProxyVector) poly.PolyProxy {
	if inputPolyProxyVector.Length() != vec.Length() {
		log.Panic("DotProduct: two vectors don't have the same length.")
	}

	switch input := inputPolyProxyVector.(type) {
	case PolyVector:
		ret := make(poly.Poly, vec[0].Length())

		for i := 0; i < len(vec); i++ {
			ret = ret.Add(vec[i].Mul(input[i])).(poly.Poly)
		}

		return ret
	case PolyQVector:
		currentPolyQVector := vec.TransformedToPolyQVector()

		newPoly := poly.NewPolyQ()

		for i := 0; i < vec.Length(); i++ {
			devkit.MainRing.NTT(*currentPolyQVector[i].Poly, *currentPolyQVector[i].Poly)
			devkit.MainRing.NTT(*input[i].Poly, *input[i].Poly)

			devkit.MainRing.MulCoeffsBarrettThenAdd(*currentPolyQVector[i].Poly, *input[i].Poly, *newPoly.Poly)

			devkit.MainRing.INTT(*currentPolyQVector[i].Poly, *currentPolyQVector[i].Poly)
			devkit.MainRing.INTT(*input[i].Poly, *input[i].Poly)
		}
		devkit.MainRing.INTT(*newPoly.Poly, *newPoly.Poly)
		return newPoly
	default:
		log.Panic("Invalid PolyProxyVector.")
		return nil
	}
}

func (vec PolyVector) Equals(other PolyProxyVector) bool {
	switch input := other.(type) {
	case PolyVector:
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
