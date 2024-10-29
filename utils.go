package latticehelper

import (
	"math/big"
)

func FloorDivision(a, b int64) int64 {
	ret := new(big.Int).Div(big.NewInt(a), big.NewInt(b)).Int64()
	return ret
}

func InvMod(d, q int64) int64 {
	ret := new(big.Int).ModInverse(big.NewInt(d), big.NewInt(q)).Int64()
	return ret
}

func MulMod(a, b, m int64) int64 {
	ret := new(big.Int).Mod(
		new(big.Int).Mul(big.NewInt(a), big.NewInt(b)),
		big.NewInt(m),
	).Int64()
	return ret
}

func PowMod(a, b, m int64) int64 {
	ret := new(big.Int).Exp(big.NewInt(a), big.NewInt(b), big.NewInt(m)).Int64()
	return ret
}

func Pow(a, b int64) int64 {
	ret := new(big.Int).Exp(big.NewInt(a), big.NewInt(b), nil).Int64()
	return ret
}

func PositiveMod(a, m int64) int64 {
	ret := a % m
	if ret < 0 {
		ret += m
	}
	return ret
}
