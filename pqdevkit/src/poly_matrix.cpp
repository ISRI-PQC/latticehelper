#include "poly_matrix.hpp"

namespace pqdevkit
{

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus>::PolyMatrix(const std::initializer_list<std::initializer_list<PolyProxy<_degree, _coeff_modulus>>> &other)
    {
        std::vector<PolyVector> result;

        for (auto row : other)
        {
            result.push_back(PolyVector(row));
        }

        this->poly_matrix = std::move(result);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus>::PolyMatrix(const std::initializer_list<std::initializer_list<std::initializer_list<typename PolyProxy<_degree, _coeff_modulus>::coeff_type>>> &other)
    {
        std::vector<PolyVector> result;

        for (auto row : other)
        {
            result.push_back(PolyVector(row));
        }

        this->poly_matrix = std::move(result);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus>::PolyMatrix(const std::vector<PolyVector<_degree, _coeff_modulus>> &other)
    {
        this->poly_matrix = other;
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus>::PolyMatrix(const PolyMatrix &other)
    {
        this->poly_matrix = other.poly_matrix;
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus>::~PolyMatrix() {}

    template <unsigned short _degree, size_t _coeff_modulus>
    const std::vector<PolyVector<_degree, _coeff_modulus>> &PolyMatrix<_degree, _coeff_modulus>::get_matrix() const
    {
        return this->poly_matrix;
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    size_t PolyMatrix<_degree, _coeff_modulus>::rows() const
    {
        return this->poly_matrix.size();
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    size_t PolyMatrix<_degree, _coeff_modulus>::cols() const
    {
        return this->poly_matrix.front().length();
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    typename PolyProxy<_degree, _coeff_modulus>::coeff_type PolyMatrix<_degree, _coeff_modulus>::infinite_norm() const
    {
        coeff_type maxNorm = std::numeric_limits<coeff_type>::min();

        for (const auto &polyVector : this->poly_matrix)
        {
            coeff_type currentNorm = polyVector.infinite_norm();
            if (currentNorm > maxNorm)
            {
                maxNorm = currentNorm;
            }
        }

        return maxNorm;
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    std::vector<typename PolyProxy<_degree, _coeff_modulus>::coeff_type> PolyMatrix<_degree, _coeff_modulus>::listize() const
    {
        std::vector<coeff_type> result;

        for (const auto &polyVector : this->poly_matrix)
        {
            std::vector<coeff_type> currentList = polyVector.listize();
            result.insert(result.end(), currentList.begin(), currentList.end());
        }

        return result;
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::transposed() const
    {
        std::vector<PolyVector> result;

        for (size_t i = 0; i < cols(); i++)
        {
            std::vector<PolyProxy> currentColumn;

            for (size_t j = 0; j < rows(); j++)
            {
                currentColumn.push_back(this->poly_matrix[j].get_vector()[i]);
            }

            result.push_back(PolyVector(currentColumn));
        }

        return PolyMatrix(result);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::scale(const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar) const
    {
        std::vector<PolyVector> scaledPolyMatrix;

        for (const auto &polyVector : this->poly_matrix)
        {
            scaledPolyMatrix.push_back(polyVector.scale(scalar));
        }

        return PolyMatrix(scaledPolyMatrix);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::scale(const typename PolyProxy<_degree, _coeff_modulus>::poly_type &poly) const
    {
        std::vector<PolyVector> scaledPolyMatrix;

        for (const auto &polyVector : this->poly_matrix)
        {
            scaledPolyMatrix.push_back(polyVector.scale(poly));
        }

        return PolyMatrix(scaledPolyMatrix);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::operator+(const PolyMatrix &other) const
    {
        if (rows() != other.rows() || cols() != other.cols())
        {
            throw std::runtime_error("PolyMatrix::operator+: PolyMatrices must have the same dimensions");
        }

        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows(); i++)
        {
            result.push_back(this->poly_matrix[i] + other.poly_matrix[i]);
        }

        return PolyMatrix(result);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::operator-(const PolyMatrix &other) const
    {
        if (rows() != other.rows() || cols() != other.cols())
        {
            throw std::runtime_error("PolyMatrix::operator-: PolyMatrices must have the same dimensions");
        }

        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows(); i++)
        {
            result.push_back(this->poly_matrix[i] - other.poly_matrix[i]);
        }

        return PolyMatrix(result);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::operator|(const PolyMatrix &other) const
    {
        if (rows() != other.rows())
        {
            throw std::runtime_error("PolyMatrix::operator|: PolyMatrices must have the same number of rows");
        }

        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows(); i++)
        {
            result.push_back(this->poly_matrix[i] | other.poly_matrix[i]);
        }

        return PolyMatrix(result);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::operator/(const PolyMatrix &other) const
    {
        if (cols() != other.cols())
        {
            throw std::runtime_error("PolyMatrix::operator/: PolyMatrices must have the same number of columns");
        }

        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows(); i++)
        {
            result.push_back(this->poly_matrix[i]);
        }

        for (size_t i = 0; i < other.rows(); i++)
        {
            result.push_back(other.get_matrix()[i]);
        }

        return result;
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::operator*(const PolyMatrix &other) const
    {
        if (cols() != other.rows())
        {
            throw std::runtime_error("PolyMatrix::operator*: Number of columns in the first PolyMatrix must be equal to the number of rows in the second PolyMatrix");
        }

        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows(); i++)
        {
            std::vector<PolyProxy> currentRow;

            for (size_t j = 0; j < other.cols(); j++)
            {
                poly_type current;

                for (size_t k = 0; k < cols(); k++)
                {
                    current = current + (this->poly_matrix[i].get_vector()[k].get_poly() *
                                         other.transposed().get_matrix()[k].get_vector()[j].get_poly());
                } // TODO: test this

                currentRow.push_back(PolyProxy(current));
            }

            result.push_back(PolyVector(currentRow));
        }

        return PolyMatrix(result);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyVector<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::operator*(const PolyVector<_degree, _coeff_modulus> &other) const
    {
        return other * *this;
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::operator*(const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar) const
    {
        return scale(scalar);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> operator*(const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar, const PolyMatrix<_degree, _coeff_modulus> &poly_matrix)
    {
        return poly_matrix.scale(scalar);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::random_poly_matrix(size_t rows, size_t cols)
    {
        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows; i++)
        {
            result.push_back(PolyVector::random_poly_vector(cols));
        }

        return PolyMatrix(result);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::identity_matrix(size_t size)
    {
        std::vector<PolyVector> result;

        for (size_t i = 0; i < size; i++)
        {
            std::vector<PolyProxy> currentRow;

            for (size_t j = 0; j < size; j++)
            {
                if (i == j)
                {
                    currentRow.push_back(PolyProxy(1));
                }
                else
                {
                    currentRow.push_back(PolyProxy(0));
                }
            }

            result.push_back(PolyVector(currentRow));
        }

        return PolyMatrix(result);
    }

    template <unsigned short _degree, size_t _coeff_modulus>
    PolyMatrix<_degree, _coeff_modulus> PolyMatrix<_degree, _coeff_modulus>::zero_matrix(size_t rows, size_t cols)
    {
        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows; i++)
        {
            std::vector<PolyProxy> currentRow;

            for (size_t j = 0; j < cols; j++)
            {
                currentRow.push_back(PolyProxy(0));
            }

            result.push_back(PolyVector(currentRow));
        }

        return PolyMatrix(result);
    }
} // namespace pqdevkit
