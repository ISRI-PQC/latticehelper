package matrix

import (
	"log"
	"strings"

	"cyber.ee/pq/devkit"
	"cyber.ee/pq/devkit/poly"
	"cyber.ee/pq/devkit/poly/vector"
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

func (mat PolyMatrix) TransformedToPolyQMatrix() PolyQMatrix {
	polyQMatrix := make(PolyQMatrix, mat.Rows())
	for i, polyVector := range mat {
		polyQMatrix[i] = polyVector.TransformedToPolyQVector()
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
	listizedVec := make([]int64, 0, mat.Rows()*mat.Cols()*devkit.MainRing.N())

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

func (mat PolyMatrix) ScaledByPolyProxy(inputPoly poly.Polynomial) PolynomialMatrix {
	switch input := inputPoly.(type) {
	case poly.Poly:
		ret := make(PolyMatrix, mat.Rows())
		for i, polyVector := range mat {
			ret[i] = polyVector.ScaledByPolyProxy(input).(vector.PolyVector)
		}
		return PolyMatrix(ret)
	case poly.PolyQ:
		currentPolyQMatrix := mat.TransformedToPolyQMatrix()

		newMat := make(PolyQMatrix, mat.Rows())

		for i, polyQVec := range currentPolyQMatrix {
			newMat[i] = polyQVec.ScaledByPolyProxy(input)
		}
		return newMat
	default:
		log.Panic("Invalid PolyProxyMatrix.")
		return nil
	}
}

func (mat PolyMatrix) ScaledByInt(input int64) PolyMatrix {
	result := make(PolyMatrix, mat.Rows())

	for i, polyVector := range mat {
		result[i] = polyVector.ScaledByInt(input)
	}

	return PolyMatrix(result)
}

func (mat PolyMatrix) Add(inputPolyProxyMat PolynomialMatrix) PolynomialMatrix {
	if mat.Cols() != inputPolyProxyMat.Cols() || mat.Rows() != inputPolyProxyMat.Rows() {
		log.Panic("Add: rows and cols of matrices are not equal")
	}

	switch input := inputPolyProxyMat.(type) {
	case PolyMatrix:
		ret := make(PolyMatrix, mat.Rows())
		for i, polyVec := range mat {
			ret[i] = polyVec.Add(input[i]).(vector.PolyVector)
		}
		return PolyMatrix(ret)
	case PolyQMatrix:
		currentPolyQMatrix := mat.TransformedToPolyQMatrix()

		newMat := make(PolyQMatrix, mat.Rows())

		for i, polyQVec := range currentPolyQMatrix {
			newMat[i] = polyQVec.Add(input[i])
		}
		return newMat
	default:
		log.Panic("Invalid PolyProxyMatrix.")
		return nil
	}
}

func (mat PolyMatrix) Sub(inputPolyProxyMat PolynomialMatrix) PolynomialMatrix {
	if mat.Cols() != inputPolyProxyMat.Cols() || mat.Rows() != inputPolyProxyMat.Rows() {
		log.Panic("Sub: rows and cols of matrices are not equal")
	}

	switch input := inputPolyProxyMat.(type) {
	case PolyMatrix:
		ret := make(PolyMatrix, mat.Rows())
		for i, polyVec := range mat {
			ret[i] = polyVec.Sub(input[i]).(vector.PolyVector)
		}
		return PolyMatrix(ret)
	case PolyQMatrix:
		currentPolyQMatrix := mat.TransformedToPolyQMatrix()

		newMat := make(PolyQMatrix, mat.Rows())

		for i, polyQVec := range currentPolyQMatrix {
			newMat[i] = polyQVec.Sub(input[i])
		}
		return newMat
	default:
		log.Panic("Invalid PolyProxyMatrix.")
		return nil
	}
}

func (mat PolyMatrix) Concat(inputPolyProxyMat PolynomialMatrix) PolynomialMatrix {
	if mat.Rows() != inputPolyProxyMat.Rows() {
		log.Panic("Concat: rows of matrices are not equal")
	}

	switch input := inputPolyProxyMat.(type) {
	case PolyMatrix:
		ret := make(PolyMatrix, mat.Rows())
		for i, polyVec := range mat {
			ret[i] = polyVec.Concat(input[i]).(vector.PolyVector)
		}
		return PolyMatrix(ret)
	case PolyQMatrix:
		currentPolyQMatrix := mat.TransformedToPolyQMatrix()

		newMat := make(PolyQMatrix, mat.Rows())

		for i, polyQVec := range currentPolyQMatrix {
			newMat[i] = polyQVec.Add(input[i])
		}
		return newMat
	default:
		log.Panic("Invalid PolyProxyMatrix.")
		return nil
	}
}

func (mat PolyMatrix) BlockCombine(inputPolyProxyMat PolynomialMatrix) PolynomialMatrix {
	if mat.Cols() != inputPolyProxyMat.Cols() {
		log.Panic("BlockCombine: cols of matrices are not equal")
	}

	switch input := inputPolyProxyMat.(type) {
	case PolyMatrix:
		ret := make(PolyMatrix, 0, mat.Rows()+inputPolyProxyMat.Rows())
		ret = append(ret, mat...)
		ret = append(ret, input...)
		return PolyMatrix(ret)
	case PolyQMatrix:
		currentPolyQMatrix := mat.TransformedToPolyQMatrix()

		newMat := make(PolyQMatrix, 0, mat.Rows()+input.Rows())
		newMat = append(newMat, currentPolyQMatrix...)
		newMat = append(newMat, input...)

		return newMat
	default:
		log.Panic("Invalid PolyProxyMatrix.")
		return nil
	}
}

func (mat PolyMatrix) MatMul(inputPolyProxyMat PolynomialMatrix) PolynomialMatrix {
	if mat.Cols() != inputPolyProxyMat.Rows() {
		log.Panic("MatMul: Number of cols in first mat is not equal to number of rows in second mat")
	}

	switch input := inputPolyProxyMat.(type) {
	case PolyMatrix:
		inputPolyMatrix := input

		rows, cols := mat.Rows(), mat.Cols()
		otherCols := inputPolyMatrix.Cols()

		newMat := make(PolyMatrix, rows)

		for i := 0; i < rows; i++ {
			currentVec := make(vector.PolyVector, otherCols)

			for j := 0; j < otherCols; j++ {
				currentPoly := poly.NewPoly()

				for k := 0; k < cols; k++ {
					currentPoly = currentPoly.Add(mat[i][k].Mul(inputPolyMatrix[k][j])).(poly.Poly)
				}

				currentVec[j] = currentPoly
			}
			newMat[i] = currentVec
		}

		return newMat
	case PolyQMatrix:
		currentPolyQMatrix := mat.TransformedToPolyQMatrix()

		rows, cols := mat.Rows(), mat.Cols()
		otherCols := input.Cols()

		newMat := make(PolyQMatrix, rows)

		r := devkit.MainRing.AtLevel(devkit.MainRing.Level())

		for i := 0; i < rows; i++ {
			currentVec := make(vector.PolyQVector, otherCols)

			for j := 0; j < otherCols; j++ {
				currentPoly := poly.NewPolyQ()

				for k := 0; k < cols; k++ {
					matNTT := r.NewPoly()
					inputNTT := r.NewPoly()

					r.NTT(currentPolyQMatrix[i][k].Poly, matNTT)
					r.NTT(input[k][j].Poly, inputNTT)

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
	default:
		log.Panic("Invalid PolyProxyMatrix.")
		return nil
	}
}

func (mat PolyMatrix) VecMul(inputPolyProxyVector vector.PolynomialVector) vector.PolynomialVector {
	if inputPolyProxyVector.Length() != mat.Cols() {
		log.Panic("VecMul: vectors don't have the same length")
	}

	switch input := inputPolyProxyVector.(type) {
	case vector.PolyVector:
		ret := make(vector.PolyVector, mat.Rows())

		for i := 0; i < len(ret); i++ {
			currentPoly := make(poly.Poly, mat[0][0].Length())

			for j := 0; j < inputPolyProxyVector.Length(); j++ {
				currentPoly = currentPoly.Add(input[j].Mul(mat[i][j])).(poly.Poly)
			}

			ret[i] = currentPoly
		}
		return ret
	case vector.PolyQVector:
		currentPolyQMatrix := mat.TransformedToPolyQMatrix()

		newVec := make(vector.PolyQVector, currentPolyQMatrix.Rows())
		r := devkit.MainRing.AtLevel(devkit.MainRing.Level())

		for i := 0; i < currentPolyQMatrix.Rows(); i++ {
			currentPoly := poly.NewPolyQ()

			for j := 0; j < inputPolyProxyVector.Length(); j++ {
				matNTT := r.NewPoly()
				inputNTT := r.NewPoly()

				r.MulCoeffsBarrettThenAdd(inputNTT, matNTT, currentPoly.Poly)
			}

			r.INTT(currentPoly.Poly, currentPoly.Poly)
			newVec[i] = currentPoly
		}

		return newVec
	default:
		log.Panic("Invalid PolyProxyMatrix.")
		return nil
	}
}

func (mat PolyMatrix) Equals(other PolyMatrix) bool {
	for i := 0; i < mat.Rows(); i++ {
		if !mat[i].Equals(other[i]) {
			return false
		}
	}
	return true
}
