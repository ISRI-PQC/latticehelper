#ifndef PQDEVKIT_POLY_MATRIX_HPP
#define PQDEVKIT_POLY_MATRIX_HPP

#include "poly_vector.hpp"

// TODO: consider using classes and having private members
namespace pqdevkit
{
    struct PolyMatrix
    {
    private:
        size_t rows;
        size_t cols;
        std::unique_ptr<std::vector<PolyVector>> poly_matrix_ptr;

    public:
        PolyMatrix(std::initializer_list<std::initializer_list<std::initializer_list<coeff_type>>> poly_matrix);
        ~PolyMatrix();

        coeff_type infinite_norm() const;
        std::vector<coeff_type> listize() const;

        void transpose();

        PolyMatrix scale(const coeff_type &scalar) const;
        PolyMatrix scale(const poly_type &poly) const;
        PolyMatrix operator+(const PolyMatrix &other) const;
        PolyMatrix operator-(const PolyMatrix &other) const;
        PolyMatrix operator|(const PolyMatrix &other) const;
        PolyMatrix operator/(const PolyMatrix &other) const;
        PolyMatrix operator*(const PolyMatrix &other) const;
        PolyMatrix operator*(const PolyVector &other) const;
        PolyMatrix operator*(const coeff_type &scalar) const;

        static PolyMatrix random_poly_matrix(size_t rows, size_t cols);
        static PolyMatrix identity_matrix(size_t size);
        static PolyMatrix zero_matrix(size_t rows, size_t cols);
    };
} // namespace pqdevkit

#endif // PQDEVKIT_POLY_MATRIX_HPP