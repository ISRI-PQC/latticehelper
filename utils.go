package devkit

import (
	"bytes"
	"encoding/gob"
	"math/big"

	"github.com/tuneinsight/lattigo/v5/utils/sampling"
)

var GobBuffer = bytes.Buffer{}

var GobEncoder = gob.NewEncoder(&GobBuffer)
var GobDecoder = gob.NewDecoder(&GobBuffer)

func SerializeObject(obj any) ([]byte, error) {
	GobBuffer.Reset()
	err := GobEncoder.Encode(obj)
	bytes := GobBuffer.Bytes()
	GobBuffer.Reset()
	return bytes, err
}

func DeserializeObject(data []byte, obj any) error {
	GobBuffer.Reset()
	GobBuffer.Write(data)
	err := GobDecoder.Decode(obj)
	GobBuffer.Reset()
	return err
}

func RandUint64() uint64 {
	return sampling.RandUint64()
}

// RandFloat64 returns a random float between min and max.
func RandFloat64(min, max float64) float64 {
	ret := sampling.RandFloat64(min, max)
	return ret
}

func FloorDivision[T int | int64 | uint64](a, b T) uint64 {
	ret := new(big.Int).Div(big.NewInt(int64(a)), big.NewInt(int64(b))).Uint64()
	return ret
}

func InvMod[T int | int64 | uint64](d, q T) uint64 {
	ret := new(big.Int).ModInverse(big.NewInt(int64(d)), big.NewInt(int64(q))).Uint64()
	return ret
}

func MulMod[T int | int64 | uint64](a, b T, m uint64) uint64 {
	ret := new(big.Int).Mod(
		new(big.Int).Mul(big.NewInt(int64(a)), big.NewInt(int64(b))),
		big.NewInt(int64(m)),
	).Uint64()
	return ret
}

func PowMod[T int | int64 | uint64](a, b, m T) uint64 {
	ret := new(big.Int).Exp(big.NewInt(int64(a)), big.NewInt(int64(b)), big.NewInt(int64(m))).Uint64()
	return ret
}

func PositiveMod(a int64, m uint64) uint64 {
	ret := a % int64(m)
	if ret < 0 {
		ret += int64(m)
	}
	return uint64(ret)
}
