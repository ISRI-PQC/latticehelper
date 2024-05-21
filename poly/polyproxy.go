package poly

type PolyProxy interface {
	String() string
	CoeffString() string
	Serialize() []byte
	Listize() []int64
	Length() int
	Neg() PolyProxy
	// InfiniteNorm() uint64
	Add(inputPolyProxy PolyProxy) PolyProxy
	Sub(inputPolyProxy PolyProxy) PolyProxy
	Mul(inputPolyProxy PolyProxy) PolyProxy
	Pow(exp int64) PolyProxy
	ScaleByInt(scalar int64) PolyProxy
	AddToFirstCoeff(input int64) PolyProxy
	Equals(other PolyProxy) bool
}
