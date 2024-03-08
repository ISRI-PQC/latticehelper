GMP_PREFIX=/opt/homebrew/opt/gmp
# LIBOQS_ROOT=liboqs
NTL_ROOT=ntl
NFL_ROOT=NFLlib

.PHONY: all clean liboqs ntl nfl python go

all: python go

clean:
	make -C pqdevkit clean
# rm -rf ntl/lib/libntl.a && make -C ntl/src -j $(shell nproc) clean
	

# $(LIBOQS_ROOT):
# 	@if [ ! -d $(LIBOQS_ROOT) ]; then \
#         git clone -b main https://github.com/open-quantum-safe/liboqs.git; \
#     else \
#         echo "Repository already exists in $(LIBOQS_ROOT)"; \
#     fi
	

# $(LIBOQS_ROOT)/build/lib/liboqs.a: $(LIBOQS_ROOT)
# 	cmake \
# 		-GNinja \
# 		-B $(LIBOQS_ROOT)/build \
# 		$(LIBOQS_ROOT)
# 	ninja \
# 		-j $(shell nproc) \
# 		-C $(LIBOQS_ROOT)/build

$(NTL_ROOT):
	@if [ ! -d $(NTL_ROOT) ]; then \
        git clone -b main https://github.com/libntl/ntl.git; \
    else \
        echo "Repository already exists in $(NTL_ROOT)"; \
    fi

$(NTL_ROOT)/lib/libntl.a: $(NTL_ROOT)
	cd $(NTL_ROOT)/src && \
		./configure \
			GMP_PREFIX=$(GMP_PREFIX) \
			NTL_EXCEPTIONS=on \
			TUNE=auto
	make -C $(NTL_ROOT)/src -j $(shell nproc)
	make -C $(NTL_ROOT)/src -j $(shell nproc) check
	mkdir -p $(NTL_ROOT)/lib
	cp -p $(NTL_ROOT)/src/ntl.a $(NTL_ROOT)/lib/libntl.a
	- chmod a+r $(NTL_ROOT)/lib/libntl.a

# $(NFL_ROOT):
# 	@if [ ! -d $(NFL_ROOT) ]; then \
#         git clone https://github.com/Muzosh/NFLlib.git; \
#     else \
#         echo "Repository already exists in $(NFL_ROOT)"; \
#     fi

# $(NFL_ROOT)/lib/libNFLlib.a: $(NFL_ROOT)
# 	mkdir -p $(NFL_ROOT)/_build
# 	cmake -S $(NFL_ROOT) -DCMAKE_BUILD_TYPE=Release -DNFL_OPTIMIZED=ON -DNFLLIB_USE_AVX=ON -B $(NFL_ROOT)/_build
# 	make -C $(NFL_ROOT)/_build
# 	cp $(NFL_ROOT)/_build/libnfllib_static.a $(NFL_ROOT)/lib/libNFLlib.a

pqdevkit: pqdevkit/build/libpqdevkit.a

pqdevkit/build/libpqdevkit.a:
	make -C pqdevkit -j $(shell nproc) build

# python: pqdevkit/build/libpqdevkit.a $(NTL_ROOT)/lib/libNFLlib.a
# 	@echo "Building Python extension.."
# 	rm -rf swigbuild/python
# 	mkdir -p swigbuild/python

# 	@echo "Running swig.."
# 	swig \
# 		-python \
# 		-c++ \
# 		-o ./swigbuild/python/pqdevkit_wrap.cxx \
# 		-I$(NTL_ROOT)/include \
# 		-Ipqdevkit/include \
# 		pqdevkit.i

# 	@echo "Compiling C++ wrapper.."
# 	g++ \
# 		-g -O0 -DDEBUG -Wall -ggdb \
# 		-std=c++11 \
# 		-O2 \
# 		-fPIC \
# 		-I$(NTL_ROOT)/include \
# 		-Ipqdevkit/include \
# 		$(shell	python-config --cflags) \
# 		-c ./swigbuild/python/pqdevkit_wrap.cxx \
# 		-o ./swigbuild/python/pqdevkit_wrap.o

# 	@echo "Building shared library.."
# 	g++ \
# 		-g -O0 -DDEBUG -Wall -ggdb \
# 		-std=c++11 \
# 		-shared \
# 		./swigbuild/python/pqdevkit_wrap.o \
# 		-Lpqdevkit/lib -lpqdevkit \
# 		-L$(NTL_ROOT)/lib -lntl \
# 		-L$(shell python-config --prefix)/lib \
# 		-l$(shell ls $(shell python-config --prefix)/lib | grep -o 'python[0-9]\+\.[0-9]\+' | tail -1) \
# 		-lssl \
# 		-lcrypto \
# 		-o ./swigbuild/python/_pqdevkit.so

# 	# temp
# 	cp ./swigbuild/python/_pqdevkit.so pypqdevkit/_pqdevkit.so
# 	cp ./swigbuild/python/pqdevkit.py pypqdevkit/pqdevkit.py


# go: pqdevkit/build/libpqdevkit.a $(NTL_ROOT)/lib/libNFLlib.a # $(LIBOQS_ROOT)/build/lib/liboqs.a
# 	@echo "Building Go extension.."
# 	rm -rf swigbuild/go
# 	mkdir -p swigbuild/go
# 	swig \
# 		-go \
# 		-cgo \
# 		-intgosize 64 \
# 		-c++ \
# 		-o ./swigbuild/go/oqsgo_wrap.cpp \
# 		-I$(NTL_ROOT) \
# 		-Ipqdevkit/include \
# 		pqdevkit.i
