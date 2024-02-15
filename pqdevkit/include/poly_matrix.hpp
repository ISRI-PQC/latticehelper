#ifndef PQDEVKIT_POLY_MATRIX_HPP
#define PQDEVKIT_POLY_MATRIX_HPP

#include "poly_vector.hpp"

// TODO: consider using classes and having private members
namespace pqdevkit
{
    class PolyMatrix
    {
    private:
        std::unique_ptr<std::vector<PolyVector>> poly_matrix_ptr;

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
        PolyMatrix(std::initializer_list<std::initializer_list<std::initializer_list<coeff_type>>> poly_matrix);
        PolyMatrix(const std::vector<PolyVector> &poly_matrix);
        ~PolyMatrix();

        std::vector<PolyVector> get_matrix() const;

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
} // namespace pqdevkit

#endif // PQDEVKIT_POLY_MATRIX_HPP