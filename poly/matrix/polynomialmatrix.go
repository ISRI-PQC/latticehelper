package matrix

type PolynomialMatrix interface {
	String() string
	CoeffString() string
	Serialize() []byte
	Rows() int
	Cols() int
	Listize() []int64
	// Transposed() PolynomialMatrix
	// ScaleByPolyProxy(inputPolyProxy poly.Polynomial) PolynomialMatrix
	// ScaleByInt(input int64) PolynomialMatrix
	// Add(inputPolyProxyMatrix PolynomialMatrix) PolynomialMatrix
	// Sub(inputPolyProxyMatrix PolynomialMatrix) PolynomialMatrix
	// Concat(inputPolyProxyMatrix PolynomialMatrix) PolynomialMatrix
	// BlockCombine(inputPolyProxyMatrix PolynomialMatrix) PolynomialMatrix
	// MatMul(inputPolyProxyMatrix PolynomialMatrix) PolynomialMatrix
	// VecMul(inputPolyProxyVector vector.PolynomialVector) vector.PolynomialVector
	// Equals(other PolynomialMatrix) bool
}
