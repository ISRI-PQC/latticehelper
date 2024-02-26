#ifndef PQDEVKIT_POLY_VECTOR_HPP
#define PQDEVKIT_POLY_VECTOR_HPP

#include "params.hpp"
#include "poly_proxy.hpp"

// TODO: consider using classes and having private members
namespace pqdevkit
{
    class PolyMatrix;

    class PolyVector
    {
    private:
        std::vector<PolyProxy> poly_vector;

    public:
        PolyVector(std::initializer_list<std::initializer_list<coeff_type>> other); // {{1,2,3}, {4,5,6}}
        PolyVector(const std::vector<PolyProxy> &other);
        PolyVector(const PolyVector &other);
        ~PolyVector();

        const std::vector<PolyProxy> &get_vector() const;

        size_t length() const;

        coeff_type infinite_norm() const;
        std::vector<coeff_type> listize() const;

        PolyVector scale(const coeff_type &scalar) const;
        PolyVector scale(const poly_type &poly) const;
        PolyVector operator+(const PolyVector &other) const;
        PolyVector operator-(const PolyVector &other) const;
        PolyVector operator|(const PolyVector &other) const;
        PolyProxy operator*(const PolyVector &other) const;
        PolyVector operator*(const PolyMatrix &other) const;
        PolyVector operator*(const coeff_type &scalar) const;

        static PolyVector random_poly_vector(size_t length);
    };

    PolyVector operator*(const coeff_type &scalar, const PolyVector &poly_vector);
} // namespace pqdevkit

#endif // PQDEVKIT_POLY_VECTOR_HPP