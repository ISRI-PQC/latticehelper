#ifndef PQDEVKIT_POLY_VECTOR_HPP
#define PQDEVKIT_POLY_VECTOR_HPP

#include "poly_proxy.hpp"

// TODO: maybe use std::list instead of std::vector
// (https://baptiste-wicht.com/posts/2012/12/cpp-benchmark-vector-list-deque.html)
namespace pqdevkit
{
  class PolyMatrix;

  class PolyVector
  {
  private:
    std::vector<PolyProxy> poly_vector;

  public:
    PolyVector(const std::initializer_list<std::initializer_list<long>> other);
    PolyVector(const std::vector<PolyProxy> &other);
    PolyVector(const PolyVector &other);
    ~PolyVector();

    const std::vector<PolyProxy> &get_vector() const;

    size_t length() const;

    long infinite_norm() const;
    std::vector<PQDEVKIT_COEFF_TYPE> listize() const;

    PolyVector scale(const PQDEVKIT_COEFF_TYPE &scalar) const;
    PolyVector scale(const PQDEVKIT_POLY_TYPE &poly) const;
    PolyVector operator+(const PolyVector &other) const;
    PolyVector operator-(const PolyVector &other) const;
    PolyVector operator|(const PolyVector &other) const;
    PolyProxy operator*(const PolyVector &other) const;
    PolyVector operator*(const PolyMatrix &other) const;
    PolyVector operator*(const PQDEVKIT_COEFF_TYPE &scalar) const;

    static PolyVector random_poly_vector(size_t length);
  };

  PolyVector
  operator*(const PQDEVKIT_COEFF_TYPE &scalar, const PolyVector &poly_vector);
} // namespace pqdevkit

#endif // PQDEVKIT_POLY_VECTOR_HPP