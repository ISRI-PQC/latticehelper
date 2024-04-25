#include "pqdevkit_math_lib.hpp"

// Define the static variables
unsigned short PQDEVKIT_DEGREE = 0;
unsigned long PQDEVKIT_COEFF_MODULUS = 0;
NTL::ZZ_pEXModulus PQDEVKIT_MODULUS;

// Implementation of the init_ntl function
void init_ntl(unsigned short degree, size_t coeff_modulus)
{
    PQDEVKIT_DEGREE = degree;
    PQDEVKIT_COEFF_MODULUS = coeff_modulus;

    using namespace NTL;

    ZZ_p::init(ZZ(coeff_modulus));

    ZZ_pX P;

    SetCoeff(P, degree, 1);
    SetCoeff(P, 0, 1);

    ZZ_pE::init(P);

    build(PQDEVKIT_MODULUS, to_ZZ_pEX(ZZ_pE::modulus()));
}