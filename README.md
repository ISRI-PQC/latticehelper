# PQDevKit

Idea to create this library came up from seeing concrete similarities between different custom PQ protocols.

The purpose is to have a common set of full-fledged functions available in multiple components. This will allow for a more consistent and efficient development of PQ components.

Under the hood, it is utilizing [LattiGo library](https://github.com/tuneinsight/lattigo) for it's ring capabilities (and most importantly NTT and INTT functions).

## Features

- Polynomial arithmetic in two different rings:
    - ring `R` over Z[X]/(X^d + 1)
        - coefficients are all natural numbers, poly is modulo X^d + 1
        - in this library, naming is regular `poly...`
    - ring `Rq` over Z_q[X]/(X^d + 1)
        - coefficients are in range from 0 to q-1, poly is modulo X^d + 1
        - in this library, naming is `polyQ...`
- Vector and matrix arithmetic in both rings.
- some util functions like Power2Round, checking bounds, norms, etc.

## Initialization

`devkit.InitSingle()` or `devkit.InitMultiple()` MUST be called at least once before using anything related to `polyQ`. Arguments are `d` and modulus `q`/moduli `\[q1, q2, q3,...\]`. `InitMultiple` can prepare parameters between multiple rings `Rq`, but has not been tested that much yet.

## Concurrency

If not stated otherwise, all functions should be thread safe and do not require any locks.

Only exception to that rule are functions `NewRandomPolyQ{matrix|vector|""}`, which require a thread-unique sampler. If these functions are used concurrently (i.e. called multiple times at the same time), create new sampler in each thread by `devkit.NewSampler` for them.
