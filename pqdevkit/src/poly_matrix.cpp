#include "poly_matrix.hpp"

namespace pqdevkit
{
    PolyMatrix::PolyMatrix(std::initializer_list<std::initializer_list<std::initializer_list<coeff_type>>> poly_matrix)
    {
        std::vector<PolyVector> result;

        for (auto row : poly_matrix)
        {
            result.push_back(PolyVector(row));
        }

        poly_matrix_ptr = std::make_unique<std::vector<PolyVector>>(result);
    }

    PolyMatrix::PolyMatrix(const std::vector<PolyVector> &poly_matrix)
    {
        poly_matrix_ptr = std::make_unique<std::vector<PolyVector>>(poly_matrix);
    }

    PolyMatrix::~PolyMatrix() {}

    std::vector<PolyVector> PolyMatrix::get_matrix() const
    {
        return *poly_matrix_ptr;
    }

    size_t PolyMatrix::rows() const
    {
        return poly_matrix_ptr->.size();
    }

    size_t PolyMatrix::cols() const
    {
        return poly_matrix_ptr->at(0).length();
    }

    coeff_type PolyMatrix::infinite_norm() const
    {
        coeff_type maxNorm = std::numeric_limits<coeff_type>::min();

        for (const auto &polyVector : *poly_matrix_ptr)
        {
            coeff_type currentNorm = polyVector.infinite_norm();
            if (currentNorm > maxNorm)
            {
                maxNorm = currentNorm;
            }
        }

        return maxNorm;
    }

    std::vector<coeff_type> PolyMatrix::listize() const
    {
        std::vector<coeff_type> result;

        for (const auto &polyVector : *poly_matrix_ptr)
        {
            std::vector<coeff_type> currentList = polyVector.listize();
            result.insert(result.end(), currentList.begin(), currentList.end());
        }

        return result;
    }

    PolyMatrix PolyMatrix::transposed() const
    {
        std::vector<PolyVector> result;

        for (size_t i = 0; i < cols(); i++)
        {
            std::vector<PolyProxy> currentColumn;

            for (size_t j = 0; j < rows(); j++)
            {
                currentColumn.push_back(poly_matrix_ptr->at(j).get_vector()[i]);
            }

            result.push_back(PolyVector(currentColumn));
        }

        return PolyMatrix(result);
    }

    PolyMatrix PolyMatrix::scale(const coeff_type &scalar) const
    {
        std::vector<PolyVector> scaledPolyMatrix;

        for (const auto &polyVector : *poly_matrix_ptr)
        {
            scaledPolyMatrix.push_back(polyVector.scale(scalar));
        }

        return PolyMatrix(scaledPolyMatrix);
    }

    PolyMatrix PolyMatrix::scale(const poly_type &poly) const
    {
        std::vector<PolyVector> scaledPolyMatrix;

        for (const auto &polyVector : *poly_matrix_ptr)
        {
            scaledPolyMatrix.push_back(polyVector.scale(poly));
        }

        return PolyMatrix(scaledPolyMatrix);
    }

    PolyMatrix PolyMatrix::operator+(const PolyMatrix &other) const
    {
        if (rows() != other.rows() || cols() != other.cols())
        {
            throw std::runtime_error("PolyMatrix::operator+: PolyMatrices must have the same dimensions");
        }

        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows(); i++)
        {
            result.push_back(poly_matrix_ptr->at(i) + other.poly_matrix_ptr->at(i));
        }

        return PolyMatrix(result);
    }

    PolyMatrix PolyMatrix::operator-(const PolyMatrix &other) const
    {
        if (rows() != other.rows() || cols() != other.cols())
        {
            throw std::runtime_error("PolyMatrix::operator-: PolyMatrices must have the same dimensions");
        }

        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows(); i++)
        {
            result.push_back(poly_matrix_ptr->at(i) - other.poly_matrix_ptr->at(i));
        }

        return PolyMatrix(result);
    }

    PolyMatrix PolyMatrix::operator|(const PolyMatrix &other) const
    {
        if (rows() != other.rows())
        {
            throw std::runtime_error("PolyMatrix::operator|: PolyMatrices must have the same number of rows");
        }

        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows(); i++)
        {
            result.push_back(poly_matrix_ptr->at(i) | other.poly_matrix_ptr->at(i));
        }

        return PolyMatrix(result);
    }

    PolyMatrix PolyMatrix::operator/(const PolyMatrix &other) const
    {
        if (cols() != other.cols())
        {
            throw std::runtime_error("PolyMatrix::operator/: PolyMatrices must have the same number of columns");
        }

        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows(); i++)
        {
            result.push_back(poly_matrix_ptr->at(i));
        }

        for (size_t i = 0; i < other.rows(); i++)
        {
            result.push_back(other.get_matrix()[i]);
        }

        return result;
    }

    PolyMatrix PolyMatrix::operator*(const PolyMatrix &other) const
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
                    current = current + (poly_matrix_ptr->at(i).get_vector()[k].get_poly() *
                                         other.transposed().get_matrix()[k].get_vector()[j].get_poly());
                } // TODO: test this

                currentRow.push_back(PolyProxy(current));
            }

            result.push_back(PolyVector(currentRow));
        }

        return PolyMatrix(result);
    }

    PolyVector PolyMatrix::operator*(const PolyVector &other) const
    {
        return other * *this;
    }

    PolyMatrix PolyMatrix::operator*(const coeff_type &scalar) const
    {
        return scale(scalar);
    }

    PolyMatrix operator*(const coeff_type &scalar, const PolyMatrix &poly_matrix)
    {
        return poly_matrix.scale(scalar);
    }

    PolyMatrix PolyMatrix::random_poly_matrix(size_t rows, size_t cols)
    {
        std::vector<PolyVector> result;

        for (size_t i = 0; i < rows; i++)
        {
            result.push_back(PolyVector::random_poly_vector(cols));
        }

        return PolyMatrix(result);
    }

    PolyMatrix PolyMatrix::identity_matrix(size_t size)
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

    PolyMatrix PolyMatrix::zero_matrix(size_t rows, size_t cols)
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
