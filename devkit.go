package devkit

import (
	"github.com/tuneinsight/lattigo/v5/ring"
	"github.com/tuneinsight/lattigo/v5/utils/sampling"
)

// If you encounter [cyber.ee/pq/common.MainRing] or
// [cyber.ee/pq/common.MainUniformSampler] being nil,
// you must initialize it first using
// [cyber.ee/pq/common.InitSingle] or
// [cyber.ee/pq/common.InitMultiple] functions!
var (
	MainRing              *ring.Ring
	DefaultUniformSampler *ring.UniformSampler
)

func InitSingle(degree int64, modulus uint64) error {
	return InitMultiple(degree, []uint64{modulus})
}

func InitMultiple(degree int64, moduli []uint64) error {
	r, err := ring.NewRing(int(degree), moduli)

	if err != nil {
		return err
	}

	MainRing = r

	s, err := GetSampler()

	if err != nil {
		return err
	}

	DefaultUniformSampler = s

	return nil
}

func GetSampler() (*ring.UniformSampler, error) {
	prng, err := sampling.NewPRNG()

	if err != nil {
		return nil, err
	}

	us := ring.NewUniformSampler(prng, MainRing.AtLevel(MainRing.Level()))

	return us, nil
}
