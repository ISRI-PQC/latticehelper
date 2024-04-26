package main

import (
	"github.com/tuneinsight/lattigo/v5/ring"
	"github.com/tuneinsight/lattigo/v5/utils/sampling"
)

var MainRing *ring.Ring
var MainUniformSampler *ring.UniformSampler

func Init(degree int, modulus uint64) error {
	r, err := ring.NewRing(degree, []uint64{modulus})

	if err != nil {
		return err
	}

	MainRing = r

	prng, err := sampling.NewPRNG()

	if err != nil {
		return err
	}

	us := ring.NewUniformSampler(prng, r)

	MainUniformSampler = us

	return nil
}