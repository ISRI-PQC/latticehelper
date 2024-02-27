#include "poly_vector.hpp"
#include "poly_matrix.hpp"

namespace pqdevkit
{
    // PolyVector
    PolyVector::PolyVector(const std::initializer_list<std::initializer_list<coeff_type>> other)
    {
        std::vector<PolyProxy> result;

        for (auto poly : other)
        {
            result.push_back(PolyProxy(poly));
        }

        this->poly_vector = result;
    }

    PolyVector::PolyVector(const std::vector<PolyProxy> &other)
    {
        this->poly_vector = other;
    }

    PolyVector::PolyVector(const PolyVector &other)
    {
        this->poly_vector = other.poly_vector;
    }

    PolyVector::~PolyVector() {}

    const std::vector<PolyProxy> &PolyVector::get_vector() const
    {
        return this->poly_vector;
    }

    coeff_type PolyVector::infinite_norm() const
    {
        coeff_type maxNorm = std::numeric_limits<coeff_type>::min();

        for (const auto &polyProxy : this->poly_vector)
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
        return this->poly_vector.size();
    }

    std::vector<coeff_type> PolyVector::listize() const
    {
        std::vector<coeff_type> result;

        for (const auto &polyProxy : this->poly_vector)
        {
            std::vector<coeff_type> currentList = polyProxy.listize();
            result.insert(result.end(), currentList.begin(), currentList.end());
        }

        return result;
    }

    PolyVector PolyVector::scale(const coeff_type &scalar) const
    {
        std::vector<PolyProxy> scaledPolyVector;

        for (const auto &polyProxy : this->poly_vector)
        {
            scaledPolyVector.push_back(polyProxy * scalar);
        }

        return scaledPolyVector;
    }

    PolyVector PolyVector::scale(const poly_type &poly) const
    {
       std::vector<PolyProxy> scaledPolyVector;

        for (const auto &polyProxy : this->poly_vector)
        {
            scaledPolyVector.push_back(polyProxy * poly);
        }

        return scaledPolyVector;
    }

    PolyVector PolyVector::operator+(const PolyVector &other) const
    {
        if (this->poly_vector.size() != other.poly_vector.size())
        {
            throw std::runtime_error("PolyVector::operator+: PolyVectors must have the same length");
        }

        std::vector<PolyProxy> result;

        for (size_t i = 0; i < this->poly_vector.size(); i++)
        {
            result.push_back(
                this->poly_vector[i] + other.poly_vector[i]);
        }

        return result;
    }

    PolyVector PolyVector::operator-(const PolyVector &other) const
    {
        if (this->poly_vector.size() != other.poly_vector.size())
        {
            throw std::runtime_error("PolyVector::operator-: PolyVectors must have the same length");
        }

         std::vector<PolyProxy> result;

        for (size_t i = 0; i < this->poly_vector.size(); i++)
        {
            result.push_back(this->poly_vector[i] - other.poly_vector[i]);
        }

        return PolyVector(result);
    }

    PolyVector PolyVector::operator|(const PolyVector &other) const
    {
        if (this->poly_vector.size() != other.poly_vector.size())
        {
            throw std::runtime_error("PolyVector::operator|: PolyVectors must have the same length");
        }

        // concatenate this and other
        std::vector<PolyProxy> result;

        for (const auto &polyProxy : this->poly_vector)
        {
            result.push_back(polyProxy);
        }

        for (const auto &polyProxy : other.poly_vector)
        {
            result.push_back(polyProxy);
        }

        return PolyVector(result);
    }

    PolyProxy PolyVector::operator*(const PolyVector &other) const
    {
        // dot product
        if (this->poly_vector.size() != other.poly_vector.size())
        {
            throw std::runtime_error("PolyVector::operator*: PolyVectors must have the same length");
        }

        poly_type result;

        for (size_t i = 0; i < this->poly_vector.size(); i++)
        {
            result = result + (this->poly_vector[i].get_poly() * other.poly_vector[i].get_poly());
        }

        return result;
    }

    PolyVector PolyVector::operator*(const PolyMatrix &other) const
    {
        if (this->poly_vector.size() != other.cols())
        {
            throw std::runtime_error("PolyVector::operator*: PolyVector and PolyMatrix must have the same length");
        }

        std::vector<PolyProxy> result;

        for (size_t i = 0; i < other.rows(); i++)
        {
            poly_type current;

            for (size_t j = 0; j < this->poly_vector.size(); j++)
            {
                current = current + (this->poly_vector[j].get_poly() *
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
