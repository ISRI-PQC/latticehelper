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
	MainRing           *ring.Ring
	MainUniformSampler *ring.UniformSampler
)

func InitSingle(degree int64, modulus uint64) error {
	return InitMultiple(degree, []uint64{modulus})
}

func InitMultiple(degree int64, moduli []uint64) error {
	r, err := ring.NewRing(int(degree), moduli)

	if err != nil {
		return err
	}

	MainRing = r.AtLevel(0)

	prng, err := sampling.NewPRNG()

	if err != nil {
		return err
	}

	us := ring.NewUniformSampler(prng, r)

	MainUniformSampler = us

	return nil
}
