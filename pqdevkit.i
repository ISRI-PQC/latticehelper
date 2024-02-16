%module(directors = "1") pqdevkit

%inline %{
    #define SWIG_FILE_WITH_INIT
    #include "pqdevkit.hpp"
%}

%include "poly_proxy.hpp"
%include "poly_vector.hpp"
%include "poly_matrix.hpp"

