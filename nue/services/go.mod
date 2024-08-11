module github.com/usalko/s2d3/services

go 1.21.2

require github.com/usalko/s2d3/models v0.1.7

require github.com/usalko/s2d3/utils v0.1.7

replace (
	github.com/usalko/s2d3/models v0.1.7 => ../models
	github.com/usalko/s2d3/utils v0.1.7 => ../utils
)
