package poly

type Polynomial interface {
	String() string
	CoeffString() string
	Serialize() []byte
	Listize() []int64
	Length() int
	// Neg() Polynomial
	// InfiniteNorm() uint64
	// Add(inputPolynomial Polynomial) Polynomial
	// Sub(inputPolynomial Polynomial) Polynomial
	// Mul(inputPolynomial Polynomial) Polynomial
	// Pow(exp int64) Polynomial
	// ScaleByInt(scalar int64) Polynomial
	// AddToFirstCoeff(input int64) Polynomial
	// ApplyToEveryCoeff(f func(int64) any) Polynomial
	// Equals(other Polynomial) bool
}