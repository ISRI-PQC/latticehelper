#include "pqdevkit.hpp"
#include <iostream>
using namespace std;
int main()
{
  auto degree = 4;
  auto modulus = 3;
  init_ntl(degree, modulus);

  pqdevkit::PolyProxy p({1,0,0,1}); // x^3 + 1
  cout << p.get_poly() << endl;
  pqdevkit::PolyProxy q({2,0,1,2}); // 2x^3 + x^2 + 2
  cout << q.get_poly() << endl;
  pqdevkit::PolyProxy r = p + q; // 3x^3 + x^2 + 3
  cout << r.get_poly() << endl;
  pqdevkit::PolyProxy t = p - q; // -1x^3 - x^2 - 1
  cout << t.get_poly() << endl;
  pqdevkit::PolyProxy s = p * q; // 2x^6 + x^5 + 4x^3 + x^2 + 2 mod x^4 + 1 = 4x^3 - x^2 - x + 2
  cout << s.get_poly() << endl;
  return 0;
}