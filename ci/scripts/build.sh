#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-legacy-redirector
  make build && mv build/dp-legacy-redirector $cwd/build
  cp Dockerfile.concourse $cwd/build
popd
