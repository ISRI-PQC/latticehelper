cmake_minimum_required(VERSION 3.5)
set(CMAKE_CXX_STANDARD 11)

project(pqdevkit CXX)


# Set cache variables
# set(NFL_LIB "${CMAKE_CURRENT_LIST_DIR}/../NFLlib/lib/libNFLlib.a" CACHE PATH "Path to other static libraries")
# set(NFL_INCLUDE "${CMAKE_CURRENT_LIST_DIR}/../NFLlib/include" CACHE PATH "Additional include directories")
set(NTL_LIB "${CMAKE_CURRENT_LIST_DIR}/../ntl/lib" PATH)
set(NTL_INCLUDE "${CMAKE_CURRENT_LIST_DIR}/../ntl/include" PATH)

find_library(NTL_LIBRARY ntl PATHS ${NTL_LIB})

file(GLOB SOURCE_FILES
    src/*.cpp
)
list(REMOVE_ITEM SOURCE_FILES "${CMAKE_CURRENT_SOURCE_DIR}/src/main.cpp")
find_package(Threads)

include_directories(include ${NTL_INCLUDE})

add_library(pqdevkit STATIC ${SOURCE_FILES})
target_compile_options(pqdevkit PUBLIC -g -O0 -DDEBUG -ggdb -Wall -Wextra -Werror -Wno-unused-parameter -pthread -march=native)
target_link_libraries(pqdevkit PUBLIC ${CMAKE_THREAD_LIBS_INIT} ${NTL_LIBRARY} gmp m)

add_executable(main src/main.cpp)
target_link_libraries(main PUBLIC pqdevkit)
target_include_directories(main PUBLIC include ${NTL_INCLUDE})