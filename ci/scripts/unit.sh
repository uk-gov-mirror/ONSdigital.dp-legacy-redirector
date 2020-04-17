#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-legacy-redirector
  make test
popd
