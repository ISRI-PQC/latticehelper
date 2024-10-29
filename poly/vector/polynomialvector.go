package vector

type PolynomialVector interface {
	String() string
	CoeffString() string
	Serialize() []byte
	Length() int
	Listize() []int64
	// ScaleByPolyProxy(inputPolyProxy poly.Polynomial) PolynomialVector
	// ScaleByInt(input int64) PolynomialVector
	// Add(inputPolyProxyVector PolynomialVector) PolynomialVector
	// Sub(inputPolyProxyVector PolynomialVector) PolynomialVector
	// Concat(inputPolyProxyVector PolynomialVector) PolynomialVector
	// DotProduct(inputPolyProxyVector PolynomialVector) poly.Polynomial
	// Equals(other PolynomialVector) bool
}
