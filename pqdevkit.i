%module(directors = "1") pqdevkit

%include <std_unique_ptr.i>
%include "std_vector.i"
%include "cpointer.i"

%feature("director");

%inline %{
    #define SWIG_FILE_WITH_INIT
    #include "pqdevkit.hpp"
%}

%include "poly_proxy.hpp"
%include "poly_vector.hpp"
%include "poly_matrix.hpp"
