package matrix

import (
	"bytes"
	"encoding/binary"
	"log"
	"strings"

	"cyber.ee/pq/devkit"
	"cyber.ee/pq/devkit/poly"
	"cyber.ee/pq/devkit/poly/vector"
	"github.com/raszia/gotiny"
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

func NewRandomPolyQMatrix(rows, cols int) PolyQMatrix {
	newMatrix := make(PolyQMatrix, rows)
	for i := 0; i < rows; i++ {
		newMatrix[i] = vector.NewRandomPolyQVector(cols)
	}
	return PolyQMatrix(newMatrix)
}

func NewIdentityPolyQMatrix(size int) PolyQMatrix {
	newMatrix := NewZeroPolyQMatrix(size, size)
	for i := 0; i < size; i++ {
		newMatrix[i][i].Poly.Coeffs[0][0] = uint64(1)
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

func (mat PolyQMatrix) Rows() int {
	return len(mat)
}

func (mat PolyQMatrix) Cols() int {
	return mat[0].Length()
}

func (mat PolyQMatrix) Listize() []int64 {
	listizedVec := make([]int64, 0, mat.Rows()*mat.Cols()*devkit.MainRing.N())

	for _, polyQVec := range mat {
		listizedVec = append(listizedVec, polyQVec.Listize()...)
	}

	return listizedVec
}

func (mat PolyQMatrix) InfiniteNorm() uint64 {
	max := uint64(0)
	for _, polyQVec := range mat {
		maxVec := polyQVec.InfiniteNorm()
		if maxVec > max {
			max = maxVec
		}
	}

	return max
}

func (mat PolyQMatrix) Transposed() PolyProxyMatrix {
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

func (mat PolyQMatrix) Power2Round(d int) (PolyQMatrix, PolyQMatrix) {
	r1vecs := make(PolyQMatrix, mat.Rows())
	r0vecs := make(PolyQMatrix, mat.Rows())

	for i, vec := range mat {
		v1, v0 := vec.Power2Round(d)

		r1vecs[i] = v1
		r0vecs[i] = v0
	}

	return r1vecs, r0vecs
}

func (mat PolyQMatrix) ScaleByPolyProxy(inputPoly poly.PolyProxy) PolyProxyMatrix {
	result := make(PolyQMatrix, mat.Rows())

	for i, polyQVector := range mat {
		result[i] = polyQVector.ScaleByPolyProxy(inputPoly).(vector.PolyQVector)
	}
	return result
}

func (mat PolyQMatrix) ScaleByInt(input int64) PolyProxyMatrix {
	result := make(PolyQMatrix, mat.Rows())

	for i, polyQVector := range mat {
		result[i] = polyQVector.ScaleByInt(input).(vector.PolyQVector)
	}
	return result
}

func (mat PolyQMatrix) Add(inputPolyProxyMat PolyProxyMatrix) PolyProxyMatrix {
	if mat.Cols() != inputPolyProxyMat.Cols() || mat.Rows() != inputPolyProxyMat.Rows() {
		log.Panic("Add: rows and cols of matrices are not equal")
	}

	var inputPolyQMatrix PolyQMatrix

	switch input := inputPolyProxyMat.(type) {
	case PolyQMatrix:
		inputPolyQMatrix = input
	case PolyMatrix:
		inputPolyQMatrix = input.TransformedToPolyQMatrix()
	}

	newMat := make(PolyQMatrix, mat.Rows())
	for i, polyQVec := range mat {
		newMat[i] = polyQVec.Add(inputPolyQMatrix[i]).(vector.PolyQVector)
	}

	return newMat
}

func (mat PolyQMatrix) Sub(inputPolyProxyMat PolyProxyMatrix) PolyProxyMatrix {
	if mat.Cols() != inputPolyProxyMat.Cols() || mat.Rows() != inputPolyProxyMat.Rows() {
		log.Panic("Sub: rows and cols of matrices are not equal")
	}

	var inputPolyQMatrix PolyQMatrix

	switch input := inputPolyProxyMat.(type) {
	case PolyQMatrix:
		inputPolyQMatrix = input
	case PolyMatrix:
		inputPolyQMatrix = input.TransformedToPolyQMatrix()
	}

	newMat := make(PolyQMatrix, mat.Rows())
	for i, polyQVec := range mat {
		newMat[i] = polyQVec.Sub(inputPolyQMatrix[i]).(vector.PolyQVector)
	}

	return newMat
}

func (mat PolyQMatrix) Concat(inputPolyProxyMat PolyProxyMatrix) PolyProxyMatrix {
	if mat.Rows() != inputPolyProxyMat.Rows() {
		log.Panic("Concat: rows of matrices are not equal")
	}

	var inputPolyQMatrix PolyQMatrix

	switch input := inputPolyProxyMat.(type) {
	case PolyQMatrix:
		inputPolyQMatrix = input
	case PolyMatrix:
		inputPolyQMatrix = input.TransformedToPolyQMatrix()
	}

	newMat := make(PolyQMatrix, mat.Rows())
	for i, polyQVec := range mat {
		newMat[i] = polyQVec.Concat(inputPolyQMatrix[i]).(vector.PolyQVector)
	}

	return newMat
}

func (mat PolyQMatrix) BlockCombine(inputPolyProxyMat PolyProxyMatrix) PolyProxyMatrix {
	if mat.Cols() != inputPolyProxyMat.Cols() {
		log.Panic("BlockCombine: cols of matrices are not equal")
	}

	var inputPolyQMatrix PolyQMatrix

	switch input := inputPolyProxyMat.(type) {
	case PolyQMatrix:
		inputPolyQMatrix = input
	case PolyMatrix:
		inputPolyQMatrix = input.TransformedToPolyQMatrix()
	}

	newMat := make(PolyQMatrix, 0, mat.Rows()+inputPolyQMatrix.Cols())

	newMat = append(newMat, mat...)
	newMat = append(newMat, inputPolyQMatrix...)

	return newMat
}

func (mat PolyQMatrix) MatMul(inputPolyProxyMat PolyProxyMatrix) PolyProxyMatrix {
	var inputPolyQMatrix PolyQMatrix

	switch input := inputPolyProxyMat.(type) {
	case PolyQMatrix:
		inputPolyQMatrix = input
	case PolyMatrix:
		inputPolyQMatrix = input.TransformedToPolyQMatrix()
	}

	if mat.Cols() != inputPolyQMatrix.Rows() {
		log.Panic("MatMul: Number of cols in first mat is not equal to number of rows in second mat")
	}

	rows, cols := mat.Rows(), mat.Cols()
	otherCols := inputPolyQMatrix.Cols()
	inputTransposed := inputPolyQMatrix.Transposed().(PolyQMatrix)

	newMat := make(PolyQMatrix, rows)

	for i := 0; i < rows; i++ {
		currentVec := make(vector.PolyQVector, otherCols)

		for j := 0; j < otherCols; j++ {
			currentPoly := poly.NewPolyQ()

			for k := 0; k < cols; k++ {
				devkit.MainRing.NTT(*mat[i][k].Poly, *mat[i][k].Poly)
				devkit.MainRing.NTT(*inputTransposed[k][j].Poly, *inputTransposed[k][j].Poly)

				devkit.MainRing.MulCoeffsBarrettThenAdd(
					*mat[i][k].Poly,
					*inputTransposed[k][j].Poly,
					*currentPoly.Poly)

				devkit.MainRing.INTT(*mat[i][k].Poly, *mat[i][k].Poly)
				devkit.MainRing.INTT(*inputTransposed[k][j].Poly, *inputTransposed[k][j].Poly)
			}

			devkit.MainRing.INTT(*currentPoly.Poly, *currentPoly.Poly)

			currentVec[j] = currentPoly
		}
		newMat[i] = currentVec
	}

	return newMat
}

func (mat PolyQMatrix) VecMul(inputPolyProxyVector vector.PolyProxyVector) vector.PolyProxyVector {
	if inputPolyProxyVector.Length() != mat.Cols() {
		log.Panic("VecMul: vectors don't have the same length")
	}
	var inputPolyQVector vector.PolyQVector

	switch input := inputPolyProxyVector.(type) {
	case vector.PolyQVector:
		inputPolyQVector = input
	case vector.PolyVector:
		inputPolyQVector = input.TransformedToPolyQVector()
	}

	newVec := make(vector.PolyQVector, mat.Rows())

	for i := 0; i < mat.Rows(); i++ {
		currentPoly := poly.NewPolyQ()

		for j := 0; j < inputPolyProxyVector.Length(); j++ {
			devkit.MainRing.NTT(*inputPolyQVector[j].Poly, *inputPolyQVector[j].Poly)
			devkit.MainRing.NTT(*mat[i][j].Poly, *mat[i][j].Poly)

			devkit.MainRing.MulCoeffsBarrettThenAdd(*inputPolyQVector[j].Poly, *mat[i][j].Poly, *currentPoly.Poly)

			devkit.MainRing.INTT(*inputPolyQVector[j].Poly, *inputPolyQVector[j].Poly)
			devkit.MainRing.INTT(*mat[i][j].Poly, *mat[i][j].Poly)
		}
		devkit.MainRing.INTT(*currentPoly.Poly, *currentPoly.Poly)

		newVec[i] = currentPoly
	}

	return newVec
}

func (mat PolyQMatrix) Equals(other PolyProxyMatrix) bool {
	switch input := other.(type) {
	case PolyQMatrix:
		if mat.Rows() != input.Rows() || mat.Cols() != input.Cols() {
			return false
		}

		for i := 0; i < mat.Rows(); i++ {
			if !mat[i].Equals(input[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
