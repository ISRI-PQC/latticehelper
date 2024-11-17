package vector

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"strings"

	"github.com/isri-pqc/latticehelper"
	"github.com/isri-pqc/latticehelper/poly"
	"github.com/raszia/gotiny"
	"github.com/tuneinsight/lattigo/v5/ring"
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

// Make sure sampler is not used concurrently. If needed, created new with latticehelper.GetSampler()
// If sampler is nil, default one will be used
func NewRandomPolyQVector(sampler *ring.UniformSampler, length int) PolyQVector {
	vec := make(PolyQVector, length)
	for i := 0; i < len(vec); i++ {
		vec[i] = poly.NewRandomPolyQ(sampler)
	}
	return vec
}

func NewRandomPolyQVectorWithMaxInfNorm(length int, maxInfNorm int64) PolyQVector {
	return NewRandomPolyQVectorWithMaxInfNormWithSeed(nil, length, maxInfNorm)
}

// Input nil seed to use random seed, otherwise, only first 32 bytes from seed will be used!
func NewRandomPolyQVectorWithMaxInfNormWithSeed(seed []byte, length int, maxInfNorm int64) PolyQVector {
	vec := make(PolyQVector, length)
	for i := 0; i < len(vec); i++ {
		vec[i] = poly.NewRandomPolyQWithMaxInfNorm(seed, maxInfNorm)
	}
	return vec
}

func (vec PolyQVector) Power2Round(d int64) (PolyQVector, PolyQVector) {
	r1polys := make(PolyQVector, vec.Length())
	r0polys := make(PolyQVector, vec.Length())

	for i, poly := range vec {
		p1, p0 := poly.Power2Round(d)

		r1polys[i] = p1
		r0polys[i] = p0
	}

	return r1polys, r0polys
}

func (vec *PolyQVector) ApplyToEveryCoeff(f func(int64) any) {
	for _, poly := range *vec {
		poly.ApplyToEveryCoeff(f)
	}
}

func (vec PolyQVector) HighBits(alpha int64) PolyQVector {
	newVec := make(PolyQVector, len(vec))
	for i := 0; i < len(newVec); i++ {
		newVec[i] = vec[i].HighBits(alpha)
	}
	return newVec
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

func (vec PolyQVector) NonQ() PolyVector {
	ret := make(PolyVector, vec.Length())
	for i, currentPoly := range vec {
		ret[i] = currentPoly.NonQ()
	}
	return ret
}

func (vec PolyQVector) Length() int {
	return len(vec)
}

func (vec PolyQVector) Listize() []int64 {
	listizedVec := make([]int64, 0, vec.Length()*latticehelper.MainRing.N())
	for _, currentPoly := range vec {
		listizedVec = append(listizedVec, currentPoly.Listize()...)
	}
	return listizedVec
}

func (vec PolyQVector) InfiniteNorm() int64 {
	max := int64(0)
	for _, currentPoly := range vec {
		maxPoly := currentPoly.InfiniteNorm()
		if maxPoly > max {
			max = maxPoly
		}
	}
	return max
}

func (vec PolyQVector) SecondNorm() float64 {
	sum := int64(0)
	for _, currentPoly := range vec {
		sum += latticehelper.Pow(currentPoly.InfiniteNorm(), 2)
	}
	return math.Sqrt(float64(sum))
}

func (vec PolyQVector) ScaledByPolyQ(inputPoly poly.PolyQ) PolyQVector {
	newVec := make(PolyQVector, vec.Length())
	for i, currentPoly := range vec {
		newVec[i] = currentPoly.Mul(inputPoly)
	}

	return newVec
}

func (vec PolyQVector) ScaledByInt(input int64) PolyQVector {
	newVec := make(PolyQVector, vec.Length())
	for i, currentPoly := range vec {
		newVec[i] = currentPoly.ScaledByInt(input)
	}

	return newVec
}

func (vec PolyQVector) Add(inputPolyQVector PolyQVector) PolyQVector {
	if vec.Length() != inputPolyQVector.Length() {
		log.Panic("Add: length of input vector is not equal to length of vector")
	}

	newVec := make(PolyQVector, vec.Length())
	for i, currentPoly := range vec {
		newVec[i] = currentPoly.Add(inputPolyQVector[i])
	}
	return newVec
}

func (vec PolyQVector) Sub(inputPolyQVector PolyQVector) PolyQVector {
	if vec.Length() != inputPolyQVector.Length() {
		log.Panic("Sub: length of input vector is not equal to length of vector")
	}

	newVec := make(PolyQVector, vec.Length())
	for i, currentPoly := range vec {
		newVec[i] = currentPoly.Sub(inputPolyQVector[i])
	}
	return newVec
}

func (vec PolyQVector) Concat(inputPolyQVector PolyQVector) PolyQVector {
	newVec := make(PolyQVector, 0, vec.Length()+inputPolyQVector.Length())

	newVec = append(newVec, vec...)
	newVec = append(newVec, inputPolyQVector...)

	return newVec
}

func (vec PolyQVector) DotProduct(inputPolyQVector PolyQVector) poly.PolyQ {

	if inputPolyQVector.Length() != vec.Length() {
		log.Panic("DotProduct: two vectors don't have the same length.")
	}

	newPoly := poly.NewPolyQ()

	for i := 0; i < vec.Length(); i++ {
		latticehelper.MainRing.NTT(vec[i].Poly, vec[i].Poly)
		latticehelper.MainRing.NTT(inputPolyQVector[i].Poly, inputPolyQVector[i].Poly)

		latticehelper.MainRing.MulCoeffsBarrettThenAdd(vec[i].Poly, inputPolyQVector[i].Poly, newPoly.Poly)

		latticehelper.MainRing.INTT(vec[i].Poly, vec[i].Poly)
		latticehelper.MainRing.INTT(inputPolyQVector[i].Poly, inputPolyQVector[i].Poly)
	}
	latticehelper.MainRing.INTT(newPoly.Poly, newPoly.Poly)

	return newPoly
}

func (vec PolyQVector) Equals(other PolyQVector) bool {
	for i := 0; i < vec.Length(); i++ {
		if !vec[i].Equals(other[i]) {
			return false
		}
	}
	return true
}
