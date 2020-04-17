#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-legacy-redirector
  make audit
popd
