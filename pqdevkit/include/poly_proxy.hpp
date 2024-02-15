#ifndef PQDEVKIT_POLY_PROXY_HPP
#define PQDEVKIT_POLY_PROXY_HPP

#include <limits>
#include "params.hpp"

// TODO: consider not using PolyProxy at all and specify helper functions -
// what about scaling with scalar?
namespace pqdevkit
{
    // NFL specific
    using poly_type = PQDEVKIT_POLY_TYPE;
    using coeff_type = PQDEVKIT_COEFF_TYPE;

    // TODO: make this a template for easy switching of the math library
    class PolyProxy
    {
    private:
        std::unique_ptr<poly_type> poly_ptr;

    public:
        unsigned short degree = PQDEVKIT_DEGREE;
        size_t coeff_modulus = PQDEVKIT_COEFF_MODULUS;

        PolyProxy(coeff_type constant);
        PolyProxy(std::initializer_list<coeff_type> coefficients); // {1,2,3}
        PolyProxy(const poly_type &poly);
        ~PolyProxy();

        poly_type &get_poly() const;

        coeff_type infinite_norm() const;
        std::vector<coeff_type> listize() const;

        PolyProxy operator+(const PolyProxy &other) const;
        PolyProxy operator-(const PolyProxy &other) const;
        PolyProxy operator*(const PolyProxy &other) const;
        PolyProxy operator*(const coeff_type &scalar) const;

        static PolyProxy random_poly();
    };

} // namespace pqdevkit

#endif // PQDEVKIT_POLY_PROXY_HPP