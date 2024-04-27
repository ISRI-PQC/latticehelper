package pqdevkit

type PolyMatrix struct {
	PolyVectors []PolyVector
	IsNTT       bool
}

func (mat PolyMatrix) Rows() int {
	return len(mat.PolyVectors)
}

func (mat PolyMatrix) Cols() int {
	return mat.PolyVectors[0].Length()
}

func (mat *PolyMatrix) ToNTT() {
	if mat.IsNTT {
		return
	}
	for _, polyVector := range mat.PolyVectors {
		polyVector.ToNTT()
	}
}

func (mat *PolyMatrix) FromNTT() {
	if !mat.IsNTT {
		return
	}
	for _, polyVector := range mat.PolyVectors {
		polyVector.FromNTT()
	}
}

func (mat PolyMatrix) Listize() []uint64 {
	if mat.IsNTT {
		mat.FromNTT()
	}

	listizedVec := make([]uint64, len(mat.PolyVectors)*mat.PolyVectors[0].Length()*MainRing.N())

	for _, polyVec := range mat.PolyVectors {
		listizedVec = append(listizedVec, polyVec.Listize()...)
	}

	return listizedVec
}

func (mat PolyMatrix) InfiniteNorm() uint64 {
	if mat.IsNTT {
		mat.FromNTT()
	}

	max := uint64(0)
	for _, polyVec := range mat.PolyVectors {
		maxVec := polyVec.InfiniteNorm()
		if maxVec > max {
			max = maxVec
		}
	}

	return max
}

func (mat PolyMatrix) Transposed() PolyMatrix {
	cols := mat.Cols()
	rows := mat.Rows()

	result := make([]PolyVector, cols)

	for i := 0; i < cols; i++ {
		polyVector := PolyVector{make([]PolyProxy, rows), mat.IsNTT}

		for j := 0; j < rows; j++ {
			polyVector.PolyProxies[j] = mat.PolyVectors[j].PolyProxies[i]
		}

		result[i] = polyVector
	}

	return PolyMatrix{result, mat.IsNTT}
}

func (mat PolyMatrix) Scale(input_poly PolyProxy) PolyMatrix {
	result := make([]PolyVector, len(mat.PolyVectors))

	for i, polyVector := range mat.PolyVectors {
		result[i] = polyVector.Scale(input_poly)
	}
	return PolyMatrix{result, mat.IsNTT}
}

func (mat PolyMatrix) ScaleScalar(input uint64) PolyMatrix {
	result := make([]PolyVector, len(mat.PolyVectors))

	for i, polyVector := range mat.PolyVectors {
		result[i] = polyVector.ScaleScalar(input)
	}
	return PolyMatrix{result, mat.IsNTT}
}

func (mat PolyMatrix) Add(input_mat PolyMatrix) PolyMatrix {
	if mat.Cols() != input_mat.Cols() || mat.Rows() != input_mat.Rows() {
		panic("Add: rows and cols of matrices are not equal")
	}

	newMat := make([]PolyVector, mat.Rows())
	for i, polyVec := range mat.PolyVectors {
		newMat[i] = polyVec.Add(input_mat.PolyVectors[i])
	}

	return PolyMatrix{newMat, mat.IsNTT}
}

func (mat PolyMatrix) Sub(input_mat PolyMatrix) PolyMatrix {
	if mat.Cols() != input_mat.Cols() || mat.Rows() != input_mat.Rows() {
		panic("Sub: rows and cols of matrices are not equal")
	}

	newMat := make([]PolyVector, mat.Rows())
	for i, polyVec := range mat.PolyVectors {
		newMat[i] = polyVec.Sub(input_mat.PolyVectors[i])
	}

	return PolyMatrix{newMat, mat.IsNTT}
}

func (mat PolyMatrix) Concat(input_mat PolyMatrix) PolyMatrix {
	if mat.Rows() != input_mat.Rows() {
		panic("Concat: rows of matrices are not equal")
	}

	if mat.IsNTT != input_mat.IsNTT {
		panic("Concat: matrices don't have the same form")
	}

	newMat := make([]PolyVector, mat.Rows())
	for i, polyVec := range mat.PolyVectors {
		newMat[i] = polyVec.Concat(input_mat.PolyVectors[i])
	}

	return PolyMatrix{newMat, mat.IsNTT}
}

func (mat PolyMatrix) BlockCombine(input_mat PolyMatrix) PolyMatrix {
	if mat.Cols() != input_mat.Cols() {
		panic("BlockCombine: cols of matrices are not equal")
	}

	if mat.IsNTT != input_mat.IsNTT {
		panic("BlockCombine: matrices don't have the same form")
	}

	newMat := make([]PolyVector, mat.Rows()+input_mat.Cols())

	newMat = append(newMat, mat.PolyVectors...)
	newMat = append(newMat, input_mat.PolyVectors...)

	return PolyMatrix{newMat, mat.IsNTT}
}

func (mat PolyMatrix) MatMul(input_mat PolyMatrix) PolyMatrix {
	if mat.Cols() != input_mat.Rows() {
		panic("MatMul: Number of cols in first mat is not equal to number of rows in second mat")
	}

	if mat.IsNTT != input_mat.IsNTT {
		panic("Concat: matrices don't have the same form")
	}

	was_ntt := mat.IsNTT

	if !was_ntt {
		mat.ToNTT()
		input_mat.ToNTT()
	}

	rows, cols := mat.Rows(), mat.Cols()
	other_cols := input_mat.Cols()

	newMat := make([]PolyVector, rows)

	for i := 0; i < rows; i++ {
		currentVec := make([]PolyProxy, other_cols)

		for j := 0; j < other_cols; j++ {
			currentPoly := MainRing.NewPoly()

			for k := 0; k < cols; k++ {
				MainRing.MulCoeffsBarrettThenAdd(mat.PolyVectors[i].PolyProxies[k].Poly, input_mat.Transposed().PolyVectors[k].PolyProxies[j].Poly, currentPoly)
			}

			currentVec[j] = PolyProxy{currentPoly, true}
		}
		newMat[i] = PolyVector{currentVec, true}
	}

	if !was_ntt {
		mat.FromNTT()
		input_mat.FromNTT()
	}

	return PolyMatrix{newMat, was_ntt}
}

func NewRandomPolyMatrix(rows, cols int) PolyMatrix {
	polyVectors := make([]PolyVector, rows)
	for i := 0; i < rows; i++ {
		polyVectors[i] = NewRandomPolyVector(cols)
	}
	return PolyMatrix{polyVectors, false}
}

func NewIdentityMatrix(size int) PolyMatrix {
	polyVectors := make([]PolyVector, size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if i == j {
				polyVectors[i].PolyProxies[j].Poly.Coeffs[0][0] = uint64(1)
			}
		}
	}
	return PolyMatrix{polyVectors, false}
}

func NewZeroMatrix(rows, cols int) PolyMatrix {
	polyVectors := make([]PolyVector, rows)
	for i := 0; i < rows; i++ {
		polyVectors[i] = PolyVector{}
	}
	return PolyMatrix{polyVectors, false}
}
