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
  go build -o "build/myapp-$GOOS-$GOARCH"
}

# Windows (Use MinGW-w64 instead of Zig)
export MINGW_CC=x86_64-w64-mingw32-gcc
export MINGW_CXX=x86_64-w64-mingw32-g++

build windows amd64 "$MINGW_CC" "$MINGW_CXX"

# Linux 64-bit
build linux amd64

# macOS x86_64 (Use Zig)
#build darwin amd64

# macOS ARM64
build darwin arm64

echo "Build completed! Binaries are in the 'build' directory."
