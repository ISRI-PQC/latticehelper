#ifndef PQDEVKIT_MATH_LIB_H
#define PQDEVKIT_MATH_LIB_H

// Include everything needed for math lib
#include "nfl.hpp"

// Initial setup or preprocessing

// Define PQDEVKIT_POLY_TYPE and PQDEVKIT_COEFF_TYPE
#define PQDEVKIT_POLY_TYPE                                                    \
  nfl::poly_from_modulus<uint64_t, _degree,                                   \
                         nfl::params<uint64_t>::kModulusBitsize>
#define PQDEVKIT_COEFF_TYPE typename PQDEVKIT_POLY_TYPE::value_type

#endif // PQDEVKIT_MATH_LIB_H