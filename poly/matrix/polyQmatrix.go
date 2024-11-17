package matrix

import (
	"bytes"
	"encoding/binary"
	"log"
	"strings"

	"github.com/isri-pqc/latticehelper"
	"github.com/isri-pqc/latticehelper/poly"
	"github.com/isri-pqc/latticehelper/poly/vector"
	"github.com/raszia/gotiny"
	"github.com/tuneinsight/lattigo/v5/ring"
)

type PolyQMatrix []vector.PolyQVector

func (mat PolyQMatrix) Serialize() []byte {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.LittleEndian, uint16(mat.Rows()))
	if err != nil {
		panic(err)
	}
	err = binary.Write(&buf, binary.LittleEndian, uint16(mat.Cols()))
	if err != nil {
		panic(err)
	}
	_, err = buf.Write(gotiny.MarshalCompress(&mat))
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func DeserializePolyQMatrix(data []byte) PolyQMatrix {
	var rows, cols uint16
	_ = binary.Read(bytes.NewReader(data[:2]), binary.LittleEndian, &rows)
	_ = binary.Read(bytes.NewReader(data[2:4]), binary.LittleEndian, &cols)

	p := NewZeroPolyQMatrix(int(rows), int(cols))
	n := gotiny.UnmarshalCompress(data[4:], &p)
	if n == 0 {
		panic("failed to deserialize PolyQVector")
	}

	return p
}

func NewPolyQMatrixFromCoeffs(coeffMat [][][]int64) PolyQMatrix {
	newMatrix := make(PolyQMatrix, len(coeffMat))
	for i := range coeffMat {
		newMatrix[i] = vector.NewPolyQVectorFromCoeffs(coeffMat[i])
	}
	return PolyQMatrix(newMatrix)
}

// Make sure sampler is not used concurrently. If needed, created new with latticehelper.GetSampler()
// If sampler is nil, default one will be used
func NewRandomPolyQMatrix(sampler *ring.UniformSampler, rows, cols int) PolyQMatrix {
	newMatrix := make(PolyQMatrix, rows)
	for i := 0; i < rows; i++ {
		newMatrix[i] = vector.NewRandomPolyQVector(sampler, cols)
	}
	return PolyQMatrix(newMatrix)
}

func NewIdentityPolyQMatrix(size int) PolyQMatrix {
	newMatrix := NewZeroPolyQMatrix(size, size)
	for i := 0; i < size; i++ {
		newMatrix[i][i].Poly.Coeffs[latticehelper.MainRing.Level()][0] = uint64(1)
	}
	return newMatrix
}

func NewZeroPolyQMatrix(rows, cols int) PolyQMatrix {
	newMatrix := make(PolyQMatrix, rows)
	for i := 0; i < rows; i++ {
		newMatrix[i] = vector.NewZeroPolyQVector(cols)
	}
	return PolyQMatrix(newMatrix)
}

