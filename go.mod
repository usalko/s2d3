module github.com/usalko/s2d3

go 1.21.2

require (
	github.com/usalko/s2d3/client v0.1.5
	github.com/usalko/s2d3/services v0.1.5
)

require (
	github.com/usalko/s2d3/models v0.1.5 // indirect
	github.com/usalko/s2d3/utils v0.1.5 // indirect
	golang.org/x/net v0.22.0 // indirect
)

replace (
	github.com/usalko/s2d3/client v0.1.5 => ./client
	github.com/usalko/s2d3/models v0.1.5 => ./models
	github.com/usalko/s2d3/services v0.1.5 => ./services
	github.com/usalko/s2d3/utils v0.1.5 => ./utils
)
