#include "poly_proxy.hpp"

namespace pqdevkit
{
    /// @brief DOES convert to NTT
    /// @param constant
    PolyProxy::PolyProxy(coeff_type constant)
    {
        poly_ptr = std::make_unique<poly_type>(constant);
        poly_ptr->ntt_pow_phi();
    }

    /// @brief DOES convert to NTT
    /// @param coefficients
    PolyProxy::PolyProxy(std::initializer_list<coeff_type> coefficients)
    {
        poly_ptr = std::make_unique<poly_type>(coefficients);
        poly_ptr->ntt_pow_phi();
    }

    /// @brief DOES NOT convert to NTT
    /// @param poly
    PolyProxy::PolyProxy(const poly_type &poly)
    {
        poly_ptr = std::make_unique<poly_type>(poly);
    }

    PolyProxy::~PolyProxy() {}

    /// @brief
    /// @return
    poly_type &PolyProxy::get_poly() const
    {
        return *poly_ptr;
    }

    coeff_type PolyProxy::infinite_norm() const
    {
        throw std::runtime_error(
            "Not implemented"); // TODO: how do I get coeffs from NFLlib?
    }

    std::vector<coeff_type> PolyProxy::listize() const
    {
        throw std::runtime_error(
            "Not implemented"); // TODO: how do I get coeffs from NFLlib?
    }

    PolyProxy PolyProxy::operator+(const PolyProxy &other) const
    {
        return poly_type(*poly_ptr + *other.poly_ptr);
    }

    PolyProxy PolyProxy::operator-(const PolyProxy &other) const
    {
        return poly_type(*poly_ptr - *other.poly_ptr);
    }

    PolyProxy PolyProxy::operator*(const PolyProxy &other) const
    {
        return poly_type(*poly_ptr * *other.poly_ptr);
    }

    PolyProxy PolyProxy::operator*(const coeff_type &scalar) const
    {
        throw std::runtime_error("Not implemented"); // TODO: how do I get, modify and save coeffs from NFLlib?
    }

    PolyProxy PolyProxy::random_poly()
    {
        return poly_type(nfl::uniform());
    }
} // namespace pqdevkit
