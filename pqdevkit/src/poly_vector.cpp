#include "poly_vector.hpp"
#include "poly_matrix.hpp"

// TODO: check algorithms according to ntl/src/mat_ZZ.cpp
namespace pqdevkit
{
  PolyVector::PolyVector(
    const std::initializer_list<std::initializer_list<long>> other)
  {
    std::vector<PolyProxy> result;

    for(auto poly : other)
      {
        result.emplace_back(PolyProxy(poly));
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

  long PolyVector::infinite_norm() const
  {
    long maxNorm = std::numeric_limits<long>::min();

    for(const auto &polyProxy : this->poly_vector)
      {
        long currentNorm = polyProxy.infinite_norm();
        if(currentNorm > maxNorm)
          {
            maxNorm = currentNorm;
          }
      }

    return maxNorm;
  }

  size_t PolyVector::length() const { return this->poly_vector.size(); }

  std::vector<PQDEVKIT_COEFF_TYPE> PolyVector::listize() const
  {
    std::vector<PQDEVKIT_COEFF_TYPE> result;

    for(const auto &polyProxy : this->poly_vector)
      {
        std::vector<PQDEVKIT_COEFF_TYPE> currentList = polyProxy.listize();
        result.insert(result.end(), currentList.begin(), currentList.end());
      }

    return result;
  }

  PolyVector PolyVector::scale(const PQDEVKIT_COEFF_TYPE &scalar) const
  {
    std::vector<PolyProxy> scaledPolyVector;

    for(const auto &polyProxy : this->poly_vector)
      {
        scaledPolyVector.push_back(polyProxy * scalar);
      }

    return scaledPolyVector;
  }

  PolyVector PolyVector::scale(const PQDEVKIT_POLY_TYPE &poly) const
  {
    std::vector<PolyProxy> scaledPolyVector;

    for(const auto &polyProxy : this->poly_vector)
      {
        scaledPolyVector.push_back(polyProxy * poly);
      }

    return scaledPolyVector;
  }

  PolyVector PolyVector::operator+(const PolyVector &other) const
  {
    if(this->poly_vector.size() != other.poly_vector.size())
      {
        throw std::runtime_error(
          "PolyVector::operator+: PolyVectors must have the same length");
      }

    std::vector<PolyProxy> result;

    for(size_t i = 0; i < this->poly_vector.size(); i++)
      {
        result.push_back(this->poly_vector[i] + other.poly_vector[i]);
      }

    return result;
  }

  PolyVector PolyVector::operator-(const PolyVector &other) const
  {
    if(this->poly_vector.size() != other.poly_vector.size())
      {
        throw std::runtime_error(
          "PolyVector::operator-: PolyVectors must have the same length");
      }

    std::vector<PolyProxy> result;

    for(size_t i = 0; i < this->poly_vector.size(); i++)
      {
        result.push_back(this->poly_vector[i] - other.poly_vector[i]);
      }

    return PolyVector(result);
  }

  PolyVector PolyVector::operator|(const PolyVector &other) const
  {
    if(this->poly_vector.size() != other.poly_vector.size())
      {
        throw std::runtime_error(
          "PolyVector::operator|: PolyVectors must have the same length");
      }

    // concatenate this and other
    std::vector<PolyProxy> result;

    for(const auto &polyProxy : this->poly_vector)
      {
        result.push_back(polyProxy);
      }

    for(const auto &polyProxy : other.poly_vector)
      {
        result.push_back(polyProxy);
      }

    return PolyVector(result);
  }

  PolyProxy PolyVector::operator*(const PolyVector &other) const
  {
    // dot product
    if(this->poly_vector.size() != other.poly_vector.size())
      {
        throw std::runtime_error(
          "PolyVector::operator*: PolyVectors must have the same length");
      }

    PQDEVKIT_POLY_TYPE result;

    for(size_t i = 0; i < this->poly_vector.size(); i++)
      {
        result = result
                 + (this->poly_vector[i].get_poly()
                    * other.poly_vector[i].get_poly());
      }

    return result;
  }

  PolyVector PolyVector::operator*(const PolyMatrix &other) const
  {
    if(this->poly_vector.size() != other.cols())
      {
        throw std::runtime_error("PolyVector::operator*: PolyVector and "
                                 "PolyMatrix must have the same length");
      }

    std::vector<PolyProxy> result;

    for(size_t i = 0; i < other.rows(); i++)
      {
        PQDEVKIT_POLY_TYPE current;

        for(size_t j = 0; j < this->poly_vector.size(); j++)
          {
            current = current
                      + (this->poly_vector[j].get_poly()
                         * other.get_matrix()[i].get_vector()[j].get_poly());
          }

        result.push_back(PolyProxy(current));
      }

    return PolyVector(result);
  }

  PolyVector PolyVector::operator*(const PQDEVKIT_COEFF_TYPE &scalar) const
  {
    return scale(scalar);
  }

  PolyVector
  operator*(const PQDEVKIT_COEFF_TYPE &scalar, const PolyVector &poly_vector)
  {
    return poly_vector.scale(scalar);
  }

  PolyVector PolyVector::random_poly_vector(size_t length)
  {
    std::vector<PolyProxy> result;

    for(size_t i = 0; i < length; i++)
      {
        result.push_back(PolyProxy::random_poly());
      }

    return PolyVector(result);
  }

} // namespace pqdevkit
