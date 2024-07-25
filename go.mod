module github.com/usalko/s2d3

go 1.21.2

require (
	github.com/usalko/s2d3/client v0.1.3
	github.com/usalko/s2d3/services v0.1.3
)

require (
	github.com/usalko/s2d3/models v0.1.3 // indirect
	github.com/usalko/s2d3/utils v0.1.3 // indirect
	golang.org/x/net v0.22.0 // indirect
)

replace (
	github.com/usalko/s2d3/client v0.1.3 => ./client
	github.com/usalko/s2d3/models v0.1.3 => ./models
	github.com/usalko/s2d3/services v0.1.3 => ./services
	github.com/usalko/s2d3/utils v0.1.3 => ./utils
)
