# PQDevKit

This repository contains all necessary files and instructions to create wrappers for common PQC constructs required for PQ engineering.

The purpose is to have a common set of full-fledged functions available in multiple components built with different languages. This will allow for a more consistent and efficient development of PQ components. Also, having these functions implemented in C/C++ by people who understand them much better allows for better performance.

## Initial remarks

At first, focus will be on lattice-based cryptography (basically everything required for [TOPCOAT](https://gitlab.cyber.ee/nsnetkov/topcoat/-/tree/python-implementation?ref_type=heads) and [pq-cast-as-intented](https://gitlab.cyber.ee/pq-ivxv/pq-cast-as-intended)).

Libraries like [libNTL](https://libntl.org/) and [BPAS](https://bpaslib.org/) should be utilized for the heavy lifting.
For ease of use, we will also include a wrapper for [liboqs](https://github.com/open-quantum-safe/liboqs)

Initial target languages are Python and Go.

## Project structure

## Prerequisites


## Go notes

```
export LIBOQS_ROOT="CHANGE_ME: <path-to-liboqs-repo>"
export DYLD_FALLBACK_LIBRARY_PATH="$DYLD_FALLBACK_LIBRARY_PATH:$LIBOQS_ROOT/build/lib" # only for macOS
export CGO_CPPFLAGS="-I$LIBOQS_ROOT/build/include"
export CGO_LDFLAGS="-L$LIBOQS_ROOT/build/lib -loqs"```