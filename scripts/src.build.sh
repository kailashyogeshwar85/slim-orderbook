#!/usr/bin/env bash

set -e

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
BUILD_DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

export CGO_ENABLED=0

pushd $BUILD_DIR

BUILD_TIME="$(date -u '+%Y-%m-%d_%I:%M:%S%p')"
TAG="current"
REVISION="current"

LD_FLAGS="-s -w -X github.com/kailashyogeshwar85/slim-orderbook/pkg/version.appVersionTag=${TAG} -X github.com/kailashyogeshwar85/slim-orderbook/pkg/version.appVersionRev=${REVISION} -X github.com/kailashyogeshwar85/slim-orderbook/pkg/version.appVersionTime=${BUILD_TIME}"

pushd $BUILD_DIR/cmd/slim-orderbook

GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="${LD_FLAGS}" -a -tags 'netgo osusergo' -o "${BUILD_DIR}/bin/linux/slim-orderbook"
GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="${LD_FLAGS}" -a -tags 'netgo osusergo' -o "${BUILD_DIR}/bin/mac/slim-orderbook"
GOOS=linux GOARCH=arm go build -trimpath -ldflags="${LD_FLAGS}" -a -tags 'netgo osusergo' -o "$BUILD_DIR/bin/linux_arm/slim-orderbook"
GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="${LD_FLAGS}" -a -tags 'netgo osusergo' -o "$BUILD_DIR/bin/linux_arm64/slim-orderbook"
popd
popd

rm -rvf ${BUILD_DIR}/dist_linux
mkdir ${BUILD_DIR}/dist_linux
cp ${BUILD_DIR}/bin/linux/slim-orderbook ${BUILD_DIR}/dist_linux/slim-orderbook

pushd $BUILD_DIR
tar -czvf dist_linux.tar.gz dist_linux
popd

rm -rfv ${BUILD_DIR}/dist_linux
