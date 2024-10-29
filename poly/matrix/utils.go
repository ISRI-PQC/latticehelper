package matrix

import (
	"cyber.ee/pq/latticehelper"
	"cyber.ee/pq/latticehelper/poly"
)

// Used to get polynomial matrices to integer matrices while
// preserving the structure
func toeplitz(input poly.PolyQ) [][]int64 {
	flist := input.Poly.Coeffs[latticehelper.MainRing.Level()]
	F := make([][]int64, len(flist))
	for i := 0; i < len(flist); i++ {
		F[i] = make([]int64, len(flist))
		for j := 0; j < len(flist); j++ {
			multiplier := int64(1)
			if j > i {
				multiplier = -1
			}
			pos := i - j

			if pos < 0 {
				pos = len(flist) + pos
			}

			F[i][j] = multiplier * int64(flist[pos])
		}
	}
	return F
}

func BigToeplitz(A PolyQMatrix, m, n int) [][]int64 {
	source := make([][][][]int64, m)
	for i := 0; i < m; i++ {
		source[i] = make([][][]int64, n)
		for j := 0; j < n; j++ {
			source[i][j] = toeplitz(A[i][j])
		}
	}

	result := make([][]int64, m*latticehelper.MainRing.N())

	for row := 0; row < m*latticehelper.MainRing.N(); row++ {
		result[row] = make([]int64, n*latticehelper.MainRing.N())
		for col := 0; col < n*latticehelper.MainRing.N(); col++ {
			result[row][col] = source[latticehelper.FloorDivision(
				int64(row), int64(latticehelper.MainRing.N()),
			)][latticehelper.FloorDivision(
				int64(col), int64(latticehelper.MainRing.N()),
			)][row%latticehelper.MainRing.N()][col%latticehelper.MainRing.N()]
		}
	}

	return result
}

func Transpose(input [][]int64) [][]int64 {
	rows, cols := len(input), len(input[0])
	result := make([][]int64, cols)

	for i := 0; i < cols; i++ {
		row := make([]int64, rows)

		for j := 0; j < rows; j++ {
			row[j] = input[j][i]
		}

		result[i] = row
	}
	return result
}
