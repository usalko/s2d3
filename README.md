# github.com/usalko/s2d3

Usage
Initialize your module
go mod init example.com/my-s2d3-demo

Get the s2d3 module
Note that you need to include the v in the version tag.

go get github.com/usalko/s2d3@v0.1.8

package main

import (
    "fmt"

    "github.com/usalko/s2d3"
)

func main() {
    s2d3.Serve("./s3data")
}
