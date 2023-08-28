cmake_minimum_required(VERSION 3.10)
project(golib)

# Adjust to match your project
set(LIBNAME "example")
set(GOSRC "${CMAKE_CURRENT_SOURCE_DIR}/../../")
set(GOMAIN "example/bridge/bridge.gen.go")

# Only edit below this line for advanced use cases

set(GOLDF "-s -w")
if (CMAKE_BUILD_TYPE STREQUAL Debug)
    set(GOLDF "")
endif ()

set(ARCH_MAP_i686 386)
set(ARCH_MAP_x86_64 amd64)
set(ARCH_MAP_armv7a arm)
set(ARCH_MAP_aarch64 arm64)

set(ANDCC_MAP_i686 i686-linux-android)
set(ANDCC_MAP_x86_64 x86_64-linux-android)
set(ANDCC_MAP_armv7a armv7a-linux-androideabi)
set(ANDCC_MAP_aarch64 aarch64-linux-android)

string(REPLACE "-" "" SYSTEM_PROCESSOR ${CMAKE_SYSTEM_PROCESSOR})
set(GOARCH ${ARCH_MAP_${SYSTEM_PROCESSOR}})

if (CMAKE_SYSTEM_NAME STREQUAL Linux)
    set(GOOS linux)
    set(GOCC "${CMAKE_C_COMPILER}")
    set(GOOUT "${CMAKE_CURRENT_BINARY_DIR}/lib${LIBNAME}.so")
elseif (CMAKE_SYSTEM_NAME STREQUAL Windows)
    set(GOOS windows)
    # Leaving GOCC to default on Windows. MSVCC is not compatible with golang.
    set(GOOUT "${CMAKE_CURRENT_BINARY_DIR}/${LIBNAME}.dll")
elseif (CMAKE_SYSTEM_NAME STREQUAL Android)
    set(GOOS android)
    set(GOCC "${ANDROID_TOOLCHAIN_ROOT}/bin/${ANDCC_MAP_${SYSTEM_PROCESSOR}}${ANDROID_NATIVE_API_LEVEL}-clang")
    set(GOOUT "${CMAKE_LIBRARY_OUTPUT_DIRECTORY}/lib${LIBNAME}.so")
endif ()

add_custom_command(
        OUTPUT ${GOOUT}
        WORKING_DIRECTORY ${GOSRC}
        COMMENT "Building library"
        VERBATIM
        COMMAND ${CMAKE_COMMAND} -E env
        CGO_ENABLED=1
        GOOS=${GOOS} GOARCH=${GOARCH}
        CC=${GOCC}
        go build -buildmode=c-shared -ldflags=${GOLDF} -o ${GOOUT}
        ${GOMAIN}
)

add_custom_target(lib${LIBNAME} DEPENDS ${GOOUT})
