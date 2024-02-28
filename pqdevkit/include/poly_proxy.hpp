#ifndef PQDEVKIT_POLY_PROXY_HPP
#define PQDEVKIT_POLY_PROXY_HPP

#include "pqdevkit_math_lib.hpp"
#include <limits>
#include <memory>

// TODO: consider not using PolyProxy at all and specify helper functions -
// what about scaling with scalar?
namespace pqdevkit
{
  template <unsigned short _degree, size_t _coeff_modulus> class PolyProxy
  {
  private:
    poly_type underlying_poly;
    bool ntt_from = false;

  public:
    // NFL specific
    typedef PQDEVKIT_POLY_TYPE poly_type;
    typedef PQDEVKIT_COEFF_TYPE coeff_type;

    unsigned short degree = _degree;
    size_t coeff_modulus = _coeff_modulus;

    PolyProxy(const coeff_type constant);
    PolyProxy(const std::initializer_list<coeff_type> coefficients);
    PolyProxy(const poly_type &other, const bool ntt_from = true);
    PolyProxy(const PolyProxy &other);
    ~PolyProxy();

    bool is_ntt() const;

    const poly_type &get_poly() const;

    coeff_type infinite_norm() const;
    std::vector<coeff_type> listize() const;

    PolyProxy operator-() const;

    PolyProxy operator+(const PolyProxy &other) const;
    PolyProxy operator-(const PolyProxy &other) const;
    PolyProxy operator*(const PolyProxy &other) const;
    PolyProxy operator*(const coeff_type &scalar) const;

    static PolyProxy random_poly();
  };

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyProxy<_degree, _coeff_modulus> operator*(
    const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar,
    const PolyProxy<_degree, _coeff_modulus> &poly_proxy);

} // namespace pqdevkit

#endif // PQDEVKIT_POLY_PROXY_HPP