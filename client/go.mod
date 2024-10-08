module github.com/usalko/s2d3/client

go 1.21.2

require (
	github.com/usalko/s2d3/models v0.1.8
	github.com/usalko/s2d3/utils v0.1.8
	golang.org/x/net v0.22.0
)

replace (
	github.com/usalko/s2d3/models v0.1.8 => ../models
	github.com/usalko/s2d3/utils v0.1.8 => ../utils
)
