package main

import (
	"flag"
	"fmt"

	"github.com/usalko/s2d3"
)

func main() {
	fmt.Println("s2d3 utility for running simple s3 compatible service")
	ipAddr := flag.String("a", "0.0.0.0", "ip address ")
	ipPort := flag.Int("p", 3333, "ip port")
	localFolder := flag.String("d", "/tmp", "local folder")

	flag.Parse()

	s2d3.Serve(*localFolder, *ipAddr, *ipPort)
}
