#ifndef PQDEVKIT_MATH_LIB_H
#define PQDEVKIT_MATH_LIB_H

// Include everything needed for math lib

#include <NTL/ZZ.h>
#include <NTL/ZZ_pE.h>
#include <NTL/ZZ_pEX.h>
#include <NTL/ZZ_pX.h>
#include <NTL/ZZ_pXFactoring.h>

extern unsigned short PQDEVKIT_DEGREE;
extern unsigned long PQDEVKIT_COEFF_MODULUS;
extern NTL::ZZ_pEXModulus PQDEVKIT_MODULUS;

void init_ntl(unsigned short degree, size_t coeff_modulus);

#define PQDEVKIT_POLY_TYPE NTL::ZZ_pEX
#define PQDEVKIT_COEFF_TYPE typename PQDEVKIT_POLY_TYPE::coeff_type

// Define PQDEVKIT_POLY_TYPE and PQDEVKIT_COEFF_TYPE
#define PQDEVKIT_POLY_TYPE NTL::ZZ_pEX
#define PQDEVKIT_COEFF_TYPE typename PQDEVKIT_POLY_TYPE::coeff_type

#endif // PQDEVKIT_MATH_LIB_H