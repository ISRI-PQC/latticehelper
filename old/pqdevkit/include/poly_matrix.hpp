#ifndef PQDEVKIT_POLY_MATRIX_HPP
#define PQDEVKIT_POLY_MATRIX_HPP

#include "poly_vector.hpp"

// TODO: maybe use std::list instead of std::vector
// (https://baptiste-wicht.com/posts/2012/12/cpp-benchmark-vector-list-deque.html)
namespace pqdevkit
{
  class PolyMatrix
  {
  private:
    std::vector<PolyVector> poly_matrix;

  public:
    PolyMatrix(
      const std::initializer_list<std::initializer_list<PolyProxy>> &other);
    PolyMatrix(const std::initializer_list<
               std::initializer_list<std::initializer_list<long>>> &other);
    PolyMatrix(const std::vector<PolyVector> &other);
    PolyMatrix(const PolyMatrix &other);
    ~PolyMatrix();

    const std::vector<PolyVector> &get_matrix() const;

    size_t rows() const;
    size_t cols() const;

    long infinite_norm() const;
    std::vector<PQDEVKIT_COEFF_TYPE> listize() const;

    PolyMatrix transposed() const;

    PolyMatrix scale(const PQDEVKIT_COEFF_TYPE &scalar) const;
    PolyMatrix scale(const PQDEVKIT_POLY_TYPE &poly) const;
    PolyMatrix operator+(const PolyMatrix &other) const;
    PolyMatrix operator-(const PolyMatrix &other) const;
    PolyMatrix operator|(const PolyMatrix &other) const;
    PolyMatrix operator/(const PolyMatrix &other) const;
    PolyMatrix operator*(const PolyMatrix &other) const;
    PolyVector operator*(const PolyVector &other) const;
    PolyMatrix operator*(const PQDEVKIT_COEFF_TYPE &scalar) const;

    static PolyMatrix random_poly_matrix(size_t rows, size_t cols);
    static PolyMatrix identity_matrix(size_t size);
    static PolyMatrix zero_matrix(size_t rows, size_t cols);
  };

  PolyMatrix
  operator*(const PQDEVKIT_COEFF_TYPE &scalar, const PolyMatrix &poly_matrix);
} // namespace pqdevkit

#endif // PQDEVKIT_POLY_MATRIX_HPP