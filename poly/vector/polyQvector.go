package vector

import (
	"bytes"
	"encoding/binary"
	"log"
	"strings"

	"cyber.ee/pq/devkit"
	"cyber.ee/pq/devkit/poly"
	"github.com/raszia/gotiny"
)

type PolyQVector []poly.PolyQ

func (vec PolyQVector) Serialize() []byte {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.LittleEndian, uint16(vec.Length()))
	if err != nil {
		panic(err)
	}
	_, err = buf.Write(gotiny.MarshalCompress(&vec))
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func DeserializePolyQVector(data []byte) PolyQVector {
	var len uint16
	_ = binary.Read(bytes.NewReader(data[:2]), binary.LittleEndian, &len)

	p := NewZeroPolyQVector(int(len))
	n := gotiny.UnmarshalCompress(data[2:], &p)
	if n == 0 {
		panic("failed to deserialize PolyQVector")
	}

	return p
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
		vec[i] = poly.NewPolyQ()
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

func NewRandomPolyQVectorWithMaxInfNorm(length int, maxInfNorm int) PolyQVector {
	vec := make(PolyQVector, length)
	for i := 0; i < len(vec); i++ {
		vec[i] = poly.NewRandomPolyQWithMaxInfNorm(maxInfNorm)
	}
	return vec
}

func (vec PolyQVector) Power2Round(d int) (PolyQVector, PolyQVector) {
	r1polys := make(PolyQVector, vec.Length())
	r0polys := make(PolyQVector, vec.Length())

	for i, poly := range vec {
		p1, p0 := poly.Power2Round(d)

		r1polys[i] = p1
		r0polys[i] = p0
	}

	return r1polys, r0polys
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
		log.Panic("Add: length of input vector is not equal to length of vector")
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
		log.Panic("Sub: length of input vector is not equal to length of vector")
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
		log.Panic("DotProduct: two vectors don't have the same length.")
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
