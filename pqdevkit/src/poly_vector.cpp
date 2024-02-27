#include "poly_vector.hpp"
#include "poly_matrix.hpp"

namespace pqdevkit
{
    // PolyVector
    PolyVector::PolyVector(const std::initializer_list<std::initializer_list<coeff_type>> other)
    {
        std::vector<PolyProxy> result(other.size());

        size_t i = 0;
        for (auto row : other)
        {
            result[i] = PolyProxy(row);
            i++;
        }

        this->poly_vector = std::move(result);
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
        size_t total_size = this->poly_vector.size() * this->poly_vector.front().degree;
        std::vector<coeff_type> result(total_size);

        size_t i = 0;
        for (const auto &polyProxy : this->poly_vector)
        {
            for (const auto &coeff : polyProxy.listize())
            {
                result[i] = coeff;
                i++;
            }
        }

        return result;
    }

    PolyVector PolyVector::scale(const coeff_type &scalar) const
    {
        std::vector<PolyProxy> scaledPolyVector(this->poly_vector.size());

        for (size_t i = 0; i < this->poly_vector.size(); i++)
        {
            scaledPolyVector[i] = this->poly_vector[i] * scalar;
        }

        return scaledPolyVector;
    }

    PolyVector PolyVector::scale(const poly_type &poly) const
    {
        std::vector<PolyProxy> scaledPolyVector(this->poly_vector.size());

        for (size_t i = 0; i < this->poly_vector.size(); i++)
        {
            scaledPolyVector[i] = this->poly_vector[i] * poly;
        }

        return scaledPolyVector;
    }

    PolyVector PolyVector::operator+(const PolyVector &other) const
    {
        if (this->poly_vector.size() != other.poly_vector.size())
        {
            throw std::runtime_error("PolyVector::operator+: PolyVectors must have the same length");
        }

        std::vector<PolyProxy> result(this->poly_vector.size());

        for (size_t i = 0; i < this->poly_vector.size(); i++)
        {
            result[i] = this->poly_vector[i] + other.poly_vector[i];
        }

        return result;
    }

    PolyVector PolyVector::operator-(const PolyVector &other) const
    {
        if (this->poly_vector.size() != other.poly_vector.size())
        {
            throw std::runtime_error("PolyVector::operator-: PolyVectors must have the same length");
        }

        std::vector<PolyProxy> result(this->poly_vector.size());

        for (size_t i = 0; i < this->poly_vector.size(); i++)
        {
            result[i] = this->poly_vector[i] - other.poly_vector[i];
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
        std::vector<PolyProxy> result(this->poly_vector.size() + other.poly_vector.size());

        for (size_t i = 0; i < this->poly_vector.size(); i++)
        {
            result[i] = this->poly_vector[i];
        }

        for (size_t i = 0; i < other.poly_vector.size(); i++)
        {
            result[i + this->poly_vector.size()] = other.poly_vector[i];
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

        std::vector<PolyProxy> result(other.rows());

        for (size_t i = 0; i < other.rows(); i++)
        {
            poly_type current;

            for (size_t j = 0; j < this->poly_vector.size(); j++)
            {
                current = current + (this->poly_vector[j].get_poly() *
                                     other.get_matrix()[i].get_vector()[j].get_poly());
            }

            result[i] = PolyProxy(current);
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
        std::vector<PolyProxy> result(length);

        for (size_t i = 0; i < length; i++)
        {
            result[i] = PolyProxy::random_poly();
        }

        return PolyVector(result);
    }

} // namespace pqdevkit
