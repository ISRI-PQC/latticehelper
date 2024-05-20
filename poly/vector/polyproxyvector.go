package vector

import "cyber.ee/pq/devkit/poly"

type PolyProxyVector interface {
	String() string
	CoeffString() string
	Length() int
	Listize() []int64
	ScaleByPolyProxy(inputPolyProxy poly.PolyProxy) PolyProxyVector
	ScaleByInt(input int64) PolyProxyVector
	Add(inputPolyProxyVector PolyProxyVector) PolyProxyVector
	Sub(inputPolyProxyVector PolyProxyVector) PolyProxyVector
	Concat(inputPolyProxyVector PolyProxyVector) PolyProxyVector
	DotProduct(inputPolyProxyVector PolyProxyVector) poly.PolyProxy
	Equals(other PolyProxyVector) bool
}
