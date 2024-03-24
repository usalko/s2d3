#!/bin/bash

GOPROXY=proxy.golang.org go list -m github.com/usalko/s2d3/client
GOPROXY=proxy.golang.org go list -m github.com/usalko/s2d3/models
GOPROXY=proxy.golang.org go list -m github.com/usalko/s2d3/services
GOPROXY=proxy.golang.org go list -m github.com/usalko/s2d3/utils
echo ---------
GOPROXY=proxy.golang.org go list -m github.com/usalko/s2d3
