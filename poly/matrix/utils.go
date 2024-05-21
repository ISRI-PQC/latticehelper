package matrix

import (
	"cyber.ee/pq/devkit"
	"cyber.ee/pq/devkit/poly"
)

// Used to get polynomial matrices to integer matrices while
// preserving the structure
func toeplitz(input poly.PolyQ) [][]int64 {
	flist := input.Poly.Coeffs[devkit.MainRing.Level()]
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

	result := make([][]int64, m*devkit.MainRing.N())

	for row := 0; row < m*devkit.MainRing.N(); row++ {
		result[row] = make([]int64, n*devkit.MainRing.N())
		for col := 0; col < n*devkit.MainRing.N(); col++ {
			result[row][col] = source[devkit.FloorDivision(
				int64(row), int64(devkit.MainRing.N()),
			)][devkit.FloorDivision(
				int64(col), int64(devkit.MainRing.N()),
			)][row%devkit.MainRing.N()][col%devkit.MainRing.N()]
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
