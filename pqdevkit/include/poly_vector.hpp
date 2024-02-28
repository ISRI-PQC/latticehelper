#ifndef PQDEVKIT_POLY_VECTOR_HPP
#define PQDEVKIT_POLY_VECTOR_HPP

#include "poly_proxy.hpp"

// TODO: maybe use std::list instead of std::vector (https://baptiste-wicht.com/posts/2012/12/cpp-benchmark-vector-list-deque.html)
namespace pqdevkit
{
    template <unsigned short _degree, size_t _coeff_modulus>
    class PolyMatrix;

    template <unsigned short _degree, size_t _coeff_modulus>
    class PolyVector
    {
    private:
        std::vector<PolyProxy> poly_vector;

    public:
        PolyVector(const std::initializer_list<std::initializer_list<typename PolyProxy<_degree, _coeff_modulus>::coeff_type>> other); // {{1,2,3}, {4,5,6}}
        PolyVector(const std::vector<PolyProxy<_degree, _coeff_modulus>> &other);
        PolyVector(const PolyVector &other);
        ~PolyVector();

        const std::vector<PolyProxy<_degree, _coeff_modulus>> &get_vector() const;

        size_t length() const;

        typename PolyProxy<_degree, _coeff_modulus>::coeff_type infinite_norm() const;
        std::vector<typename PolyProxy<_degree, _coeff_modulus>::coeff_type> listize() const;

        PolyVector scale(const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar) const;
        PolyVector scale(const typename PolyProxy<_degree, _coeff_modulus>::poly_type &poly) const;
        PolyVector operator+(const PolyVector &other) const;
        PolyVector operator-(const PolyVector &other) const;
        PolyVector operator|(const PolyVector &other) const;
        PolyProxy<_degree, _coeff_modulus> operator*(const PolyVector &other) const;
        PolyVector operator*(const PolyMatrix<_degree, _coeff_modulus> &other) const;
        PolyVector operator*(const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar) const;

        static PolyVector random_poly_vector(size_t length);
    };

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyVector<_degree, _coeff_modulus> operator*(const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar, const PolyVector<_degree, _coeff_modulus> &poly_vector);
} // namespace pqdevkit

#endif // PQDEVKIT_POLY_VECTOR_HPP