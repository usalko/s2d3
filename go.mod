module github.com/usalko/s2d3

go 1.21.2

require (
	github.com/usalko/s2d3/client v0.1.8
	github.com/usalko/s2d3/services v0.1.8
)

require (
	github.com/usalko/s2d3/models v0.1.8 // indirect
	github.com/usalko/s2d3/utils v0.1.8 // indirect
	golang.org/x/net v0.22.0 // indirect
)

replace (
	github.com/usalko/s2d3/client v0.1.8 => ./client
	github.com/usalko/s2d3/models v0.1.8 => ./models
	github.com/usalko/s2d3/services v0.1.8 => ./services
	github.com/usalko/s2d3/utils v0.1.8 => ./utils
)
