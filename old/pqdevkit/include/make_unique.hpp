#ifndef PQDEVKIT_MAKE_UNIQUE_HPP
#define PQDEVKIT_MAKE_UNIQUE_HPP

#include <memory>
#include <type_traits>

template <typename T, typename... Args>
typename std::enable_if<!std::is_array<T>::value, std::unique_ptr<T>>::type
make_unique(Args &&...args)
{
  return std::unique_ptr<T>(new T(std::forward<Args>(args)...));
}
// Use make_unique for arrays
template <typename T>
typename std::enable_if<std::is_array<T>::value, std::unique_ptr<T>>::type
make_unique(size_t n)
{
  return std::unique_ptr<T>(new typename std::remove_extent<T>::type[n]());
}

#endif // PQDEVKIT_MAKE_UNIQUE_HPP