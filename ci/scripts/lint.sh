#!/bin/bash -eux

pushd dp-legacy-redirector
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
  make lint
popd