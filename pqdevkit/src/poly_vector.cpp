#include "poly_vector.hpp"

namespace pqdevkit
{
  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::PolyVector(
    const std::initializer_list<std::initializer_list<
      typename PolyProxy<_degree, _coeff_modulus>::coeff_type>>
      other)
  {
    std::vector<PolyProxy> result;

    for(auto poly : other)
      {
        result.push_back(PolyProxy(poly));
      }

    this->poly_vector = result;
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::PolyVector(
    const std::vector<PolyProxy<_degree, _coeff_modulus>> &other)
  {
    this->poly_vector = other;
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::PolyVector(const PolyVector &other)
  {
    this->poly_vector = other.poly_vector;
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::~PolyVector()
  {}

  template <unsigned short _degree, size_t _coeff_modulus>
  const std::vector<PolyProxy<_degree, _coeff_modulus>> &
  PolyVector<_degree, _coeff_modulus>::get_vector() const
  {
    return this->poly_vector;
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  typename PolyProxy<_degree, _coeff_modulus>::coeff_type
  PolyVector<_degree, _coeff_modulus>::infinite_norm() const
  {
    coeff_type maxNorm = std::numeric_limits<coeff_type>::min();

    for(const auto &polyProxy : this->poly_vector)
      {
        coeff_type currentNorm = polyProxy.infinite_norm();
        if(currentNorm > maxNorm)
          {
            maxNorm = currentNorm;
          }
      }

    return maxNorm;
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  size_t PolyVector<_degree, _coeff_modulus>::length() const
  {
    return this->poly_vector.size();
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  std::vector<typename PolyProxy<_degree, _coeff_modulus>::coeff_type>
  PolyVector<_degree, _coeff_modulus>::listize() const
  {
    std::vector<coeff_type> result;

    for(const auto &polyProxy : this->poly_vector)
      {
        std::vector<coeff_type> currentList = polyProxy.listize();
        result.insert(result.end(), currentList.begin(), currentList.end());
      }

    return result;
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::scale(
    const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar) const
  {
    std::vector<PolyProxy> scaledPolyVector;

    for(const auto &polyProxy : this->poly_vector)
      {
        scaledPolyVector.push_back(polyProxy * scalar);
      }

    return scaledPolyVector;
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::scale(
    const typename PolyProxy<_degree, _coeff_modulus>::poly_type &poly) const
  {
    std::vector<PolyProxy> scaledPolyVector;

    for(const auto &polyProxy : this->poly_vector)
      {
        scaledPolyVector.push_back(polyProxy * poly);
      }

    return scaledPolyVector;
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::operator+(const PolyVector &other) const
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

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::operator-(const PolyVector &other) const
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

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::operator|(const PolyVector &other) const
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

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyProxy<_degree, _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::operator*(const PolyVector &other) const
  {
    // dot product
    if(this->poly_vector.size() != other.poly_vector.size())
      {
        throw std::runtime_error(
          "PolyVector::operator*: PolyVectors must have the same length");
      }

    poly_type result;

    for(size_t i = 0; i < this->poly_vector.size(); i++)
      {
        result = result
                 + (this->poly_vector[i].get_poly()
                    * other.poly_vector[i].get_poly());
      }

    return result;
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::operator*(
    const PolyMatrix<_degree, _coeff_modulus> &other) const
  {
    if(this->poly_vector.size() != other.cols())
      {
        throw std::runtime_error("PolyVector::operator*: PolyVector and "
                                 "PolyMatrix must have the same length");
      }

    std::vector<PolyProxy> result;

    for(size_t i = 0; i < other.rows(); i++)
      {
        poly_type current;

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

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::operator*(
    const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar) const
  {
    return scale(scalar);
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus> operator*(
    const typename PolyProxy<_degree, _coeff_modulus>::coeff_type &scalar,
    const PolyVector<_degree, _coeff_modulus> &poly_vector)
  {
    return poly_vector.scale(scalar);
  }

  template <unsigned short _degree, size_t _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>
  PolyVector<_degree, _coeff_modulus>::random_poly_vector(size_t length)
  {
    std::vector<PolyProxy> result;

    for(size_t i = 0; i < length; i++)
      {
        result.push_back(PolyProxy::random_poly());
      }

    return PolyVector(result);
  }

} // namespace pqdevkit
