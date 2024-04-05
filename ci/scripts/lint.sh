#!/bin/bash -eux

pushd dp-legacy-redirector
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2
  make lint
popd