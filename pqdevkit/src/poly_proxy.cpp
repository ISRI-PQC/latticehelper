#include "poly_proxy.hpp"

namespace pqdevkit
{
    // PolyProxy
    PolyProxy::PolyProxy(coeff_type constant)
    {
        this->poly_p = std::make_unique<poly_type>(constant);
        this->poly_p->ntt_pow_phi();
    }

    PolyProxy::PolyProxy(std::initializer_list<coeff_type> coefficients)
    {
        this->poly_p = std::make_unique<poly_type>(coefficients);
        this->poly_p->ntt_pow_phi();
    }

    PolyProxy::PolyProxy(const poly_type &poly_p)
    {
        this->poly_p = std::make_unique<poly_type>(poly_p);
    }

    PolyProxy::~PolyProxy() {}

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
        return PolyProxy(poly_type(*this->poly_p + *other.poly_p));
    }

    PolyProxy PolyProxy::operator-(const PolyProxy &other) const
    {
        return PolyProxy(poly_type(*this->poly_p - *other.poly_p));
    }

    PolyProxy PolyProxy::operator*(const PolyProxy &other) const
    {
        return PolyProxy(poly_type(*this->poly_p * *other.poly_p));
    }

    PolyProxy PolyProxy::operator*(const coeff_type &scalar) const
    {
        throw std::runtime_error("Not implemented"); // TODO: how do I get, modify and save coeffs from NFLlib?
    }

    PolyProxy PolyProxy::random_poly()
    {
        return PolyProxy(poly_type(nfl::uniform()));
    }
} // namespace pqdevkit
