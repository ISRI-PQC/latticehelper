#include "poly_proxy.hpp"

namespace pqdevkit
{
    /// @brief DOES convert to NTT
    /// @param constant
    PolyProxy::PolyProxy(const coeff_type constant)
    {
        this->underlying_poly = poly_type(constant);
        this->underlying_poly.ntt_pow_phi();
        this->ntt_from = true;
    }

    /// @brief DOES convert to NTT
    /// @param coefficients
    PolyProxy::PolyProxy(const std::initializer_list<coeff_type> coefficients)
    {
        this->underlying_poly = poly_type(coefficients);
        this->underlying_poly.ntt_pow_phi();
        this->ntt_from = true;
    }

    /// @brief DOES NOT convert to NTT
    /// @param poly
    PolyProxy::PolyProxy(const poly_type &other, const bool ntt_from)
    {
        this->underlying_poly = poly_type(other);
        this->ntt_from = ntt_from;
    }

    /// @brief DOES NOT convert to NTT
    /// @param other
    PolyProxy::PolyProxy(const PolyProxy &other)
    {
        this->underlying_poly = poly_type(other.underlying_poly);
        this->ntt_from = other.ntt_from;
    }

    PolyProxy::~PolyProxy() {}

    bool PolyProxy::is_ntt() const
    {
        return this->ntt_from;
    }

    const poly_type &PolyProxy::get_poly() const
    {
        return this->underlying_poly;
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
        return poly_type(this->underlying_poly + other.underlying_poly);
    }

    PolyProxy PolyProxy::operator-(const PolyProxy &other) const
    {
        return poly_type(this->underlying_poly - other.underlying_poly);
    }

    PolyProxy PolyProxy::operator*(const PolyProxy &other) const
    {
        return poly_type(this->underlying_poly * other.underlying_poly);
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
