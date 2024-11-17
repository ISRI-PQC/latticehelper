package matrix

import (
	"log"
	"strings"

	"github.com/isri-pqc/latticehelper"
	"github.com/isri-pqc/latticehelper/poly"
	"github.com/isri-pqc/latticehelper/poly/vector"
	"github.com/raszia/gotiny"
)

type PolyMatrix []vector.PolyVector

func (mat PolyMatrix) Serialize() []byte {
	return gotiny.MarshalCompress(&mat)
}

func DeserializePolyMatrix(data []byte) PolyMatrix {
	var mat PolyMatrix
	n := gotiny.UnmarshalCompress(data, &mat)
	if n == 0 {
		panic("failed to deserialize")
	}
	return mat
}

func NewPolyMatrixFromCoeffs(coeffMat [][][]int64) PolyMatrix {
	newMatrix := make(PolyMatrix, len(coeffMat))
	for i := range coeffMat {
		newMatrix[i] = vector.NewPolyVectorFromCoeffs(coeffMat[i])
	}
	return PolyMatrix(newMatrix)
}

func NewRandomPolyMatrix(rows, cols int) PolyMatrix {
	polyVectors := make(PolyMatrix, rows)
	for i := 0; i < rows; i++ {
		polyVectors[i] = vector.NewRandomPolyVector(cols)
	}
	return polyVectors
}

func NewIdentityPolyMatrix(size int) PolyMatrix {
	newMatrix := NewZeroPolyMatrix(size, size)
	for i := 0; i < size; i++ {
		newMatrix[i][i][0] = int64(1)
	}
	return newMatrix
}

func NewZeroPolyMatrix(rows, cols int) PolyMatrix {
	polyVectors := make(PolyMatrix, rows)
	for i := 0; i < rows; i++ {
		polyVectors[i] = vector.NewZeroPolyVector(cols)
	}
	return polyVectors
}

