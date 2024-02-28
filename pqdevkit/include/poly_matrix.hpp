#ifndef PQDEVKIT_POLY_MATRIX_HPP
#define PQDEVKIT_POLY_MATRIX_HPP

#include "poly_vector.hpp"

// TODO: maybe use std::list instead of std::vector
// (https://baptiste-wicht.com/posts/2012/12/cpp-benchmark-vector-list-deque.html)
namespace pqdevkit
{
  template <unsigned short _degree, size_t _coeff_modulus> class PolyMatrix
  {
  private:
    std::vector<PolyVector> poly_matrix;

  public:
    PolyMatrix(
      const std::initializer_list<
        std::initializer_list<PolyProxy<_degree, _coeff_modulus>>> &other);
    PolyMatrix(
      const std::initializer_list<std::initializer_list<std::initializer_list<
        typename PolyProxy<_degree, _coeff_modulus>::coeff_type>>> &other);
    PolyMatrix(const std::vector<PolyVector<_degree, _coeff_modulus>> &other);
    PolyMatrix(const PolyMatrix &other);
    ~PolyMatrix();

    const std::vector<PolyVector<_degree, _coeff_modulus>> &get_matrix() const;

    size_t rows() const;
    size_t cols() const;

    typename PolyProxy<_degree, _coeff_modulus>::coeff_type
    infinite_norm() const;
    std::vector<typename PolyProxy<_degree, _coeff_modulus>::coeff_type>
    listize() const;

    PolyMatrix transposed() const;

    PolyMatrix
    scale(const typename PolyProxy<_degree, _coeff_modulus>::coeff_type
            &scalar) const;
    PolyMatrix
    scale(const typename PolyProxy<_degree, _coeff_modulus>::poly_type &poly)
      const;
    PolyMatrix operator+(const PolyMatrix &other) const;
    PolyMatrix operator-(const PolyMatrix &other) const;
    PolyMatrix operator|(const PolyMatrix &other) const;
    PolyMatrix operator/(const PolyMatrix &other) const;
    PolyMatrix operator*(const PolyMatrix &other) const;
    PolyVector<_degree, _coeff_modulus>
    operator*(const PolyVector<_degree, _coeff_modulus> &other) const;
    PolyMatrix
    operator*(const typename PolyProxy<_degree, _coeff_modulus>::coeff_type
                &scalar) const;

    static PolyMatrix random_poly_matrix(size_t rows, size_t cols);
    static PolyMatrix identity_matrix(size_t size);
    static PolyMatrix zero_matrix(size_t rows, size_t cols);
  };

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyMatrix<_degree, _coeff_modulus> operator*(
    const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar,
    const PolyMatrix<_degree, _coeff_modulus> &poly_matrix);
} // namespace pqdevkit

#endif // PQDEVKIT_POLY_MATRIX_HPP