#ifndef PQDEVKIT_PARAMS_H
#define PQDEVKIT_PARAMS_H

#define PQDEVKIT_DEGREE 128
#define PQDEVKIT_COEFF_MODULUS 4294954753 // prime under 2^32

#include "nfl.hpp"
#define PQDEVKIT_POLY_TYPE nfl::poly_from_modulus<uint64_t, PQDEVKIT_DEGREE, nfl::params<uint64_t>::kModulusBitsize>
#define PQDEVKIT_COEFF_TYPE typename PQDEVKIT_POLY_TYPE::value_type

#endif // PQDEVKIT_PARAMS_H