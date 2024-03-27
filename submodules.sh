#!/bin/bash

# @see https://stackoverflow.com/questions/67678203/why-does-go-get-fail-with-invalid-version-unknown-revision

cd utils || exit
GOPRIVATE=github.com/usalko/s2d3/utils go get github.com/usalko/s2d3/utils@v0.1.0-alpha.
cd ..

GOPROXY=proxy.golang.org go list -m github.com/usalko/s2d3/client
GOPROXY=proxy.golang.org go list -m github.com/usalko/s2d3/models
GOPROXY=proxy.golang.org go list -m github.com/usalko/s2d3/services
GOPROXY=proxy.golang.org go list -m github.com/usalko/s2d3/utils
echo ---------
GOPROXY=proxy.golang.org go list -m github.com/usalko/s2d3
