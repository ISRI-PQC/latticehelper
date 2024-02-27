#ifndef PQDEVKIT_POLY_MATRIX_HPP
#define PQDEVKIT_POLY_MATRIX_HPP

#include "pqdevkit_params.hpp"
#include "poly_vector.hpp"

// TODO: consider using classes and having private members
// TODO: maybe use std::list instead of std::vector (https://baptiste-wicht.com/posts/2012/12/cpp-benchmark-vector-list-deque.html)
namespace pqdevkit
{
    class PolyMatrix
    {
    private:
        std::vector<PolyVector> poly_matrix;

    public:
        /*
        {
            {
                {1, 2, 3},
                {4, 5, 6}
            },
            {
                {7, 8, 9},
                {10, 11, 12}
            }
        }
        */
        PolyMatrix(const std::initializer_list<std::initializer_list<std::initializer_list<coeff_type>>> other);
        PolyMatrix(const std::vector<PolyVector> &other);
        PolyMatrix(const PolyMatrix &other);
        ~PolyMatrix();

        const std::vector<PolyVector> &get_matrix() const;

        size_t rows() const;
        size_t cols() const;

        coeff_type infinite_norm() const;
        std::vector<coeff_type> listize() const;

        PolyMatrix transposed() const;

        PolyMatrix scale(const coeff_type &scalar) const;
        PolyMatrix scale(const poly_type &poly) const;
        PolyMatrix operator+(const PolyMatrix &other) const;
        PolyMatrix operator-(const PolyMatrix &other) const;
        PolyMatrix operator|(const PolyMatrix &other) const;
        PolyMatrix operator/(const PolyMatrix &other) const;
        PolyMatrix operator*(const PolyMatrix &other) const;
        PolyVector operator*(const PolyVector &other) const;
        PolyMatrix operator*(const coeff_type &scalar) const;

        static PolyMatrix random_poly_matrix(size_t rows, size_t cols);
        static PolyMatrix identity_matrix(size_t size);
        static PolyMatrix zero_matrix(size_t rows, size_t cols);
    };

    PolyMatrix operator*(const coeff_type &scalar, const PolyMatrix &poly_matrix);
} // namespace pqdevkit

#endif // PQDEVKIT_POLY_MATRIX_HPP