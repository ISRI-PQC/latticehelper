#!/bin/bash

rm -rf swigbuild
mkdir -p swigbuild

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

GMP_PREFIX=/opt/homebrew/opt/gmp
LIBOQS_ROOT=$script_dir/liboqs
NTL_ROOT=$script_dir/ntl


# Ensure libOQS is built
if [[ -e "$LIBOQS_ROOT" ]]; then
    echo "liboqs directory already exists, skipping cloning"; \
else \
    git clone -b main https://github.com/open-quantum-safe/liboqs.git; \
fi

if [[ -e "$LIBOQS_ROOT/build/lib/liboqs.a" ]]; then
    echo "liboqs library already builded, skipping compilation"; \
else \
    cmake -GNinja -B $LIBOQS_ROOT/build $LIBOQS_ROOT && ninja -j $(nproc) -C $LIBOQS_ROOT/build; \
fi

# Ensure NTL is built
if [[ -e "$NTL_ROOT" ]]; then
    echo "ntl directory already exists, skipping cloning"; \
else \
    git clone -b main https://github.com/libntl/ntl.git; \
fi

if [[ -e "$NTL_ROOT/lib/libntl.a" ]]; then
    echo "ntl library already builded, skipping compilation"; \
else \
    (
        cd ntl/src \
        && ./configure GMP_PREFIX=$GMP_PREFIX \
        && make \
        && mkdir -p ../lib \
        && cp ntl.a ../lib/libntl.a
    ); \
fi

if [ "$1" = "python" ]; then
    echo "Building Python extension.."; \

    rm -rf swigbuild/python; \
    mkdir -p swigbuild/python; \

    # Compile the C++ wrapper
    swig -python -c++ -o ./swigbuild/python/pqdevkit_wrap.cxx -I$LIBOQS_ROOT/build/include pqdevkit.i && \

    # Compile the C++ wrapper and link it with liboqs
    g++ -std=c++20 -O2 -fPIC -I$LIBOQS_ROOT/build/include -I$NTL_ROOT/include $(python-config --cflags) -c ./swigbuild/python/pqdevkit_wrap.cxx -o ./swigbuild/python/pqdevkit_wrap.o && \
    # Manual working version: g++ -std=c++20 -O2 -fPIC -I$LIBOQS_ROOT/build/include -I$NTL_ROOT/include -I/Users/petr/.pyenv/versions/3.11.5/include/python3.11 -c ./swigbuild/python/pqdevkit_wrap.cxx -o ./swigbuild/python/pqdevkit_wrap.o

    # Link the C++ wrapper with liboqs, libntl and Python+OpenSSL+Intl
    g++ -std=c++20 -shared ./swigbuild/python/pqdevkit_wrap.o -L$LIBOQS_ROOT/build/lib -loqs -L$NTL_ROOT/lib -lntl -L$(python-config --prefix)/lib -l$(ls $(python-config --prefix)/lib | grep -o 'python[0-9]\+\.[0-9]\+' | tail -1) -lssl -lcrypto -o ./swigbuild/python/pqdevkit.so; \
    # Manual working version: g++ -std=c++20 -shared ./swigbuild/python/pqdevkit_wrap.o -L$LIBOQS_ROOT/build/lib -loqs -L$NTL_ROOT/lib -lntl -L/Users/petr/.pyenv/versions/3.11.5/lib -lpython3.11 -L/opt/homebrew/lib -lintl -L/opt/homebrew/opt/openssl@1.1/lib -lssl -lcrypto -o ./swigbuild/python/pqdevkit.so

elif [ "$1" = "go" ]; then
    echo "Building Go extension.."; \

    rm -rf swigbuild/go; \
    mkdir -p swigbuild/go; \
    
    swig -go -cgo -intgosize 64 -c++ -o ./swigbuild/go/oqsgo_wrap.cpp -I$LIBOQS_ROOT/build/include -I$LIBNTL_ROOT pqdevkit.i
else
    echo "Usage: build.sh [python|go]"; \
    exit 1; \
fi