func (mat PolyQMatrix) CoeffString() string {
	var sb strings.Builder

	sb.WriteString("[")
	for i, polyQVector := range mat {
		sb.WriteString(polyQVector.CoeffString())
		if i != len(mat)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (mat PolyQMatrix) String() string {
	var sb strings.Builder

	sb.WriteString("PolyQMatrix{\n")
	for _, polyQVector := range mat {
		sb.WriteString("\t" + polyQVector.String() + "\n")
	}
	sb.WriteString("}")
	return sb.String()
}

func (mat PolyQMatrix) NonQ() PolyMatrix {
	polyMatrix := make(PolyMatrix, mat.Rows())
	for i, polyQVector := range mat {
		polyMatrix[i] = polyQVector.NonQ()
	}
	return polyMatrix
}

func (mat PolyQMatrix) Rows() int {
	return len(mat)
}

func (mat PolyQMatrix) Cols() int {
	return mat[0].Length()
}

func (mat PolyQMatrix) Listize() []int64 {
	listizedVec := make([]int64, 0, mat.Rows()*mat.Cols()*latticehelper.MainRing.N())

	for _, polyQVec := range mat {
		listizedVec = append(listizedVec, polyQVec.Listize()...)
	}

	return listizedVec
}

func (mat PolyQMatrix) InfiniteNorm() int64 {
	max := int64(0)
	for _, polyQVec := range mat {
		maxVec := polyQVec.InfiniteNorm()
		if maxVec > max {
			max = maxVec
		}
	}

	return max
}

func (mat PolyQMatrix) Transposed() PolyQMatrix {
	cols := mat.Cols()
	rows := mat.Rows()

	result := make(PolyQMatrix, cols)

	for i := 0; i < cols; i++ {
		polyQVector := make(vector.PolyQVector, rows)

		for j := 0; j < rows; j++ {
			polyQVector[j] = mat[j][i]
		}

		result[i] = polyQVector
	}

	return result
}

func (mat PolyQMatrix) Power2Round(d int64) (PolyQMatrix, PolyQMatrix) {
	r1vecs := make(PolyQMatrix, mat.Rows())
	r0vecs := make(PolyQMatrix, mat.Rows())

	for i, vec := range mat {
		v1, v0 := vec.Power2Round(d)

		r1vecs[i] = v1
		r0vecs[i] = v0
	}

	return r1vecs, r0vecs
}

func (mat PolyQMatrix) HighBits(alpha int64) PolyQMatrix {
	newVec := make(PolyQMatrix, mat.Rows())
	for i := 0; i < mat.Rows(); i++ {
		newVec[i] = mat[i].HighBits(alpha)
	}
	return newVec
}

func (mat PolyQMatrix) ScaledByPolyQ(inputPoly poly.PolyQ) PolyQMatrix {
	result := make(PolyQMatrix, mat.Rows())

	for i, polyQVector := range mat {
		result[i] = polyQVector.ScaledByPolyQ(inputPoly)
	}
	return result
}

func (mat PolyQMatrix) ScaleByInt(input int64) PolyQMatrix {
	result := make(PolyQMatrix, mat.Rows())

	for i, polyQVector := range mat {
		result[i] = polyQVector.ScaledByInt(input)
	}
	return result
}

func (mat PolyQMatrix) Add(inputPolyQMatrix PolyQMatrix) PolyQMatrix {
	if mat.Cols() != inputPolyQMatrix.Cols() || mat.Rows() != inputPolyQMatrix.Rows() {
		log.Panic("Add: rows and cols of matrices are not equal")
	}

	newMat := make(PolyQMatrix, mat.Rows())
	for i, polyQVec := range mat {
		newMat[i] = polyQVec.Add(inputPolyQMatrix[i])
	}

	return newMat
}

func (mat PolyQMatrix) Sub(inputPolyQMatrix PolyQMatrix) PolyQMatrix {
	if mat.Cols() != inputPolyQMatrix.Cols() || mat.Rows() != inputPolyQMatrix.Rows() {
		log.Panic("Sub: rows and cols of matrices are not equal")
	}

	newMat := make(PolyQMatrix, mat.Rows())
	for i, polyQVec := range mat {
		newMat[i] = polyQVec.Sub(inputPolyQMatrix[i])
	}

	return newMat
}

func (mat PolyQMatrix) Concat(inputPolyQMatrix PolyQMatrix) PolyQMatrix {
	if mat.Rows() != inputPolyQMatrix.Rows() {
		log.Panic("Concat: rows of matrices are not equal")
	}

	newMat := make(PolyQMatrix, mat.Rows())
	for i, polyQVec := range mat {
		newMat[i] = polyQVec.Concat(inputPolyQMatrix[i])
	}

	return newMat
}

func (mat PolyQMatrix) BlockCombine(inputPolyQMatrix PolyQMatrix) PolyQMatrix {
	if mat.Cols() != inputPolyQMatrix.Cols() {
		log.Panic("BlockCombine: cols of matrices are not equal")
	}

	newMat := make(PolyQMatrix, 0, mat.Rows()+inputPolyQMatrix.Cols())

	newMat = append(newMat, mat...)
	newMat = append(newMat, inputPolyQMatrix...)

	return newMat
}

func (mat PolyQMatrix) MatMul(inputPolyQMatrix PolyQMatrix) PolyQMatrix {
	if mat.Cols() != inputPolyQMatrix.Rows() {
		log.Panic("MatMul: Number of cols in first mat is not equal to number of rows in second mat")
	}

	rows, cols := mat.Rows(), mat.Cols()
	otherCols := inputPolyQMatrix.Cols()

	newMat := make(PolyQMatrix, rows)
	r := latticehelper.MainRing.AtLevel(latticehelper.MainRing.Level())

	for i := 0; i < rows; i++ {
		currentVec := make(vector.PolyQVector, otherCols)

		for j := 0; j < otherCols; j++ {
			currentPoly := poly.NewPolyQ()

			for k := 0; k < cols; k++ {
				matNTT := r.NewPoly()
				inputNTT := r.NewPoly()

				r.NTT(mat[i][k].Poly, matNTT)
				r.NTT(inputPolyQMatrix[k][j].Poly, inputNTT)

				r.MulCoeffsBarrettThenAdd(
					matNTT,
					inputNTT,
					currentPoly.Poly)
			}

			r.INTT(currentPoly.Poly, currentPoly.Poly)

			currentVec[j] = currentPoly
		}
		newMat[i] = currentVec
	}

	return newMat
}

func (mat PolyQMatrix) VecMul(inputPolyQVector vector.PolyQVector) vector.PolyQVector {
	if inputPolyQVector.Length() != mat.Cols() {
		log.Panic("VecMul: vectors don't have the same length")
	}
	newVec := make(vector.PolyQVector, mat.Rows())

	r := latticehelper.MainRing.AtLevel(latticehelper.MainRing.Level())
	for i := 0; i < mat.Rows(); i++ {
		currentPoly := poly.NewPolyQ()

		for j := 0; j < inputPolyQVector.Length(); j++ {
			matNTT := r.NewPoly()
			inputNTT := r.NewPoly()

			r.NTT(inputPolyQVector[j].Poly, inputNTT)
			r.NTT(mat[i][j].Poly, matNTT)

			r.MulCoeffsBarrettThenAdd(inputNTT, matNTT, currentPoly.Poly)
		}
		r.INTT(currentPoly.Poly, currentPoly.Poly)

		newVec[i] = currentPoly
	}

	return newVec
}

func (mat PolyQMatrix) Equals(other PolyQMatrix) bool {
	if mat.Rows() != other.Rows() || mat.Cols() != other.Cols() {
		return false
	}

	for i := 0; i < mat.Rows(); i++ {
		if !mat[i].Equals(other[i]) {
			return false
		}
	}
	return true
}
