#include "poly_proxy.hpp"

namespace pqdevkit
{
    /// @brief DOES convert to NTT
    /// @param constant
    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus>::PolyProxy(const typename PolyProxy<_degree, _coeff_modulus>::coeff_type constant)
    {
        this->underlying_poly = poly_type<_degree, _coeff_modulus>(constant);
        this->underlying_poly.ntt_pow_phi();
        this->ntt_from = true;
    }

    /// @brief DOES convert to NTT
    /// @param coefficients
    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus>::PolyProxy(const std::initializer_list<typename PolyProxy<_degree, _coeff_modulus>::coeff_type> coefficients)
    {
        this->underlying_poly = poly_type<_degree, _coeff_modulus>(coefficients);
        this->underlying_poly.ntt_pow_phi();
        this->ntt_from = true;
    }

    /// @brief DOES NOT convert to NTT
    /// @param poly
    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus>::PolyProxy(const poly_type &other, const bool ntt_from)
    {
        this->underlying_poly = poly_type<_degree, _coeff_modulus>(other);
        this->ntt_from = ntt_from;
    }

    /// @brief DOES NOT convert to NTT
    /// @param other
    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus>::PolyProxy(const PolyProxy &other)
    {
        this->underlying_poly = poly_type<_degree, _coeff_modulus>(other.underlying_poly);
        this->ntt_from = other.ntt_from;
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus>::~PolyProxy() {}

    template <unsigned short _degree, size_t _coeff_modulus>
    bool PolyProxy<_degree, _coeff_modulus>::is_ntt() const
    {
        return this->ntt_from;
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    const PolyProxy<_degree, _coeff_modulus>::poly_type &PolyProxy<_degree, _coeff_modulus>::get_poly() const
    {
        return this->underlying_poly;
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    typename PolyProxy<_degree, _coeff_modulus>::coeff_type PolyProxy<_degree, _coeff_modulus>::infinite_norm() const
    {
        throw std::runtime_error(
            "Not implemented"); // TODO: how do I get coeffs from NFLlib?
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    std::vector<typename PolyProxy<_degree, _coeff_modulus>::coeff_type> PolyProxy<_degree, _coeff_modulus>::listize() const
    {
        throw std::runtime_error(
            "Not implemented"); // TODO: how do I get coeffs from NFLlib?
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus> PolyProxy<_degree, _coeff_modulus>::operator-() const
    {
        throw std::runtime_error(
            "Not implemented"); // TODO: how do I get coeffs from NFLlib?
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus> PolyProxy<_degree, _coeff_modulus>::operator+(const PolyProxy &other) const
    {
        return poly_type<_degree, _coeff_modulus>(this->underlying_poly + other.underlying_poly);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus> PolyProxy<_degree, _coeff_modulus>::operator-(const PolyProxy &other) const
    {
        return poly_type<_degree, _coeff_modulus>(this->underlying_poly - other.underlying_poly);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus> PolyProxy<_degree, _coeff_modulus>::operator*(const PolyProxy &other) const
    {
        return poly_type<_degree, _coeff_modulus>(this->underlying_poly * other.underlying_poly);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus> PolyProxy<_degree, _coeff_modulus>::operator*(const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar) const
    {
        throw std::runtime_error("Not implemented"); // TODO: how do I get, modify and save coeffs from NFLlib?
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus> operator*(const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar, const PolyProxy<_degree, _coeff_modulus> &poly_proxy)
    {
        return poly_proxy * scalar;
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyProxy<_degree, _coeff_modulus> PolyProxy<_degree, _coeff_modulus>::random_poly()
    {
        return poly_type<_degree, _coeff_modulus>(nfl::uniform());
    }
} // namespace pqdevkit
