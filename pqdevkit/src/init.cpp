#include "init.h"

void init_ntl_context(pqdevkit::params currentParams)
{
    NTL::ZZ_p::init(NTL::ZZ(currentParams.COEFF_MODULUS));
    NTL::ZZ_pX P;

    NTL::BuildIrred(P, currentParams.MAX_DEGREE);

    NTL::ZZ_pE::init(P);

    NTL::ZZ_pEX a()
}

void pqdevkit::init()
{
    pqdevkit::params myParams;
}

void pqdevkit::init(unsigned short max_degree, unsigned long long coeff_modulus)
{
    pqdevkit::params myParams;
    myParams.MAX_DEGREE = max_degree;
    myParams.COEFF_MODULUS = coeff_modulus;
}