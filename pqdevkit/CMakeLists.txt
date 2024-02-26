cmake_minimum_required(VERSION 3.5)

project(pqdevkit CXX)

# Set cache variables
set(NFL_LIB "${CMAKE_CURRENT_LIST_DIR}/../NFLlib/lib/libNFLlib.a" CACHE PATH "Path to other static libraries")
set(NFL_INCLUDE "${CMAKE_CURRENT_LIST_DIR}/../NFLlib/include" CACHE PATH "Additional include directories")

file(GLOB SOURCE_FILES
    src/*.cpp
)

include_directories(include)

add_library(pqdevkit STATIC ${SOURCE_FILES})

target_compile_options(pqdevkit PRIVATE -std=c++11 -Wall -Wextra -Werror -Wno-unused-parameter)

# Link other static libraries
# target_link_libraries(pqdevkit PRIVATE ${NFL_LIB})

# Include additional directories
target_include_directories(pqdevkit PRIVATE ${NFL_INCLUDE})

add_executable(pqdevkit_main src/main.cpp ${SOURCE_FILES})
target_compile_options(pqdevkit_main PRIVATE -std=c++11 -Wall -Wextra -Werror -Wno-unused-parameter)

# Link other static libraries
target_link_libraries(pqdevkit_main PRIVATE ${NFL_LIB})

# Include additional directories
target_include_directories(pqdevkit_main PRIVATE ${NFL_INCLUDE})