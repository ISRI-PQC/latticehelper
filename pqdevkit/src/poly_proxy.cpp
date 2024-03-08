#include "poly_proxy.hpp"

namespace pqdevkit
{
  PolyProxy::PolyProxy(const PQDEVKIT_COEFF_TYPE constant)
  {
    this->underlying_poly = PQDEVKIT_POLY_TYPE(constant);
  }

  PolyProxy::PolyProxy(const std::initializer_list<long> coefficients)
  {
    if(coefficients.size() > PQDEVKIT_DEGREE)
      {
        throw std::runtime_error("Degree of polynomial exceeds the limit");
      }

    this->underlying_poly
      = PQDEVKIT_POLY_TYPE(NTL::INIT_SIZE, coefficients.size());

    for(auto it = coefficients.begin(); it != coefficients.end(); it++)
      {
        NTL::SetCoeff(this->underlying_poly, it - coefficients.begin(), *it);
      }
  }

  PolyProxy::PolyProxy(const PQDEVKIT_POLY_TYPE &other)
  {
    this->underlying_poly = other;
  }

  PolyProxy::PolyProxy(const PolyProxy &other)
  {
    this->underlying_poly = other.underlying_poly;
  }

  PolyProxy::~PolyProxy() {}

  const PQDEVKIT_POLY_TYPE &PolyProxy::get_poly() const
  {
    return this->underlying_poly;
  }

  long PolyProxy::infinite_norm() const
  {
    throw std::runtime_error("Not implemented");
  }

  std::vector<PQDEVKIT_COEFF_TYPE> PolyProxy::listize() const
  {
    std::vector<PQDEVKIT_COEFF_TYPE> coefficients;
    for(int i = 0; i <= this->underlying_poly.rep.length(); i++)
      {
        coefficients.emplace_back(NTL::coeff(this->underlying_poly, i));
      }
    return coefficients;
  }

  PolyProxy PolyProxy::operator-() const { return -this->underlying_poly; }

  PolyProxy PolyProxy::operator+(const PolyProxy &other) const
  {
    return this->underlying_poly + other.underlying_poly;
  }

  PolyProxy PolyProxy::operator-(const PolyProxy &other) const
  {
    return this->underlying_poly - other.underlying_poly;
  }

  PolyProxy PolyProxy::operator*(const PolyProxy &other) const
  {
    return NTL::MulMod(this->underlying_poly, other.underlying_poly,
                       PQDEVKIT_MODULUS);
  }

  PolyProxy PolyProxy::operator*(const long &scalar) const
  {
    return this->underlying_poly * scalar;
  }

  PolyProxy operator*(const long &scalar, const PolyProxy &poly_proxy)
  {
    return poly_proxy * scalar;
  }

  PolyProxy PolyProxy::random_poly() { return NTL::random_ZZ_pE(); }
} // namespace pqdevkit
