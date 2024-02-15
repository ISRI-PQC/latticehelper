#ifndef PQDEVKIT_POLY_PROXY_HPP
#define PQDEVKIT_POLY_PROXY_HPP

#include <limits>

#include "nfl.hpp"
#include "params.hpp"

// TODO: consider using classes and having private members
// TODO: consider not using PolyProxy at all and specify helper functions -
// what about scaling with scalar?
namespace pqdevkit
{
    using poly_type =
        nfl::poly_from_modulus<uint64_t, PQDEVKIT_DEGREE,
                               nfl::params<uint64_t>::kModulusBitsize>;
    using coeff_type = typename poly_type::value_type;

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

        poly_type& get_poly() const;

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