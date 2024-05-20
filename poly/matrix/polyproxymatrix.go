package matrix

import (
	"cyber.ee/pq/devkit/poly"
	"cyber.ee/pq/devkit/poly/vector"
)

type PolyProxyMatrix interface {
	String() string
	CoeffString() string
	Rows() int
	Cols() int
	Listize() []int64
	Transposed() PolyProxyMatrix
	ScaleByPolyProxy(inputPolyProxy poly.PolyProxy) PolyProxyMatrix
	ScaleByInt(input int64) PolyProxyMatrix
	Add(inputPolyProxyMatrix PolyProxyMatrix) PolyProxyMatrix
	Sub(inputPolyProxyMatrix PolyProxyMatrix) PolyProxyMatrix
	Concat(inputPolyProxyMatrix PolyProxyMatrix) PolyProxyMatrix
	BlockCombine(inputPolyProxyMatrix PolyProxyMatrix) PolyProxyMatrix
	MatMul(inputPolyProxyMatrix PolyProxyMatrix) PolyProxyMatrix
	VecMul(inputPolyProxyVector vector.PolyProxyVector) vector.PolyProxyVector
	Equals(other PolyProxyMatrix) bool
}
