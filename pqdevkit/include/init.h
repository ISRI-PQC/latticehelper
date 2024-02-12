#ifndef PQDEVKIT_INIT_H
#define PQDEVKIT_INIT_H

#include "params.h"
#include "ZZ_pEX.h"
#include "ZZ_pXFactoring.h"

void init_ntl_context(pqdevkit::params currentParams);

namespace pqdevkit
{
    void init();
    void init(unsigned short max_degree, unsigned long long coeff_modulus);
}

#endif // PQDEVKIT_INIT_H