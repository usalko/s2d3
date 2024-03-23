# github.com/usalko/s2d3

Usage
Initialize your module
go mod init example.com/my-s2d3-demo

Get the go-lib module
Note that you need to include the v in the version tag.

go get github.com/usalko/s2d3@v0.1.0

package main

import (
    "fmt"

    "github.com/usalko/s2d3"
)

func main() {
    fmt.Println(golib.Add(2,3))
    fmt.Println(golib.Subtract(2,3))
}

Testing
Tagging
git tag v0.1.0
git push origin --tags
