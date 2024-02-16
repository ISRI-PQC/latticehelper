#include "poly_vector.hpp"
#include "poly_matrix.hpp"

namespace pqdevkit
{
    // PolyVector
    PolyVector::PolyVector(std::initializer_list<std::initializer_list<coeff_type>> poly_vector)
    {
        std::vector<PolyProxy> result;

        for (auto poly : poly_vector)
        {
            result.push_back(PolyProxy(poly));
        }

        poly_vector_ptr = std::unique_ptr<std::vector<PolyProxy>>(new std::vector<PolyProxy>(result));
    }

    PolyVector::PolyVector(const std::vector<PolyProxy> &poly_vector)
    {
        poly_vector_ptr = std::unique_ptr<std::vector<PolyProxy>>(new std::vector<PolyProxy>(poly_vector));
    }

    PolyVector::PolyVector(const PolyVector &other)
    {
        poly_vector_ptr = std::unique_ptr<std::vector<PolyProxy>>(new std::vector<PolyProxy>(*other.poly_vector_ptr));
    }

    PolyVector::~PolyVector() {}

    std::vector<PolyProxy> &PolyVector::get_vector() const
    {
        return *poly_vector_ptr;
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

    size_t PolyVector::length() const
    {
        return (*poly_vector_ptr).size();
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

        return scaledPolyVector;
    }

    PolyVector PolyVector::scale(const poly_type &poly) const
    {
        std::vector<PolyProxy> scaledPolyVector;

        for (const auto &polyProxy : *poly_vector_ptr)
        {
            scaledPolyVector.push_back(polyProxy * poly);
        }

        return scaledPolyVector;
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
            result.push_back(
                poly_vector_ptr->at(i) + other.poly_vector_ptr->at(i));
        }

        return result;
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
            result.push_back(poly_vector_ptr->at(i) - other.poly_vector_ptr->at(i));
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

        poly_type result;

        for (size_t i = 0; i < poly_vector_ptr->size(); i++)
        {
            result = result + (poly_vector_ptr->at(i).get_poly() * other.poly_vector_ptr->at(i).get_poly());
        }

        return result;
    }

    PolyVector PolyVector::operator*(const PolyMatrix &other) const
    {
        if (poly_vector_ptr->size() != other.cols())
        {
            throw std::runtime_error("PolyVector::operator*: PolyVector and PolyMatrix must have the same length");
        }

        std::vector<PolyProxy> result;

        for (size_t i = 0; i < other.rows(); i++)
        {
            poly_type current;

            for (size_t j = 0; j < poly_vector_ptr->size(); j++)
            {
                current = current + (poly_vector_ptr->at(j).get_poly() *
                                     other.get_matrix()[i].get_vector()[j].get_poly());
            }

            result.push_back(PolyProxy(current));
        }

        return PolyVector(result);
    }

    PolyVector PolyVector::operator*(const coeff_type &scalar) const
    {
        return scale(scalar);
    }

    PolyVector operator*(const coeff_type &scalar, const PolyVector &poly_vector)
    {
        return poly_vector.scale(scalar);
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

} // namespace pqdevkit
