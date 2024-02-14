#include "poly_vector.hpp"

namespace pqdevkit
{

    // PolyVector
    PolyVector::PolyVector(std::initializer_list<std::initializer_list<coeff_type>> poly_vector)
    {
        this->poly_vector_ptr = std::make_unique<std::vector<PolyProxy>>(poly_vector);
    }

    PolyVector::PolyVector(const std::vector<PolyProxy> &poly_vector)
    {
        this->poly_vector_ptr = std::make_unique<std::vector<PolyProxy>>(poly_vector);
    }

    PolyVector::~PolyVector()
    {
    }

    coeff_type PolyVector::infinite_norm() const
    {
        coeff_type maxNorm = std::numeric_limits<coeff_type>::min();

        for (const auto &polyProxy : *poly_vector_ptr)
        {
            coeff_type currentNorm = polyProxy.infinite_norm();
            if (currentNorm > maxNorm)
            {
                maxNorm = currentNorm;
            }
        }

        return maxNorm;
    }

    std::vector<coeff_type> PolyVector::listize() const
    {
        std::vector<coeff_type> result;

        for (const auto &polyProxy : *poly_vector_ptr)
        {
            std::vector<coeff_type> currentList = polyProxy.listize();
            result.insert(result.end(), currentList.begin(), currentList.end());
        }

        return result;
    }

    PolyVector PolyVector::scale(const coeff_type &scalar) const
    {
        std::vector<PolyProxy> scaledPolyVector;

        for (const auto &polyProxy : *poly_vector_ptr)
        {
            scaledPolyVector.push_back(polyProxy * scalar);
        }

        return PolyVector(scaledPolyVector);
    }

    PolyVector PolyVector::scale(const poly_type &poly) const
    {
        std::vector<PolyProxy> scaledPolyVector;

        for (const auto &polyProxy : *poly_vector_ptr)
        {
            scaledPolyVector.push_back(polyProxy * poly);
        }

        return PolyVector(scaledPolyVector);
    }

    PolyVector PolyVector::operator+(const PolyVector &other) const
    {
        if (poly_vector_ptr->size() != other.poly_vector_ptr->size())
        {
            throw std::runtime_error("PolyVector::operator+: PolyVectors must have the same length");
        }

        std::vector<PolyProxy> result;

        for (size_t i = 0; i < poly_vector_ptr->size(); i++)
        {
            result.push_back((*poly_vector_ptr)[i] + (*other.poly_vector_ptr)[i]);
        }

        return PolyVector(result);
    }

    PolyVector PolyVector::operator-(const PolyVector &other) const
    {
        if (poly_vector_ptr->size() != other.poly_vector_ptr->size())
        {
            throw std::runtime_error("PolyVector::operator-: PolyVectors must have the same length");
        }

        std::vector<PolyProxy> result;

        for (size_t i = 0; i < poly_vector_ptr->size(); i++)
        {
            result.push_back((*poly_vector_ptr)[i] - (*other.poly_vector_ptr)[i]);
        }

        return PolyVector(result);
    }

    PolyVector PolyVector::operator|(const PolyVector &other) const
    {
        if (poly_vector_ptr->size() != other.poly_vector_ptr->size())
        {
            throw std::runtime_error("PolyVector::operator|: PolyVectors must have the same length");
        }

        // concatenate this and other
        std::vector<PolyProxy> result;

        for (const auto &polyProxy : *poly_vector_ptr)
        {
            result.push_back(polyProxy);
        }

        for (const auto &polyProxy : *other.poly_vector_ptr)
        {
            result.push_back(polyProxy);
        }

        return PolyVector(result);
    }

    PolyProxy PolyVector::operator*(const PolyVector &other) const
    {
        // dot product
        if (poly_vector_ptr->size() != other.poly_vector_ptr->size())
        {
            throw std::runtime_error("PolyVector::operator*: PolyVectors must have the same length");
        }

        PolyProxy result = PolyProxy(0);

        for (size_t i = 0; i < poly_vector_ptr->size(); i++)
        {
            result = result + ((*poly_vector_ptr)[i] * (*other.poly_vector_ptr)[i]);
        }

        return result;
    }

    PolyVector PolyVector::operator*(const PolyMatrix &other) const
    {
        if (poly_vector_ptr->size() != other.poly_matrix_ptr->size())
        {
            throw std::runtime_error("PolyVector::operator*: PolyVector and PolyMatrix must have the same length");
        }

        std::vector<PolyProxy> result;

        for (size_t i = 0; i < poly_vector_ptr->size(); i++)
        {
            result.push_back((*poly_vector_ptr)[i] * (*other.poly_matrix_ptr)[i]);
        }

        return PolyVector(result);
    }

    PolyVector PolyVector::operator*(const coeff_type &scalar) const
    {
        return this->scale(scalar);
    }

    PolyVector PolyVector::random_poly_vector(size_t length)
    {
        std::vector<PolyProxy> result;

        for (size_t i = 0; i < length; i++)
        {
            result.push_back(PolyProxy::random_poly());
        }

        return PolyVector(result);
    }

    PolyVector operator*(const coeff_type &scalar, const PolyVector &poly_vector)
    {
        return poly_vector.scale(scalar);
    }

} // namespace pqdevkit