func (mat PolyMatrix) CoeffString() string {
	var sb strings.Builder

	sb.WriteString("[")
	for i, polyVector := range mat {
		sb.WriteString(polyVector.CoeffString())
		if i != len(mat)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (mat PolyMatrix) String() string {
	var sb strings.Builder

	sb.WriteString("PolyMatrix{\n")
	for _, polyVector := range mat {
		sb.WriteString("\t" + polyVector.String() + "\n")
	}
	sb.WriteString("}")
	return sb.String()
}

func (mat PolyMatrix) Q() PolyQMatrix {
	polyQMatrix := make(PolyQMatrix, mat.Rows())
	for i, polyVector := range mat {
		polyQMatrix[i] = polyVector.Q()
	}
	return polyQMatrix
}

func (mat PolyMatrix) Rows() int {
	return len(mat)
}

func (mat PolyMatrix) Cols() int {
	return mat[0].Length()
}

func (mat PolyMatrix) Listize() []int64 {
	listizedVec := make([]int64, 0, mat.Rows()*mat.Cols()*latticehelper.MainRing.N())

	for _, polyQVec := range mat {
		listizedVec = append(listizedVec, polyQVec.Listize()...)
	}

	return listizedVec
}

func (mat PolyMatrix) CheckNormBound(bound int64) bool {
	for _, vec := range mat {
		if vec.CheckNormBound(bound) {
			return true
		}
	}
	return false
}
func (mat PolyMatrix) LowBits(alpha int64) PolyMatrix {
	newVec := make(PolyMatrix, mat.Rows())
	for i := 0; i < mat.Rows(); i++ {
		newVec[i] = mat[i].LowBits(alpha)
	}
	return newVec
}

func (mat PolyMatrix) Transposed() PolyMatrix {
	cols := mat.Cols()
	rows := mat.Rows()

	result := make(PolyMatrix, cols)

	for i := 0; i < cols; i++ {
		polyVector := make(vector.PolyVector, rows)

		for j := 0; j < rows; j++ {
			polyVector[j] = mat[j][i]
		}

		result[i] = polyVector
	}

	return result
}

func (mat PolyMatrix) ScaledByPoly(inputPoly poly.Poly) PolyMatrix {
	ret := make(PolyMatrix, mat.Rows())
	for i, polyVector := range mat {
		ret[i] = polyVector.ScaledByPoly(inputPoly)
	}
	return PolyMatrix(ret)
}

func (mat PolyMatrix) ScaledByInt(input int64) PolyMatrix {
	result := make(PolyMatrix, mat.Rows())

	for i, polyVector := range mat {
		result[i] = polyVector.ScaledByInt(input)
	}

	return PolyMatrix(result)
}

func (mat PolyMatrix) Add(inputPolyMatrix PolyMatrix) PolyMatrix {
	if mat.Cols() != inputPolyMatrix.Cols() || mat.Rows() != inputPolyMatrix.Rows() {
		log.Panic("Add: rows and cols of matrices are not equal")
	}
	ret := make(PolyMatrix, mat.Rows())
	for i, polyVec := range mat {
		ret[i] = polyVec.Add(inputPolyMatrix[i])
	}
	return PolyMatrix(ret)
}

func (mat PolyMatrix) Sub(inputPolyMatrix PolyMatrix) PolyMatrix {
	if mat.Cols() != inputPolyMatrix.Cols() || mat.Rows() != inputPolyMatrix.Rows() {
		log.Panic("Sub: rows and cols of matrices are not equal")
	}

	ret := make(PolyMatrix, mat.Rows())
	for i, polyVec := range mat {
		ret[i] = polyVec.Sub(inputPolyMatrix[i])
	}
	return PolyMatrix(ret)
}

func (mat PolyMatrix) Concat(inputPolyMatrix PolyMatrix) PolyMatrix {
	if mat.Rows() != inputPolyMatrix.Rows() {
		log.Panic("Concat: rows of matrices are not equal")
	}

	ret := make(PolyMatrix, mat.Rows())
	for i, polyVec := range mat {
		ret[i] = polyVec.Concat(inputPolyMatrix[i])
	}
	return PolyMatrix(ret)
}

func (mat PolyMatrix) BlockCombine(inputPolyMatrix PolyMatrix) PolyMatrix {
	if mat.Cols() != inputPolyMatrix.Cols() {
		log.Panic("BlockCombine: cols of matrices are not equal")
	}

	ret := make(PolyMatrix, 0, mat.Rows()+inputPolyMatrix.Rows())
	ret = append(ret, mat...)
	ret = append(ret, inputPolyMatrix...)
	return PolyMatrix(ret)
}

func (mat PolyMatrix) MatMul(inputPolyMatrix PolyMatrix) PolyMatrix {
	if mat.Cols() != inputPolyMatrix.Rows() {
		log.Panic("MatMul: Number of cols in first mat is not equal to number of rows in second mat")
	}

	rows, cols := mat.Rows(), mat.Cols()
	otherCols := inputPolyMatrix.Cols()

	newMat := make(PolyMatrix, rows)

	for i := 0; i < rows; i++ {
		currentVec := make(vector.PolyVector, otherCols)

		for j := 0; j < otherCols; j++ {
			currentPoly := poly.NewPoly()

			for k := 0; k < cols; k++ {
				currentPoly = currentPoly.Add(mat[i][k].Mul(inputPolyMatrix[k][j]))
			}

			currentVec[j] = currentPoly
		}
		newMat[i] = currentVec
	}

	return newMat
}

func (mat PolyMatrix) VecMul(inputPolyVector vector.PolyVector) vector.PolyVector {
	if inputPolyVector.Length() != mat.Cols() {
		log.Panic("VecMul: vectors don't have the same length")
	}

	ret := make(vector.PolyVector, mat.Rows())

	for i := 0; i < len(ret); i++ {
		currentPoly := make(poly.Poly, mat[0][0].Length())

		for j := 0; j < inputPolyVector.Length(); j++ {
			currentPoly = currentPoly.Add(inputPolyVector[j].Mul(mat[i][j]))
		}

		ret[i] = currentPoly
	}
	return ret
}

func (mat PolyMatrix) Equals(other PolyMatrix) bool {
	for i := 0; i < mat.Rows(); i++ {
		if !mat[i].Equals(other[i]) {
			return false
		}
	}
	return true
}
