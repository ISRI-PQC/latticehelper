package matrix

import (
	"strings"

	"cyber.ee/muzosh/pq/devkit"
	"cyber.ee/muzosh/pq/devkit/poly"
	"cyber.ee/muzosh/pq/devkit/poly/vector"
)

type PolyMatrix []vector.PolyVector

func (mat PolyMatrix) Serialize() ([]byte, error) {
	return devkit.SerializeObject(mat)
}

func DeserializePolyMatrix(data []byte) (PolyMatrix, error) {
	var mat PolyMatrix
	err := devkit.DeserializeObject(data, &mat)
	return mat, err
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

// func (mat PolyMatrix) InfiniteNorm() uint64 {

// 	max := uint64(0)
// 	for _, polyQVec := range mat {
// 		maxVec := polyQVec.InfiniteNorm()
// 		if maxVec > max {
// 			max = maxVec
// 		}
// 	}

// 	return max
// }

func (mat PolyMatrix) Transposed() PolyProxyMatrix {
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

	return PolyMatrix(result)
}

func (mat PolyMatrix) ScaleByPolyProxy(inputPoly poly.PolyProxy) PolyProxyMatrix {
	result := make(PolyMatrix, mat.Rows())

	for i, polyVector := range mat {
		result[i] = polyVector.ScaleByPolyProxy(inputPoly).(vector.PolyVector)
	}

	return PolyMatrix(result)
}

func (mat PolyMatrix) ScaleByInt(input int64) PolyProxyMatrix {
	result := make(PolyMatrix, mat.Rows())

	for i, polyVector := range mat {
		result[i] = polyVector.ScaleByInt(input).(vector.PolyVector)
	}

	return PolyMatrix(result)
}

func (mat PolyMatrix) Add(inputPolyProxyMat PolyProxyMatrix) PolyProxyMatrix {
	if mat.Cols() != inputPolyProxyMat.Cols() || mat.Rows() != inputPolyProxyMat.Rows() {
		panic("Add: rows and cols of matrices are not equal")
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
			newMat[i] = polyQVec.Add(input[i]).(vector.PolyQVector)
		}
		return newMat
	default:
		panic("Invalid PolyProxyMatrix.")
	}
}

func (mat PolyMatrix) Sub(inputPolyProxyMat PolyProxyMatrix) PolyProxyMatrix {
	if mat.Cols() != inputPolyProxyMat.Cols() || mat.Rows() != inputPolyProxyMat.Rows() {
		panic("Sub: rows and cols of matrices are not equal")
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
			newMat[i] = polyQVec.Sub(input[i]).(vector.PolyQVector)
		}
		return newMat
	default:
		panic("Invalid PolyProxyMatrix.")
	}
}

func (mat PolyMatrix) Concat(inputPolyProxyMat PolyProxyMatrix) PolyProxyMatrix {
	if mat.Rows() != inputPolyProxyMat.Rows() {
		panic("Concat: rows of matrices are not equal")
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
			newMat[i] = polyQVec.Add(input[i]).(vector.PolyQVector)
		}
		return newMat
	default:
		panic("Invalid PolyProxyMatrix.")
	}
}

func (mat PolyMatrix) BlockCombine(inputPolyProxyMat PolyProxyMatrix) PolyProxyMatrix {
	if mat.Cols() != inputPolyProxyMat.Cols() {
		panic("BlockCombine: cols of matrices are not equal")
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
		panic("Invalid PolyProxyMatrix.")
	}
}

func (mat PolyMatrix) MatMul(inputPolyProxyMat PolyProxyMatrix) PolyProxyMatrix {
	if mat.Cols() != inputPolyProxyMat.Rows() {
		panic("MatMul: Number of cols in first mat is not equal to number of rows in second mat")
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

		for i := 0; i < rows; i++ {
			currentVec := make(vector.PolyQVector, otherCols)

			for j := 0; j < otherCols; j++ {
				currentPoly := poly.NewPolyQ()

				for k := 0; k < cols; k++ {
					devkit.MainRing.NTT(*currentPolyQMatrix[i][k].Poly, *currentPolyQMatrix[i][k].Poly)
					devkit.MainRing.NTT(*input[k][j].Poly, *input[k][j].Poly)

					devkit.MainRing.MulCoeffsBarrettThenAdd(
						*currentPolyQMatrix[i][k].Poly,
						*input[k][j].Poly,
						*currentPoly.Poly)

					devkit.MainRing.INTT(*currentPolyQMatrix[i][k].Poly, *currentPolyQMatrix[i][k].Poly)
					devkit.MainRing.INTT(*input[k][j].Poly, *input[k][j].Poly)
				}

				devkit.MainRing.INTT(*currentPoly.Poly, *currentPoly.Poly)
				currentVec[j] = currentPoly
			}
			newMat[i] = currentVec
		}

		return newMat
	default:
		panic("Invalid PolyProxyMatrix.")
	}
}

func (mat PolyMatrix) VecMul(inputPolyProxyVector vector.PolyProxyVector) vector.PolyProxyVector {
	if inputPolyProxyVector.Length() != mat.Cols() {
		panic("VecMul: vectors don't have the same length")
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

		for i := 0; i < currentPolyQMatrix.Rows(); i++ {
			currentPoly := poly.NewPolyQ()

			for j := 0; j < inputPolyProxyVector.Length(); j++ {
				devkit.MainRing.NTT(*currentPolyQMatrix[i][j].Poly, *currentPolyQMatrix[i][j].Poly)
				devkit.MainRing.NTT(*input[j].Poly, *input[j].Poly)

				devkit.MainRing.MulCoeffsBarrettThenAdd(*input[j].Poly, *currentPolyQMatrix[i][j].Poly, *currentPoly.Poly)

				devkit.MainRing.INTT(*currentPolyQMatrix[i][j].Poly, *currentPolyQMatrix[i][j].Poly)
				devkit.MainRing.INTT(*input[j].Poly, *input[j].Poly)
			}

			devkit.MainRing.INTT(*currentPoly.Poly, *currentPoly.Poly)
			newVec[i] = currentPoly
		}

		return newVec
	default:
		panic("Invalid PolyProxyMatrix.")
	}
}

func (mat PolyMatrix) Equals(other PolyProxyMatrix) bool {
	switch input := other.(type) {
	case PolyMatrix:
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
