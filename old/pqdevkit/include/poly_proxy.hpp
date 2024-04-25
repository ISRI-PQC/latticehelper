#ifndef PQDEVKIT_POLY_PROXY_HPP
#define PQDEVKIT_POLY_PROXY_HPP

#include "pqdevkit_math_lib.hpp"
#include <limits>
#include <memory>
#include <vector>

// TODO: consider not using PolyProxy at all and specify helper functions -
// what about scaling with scalar?
namespace pqdevkit
{
  class PolyProxy
  {
  private:
    PQDEVKIT_POLY_TYPE underlying_poly;

  public:
    PolyProxy(const PQDEVKIT_COEFF_TYPE constant);
    PolyProxy(const std::initializer_list<long> coefficients);
    PolyProxy(const PQDEVKIT_POLY_TYPE &other);
    PolyProxy(const PolyProxy &other);
    ~PolyProxy();

    const PQDEVKIT_POLY_TYPE &get_poly() const;

    long infinite_norm() const;
    std::vector<PQDEVKIT_COEFF_TYPE> listize() const;

    PolyProxy operator-() const;

    PolyProxy operator+(const PolyProxy &other) const;
    PolyProxy operator-(const PolyProxy &other) const;
    PolyProxy operator*(const PolyProxy &other) const;
    PolyProxy operator*(const long &scalar) const;

    static PolyProxy random_poly();
  };

  PolyProxy
  operator*(const long &scalar, const PolyProxy &poly_proxy);

} // namespace pqdevkit

#endif // PQDEVKIT_POLY_PROXY_HPP