module github.com/usalko/s2d3/services

go 1.21.2

require github.com/usalko/s2d3/models v0.1.5

require github.com/usalko/s2d3/utils v0.1.5

replace (
	github.com/usalko/s2d3/models v0.1.5 => ../models
	github.com/usalko/s2d3/utils v0.1.5 => ../utils
)
