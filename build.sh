#!/bin/bash

set -e  # Exit on error

mkdir -p build

build() {
  export GOOS=$1
  export GOARCH=$2
  export CC=$3
  export CXX=$4
  export CGO_ENABLED=1

  echo "Building for $GOOS-$GOARCH..."
  go build -o "build/space-shmup-$GOOS-$GOARCH"
}

# Windows (Use MinGW-w64 it must be installed on build system)
export MINGW_CC=x86_64-w64-mingw32-gcc
export MINGW_CXX=x86_64-w64-mingw32-g++

build windows amd64 "$MINGW_CC" "$MINGW_CXX"
mv build/space-shmup-windows-amd64 build/space-shmup-windows-amd64.exe

# Linux 64-bit
build linux amd64

# macOS builds are not working because of cgo

# macOS x86_64
#build darwin amd64

# macOS ARM64
#build darwin arm64

echo "Build completed! Binaries are in the 'build' directory."
