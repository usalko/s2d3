package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/usalko/s2d3"
)

const LOGO_ASCII_GRAPHIC = "\n" +
	"   .---.  2/3\n" +
	"  ( .-._)    \n" +
	" (_) \\       \n" +
	" _  \\ \\      \n" +
	"( `-'  )     \n" +
	" `----'      \n" +
	"             \n"

func main() {
	fmt.Println("s2d3 utility for running simple s3 compatible service")
	ipAddr := flag.String("a", "127.0.0.1", "ip address ")
	ipPort := flag.Int("p", 3333, "ip port")
	localFolder := flag.String("d", "/tmp", "local folder")
	urlContext := flag.String("u", "/", "url context")

	flag.Parse()

	http.Handle(*urlContext, &s2d3.ServeLocalFolder{
		RootFolder: *localFolder,
		UrlContext: *urlContext,
		ServerAddr: fmt.Sprintf("%s:%d", *ipAddr, *ipPort),
	})
	fmt.Print(LOGO_ASCII_GRAPHIC)
	fmt.Printf("Serve local folder '%s' \n", *localFolder)
	fmt.Printf("Please check url: http://%s:%d%s\n", *ipAddr, *ipPort, *urlContext)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *ipAddr, *ipPort), nil); err != nil {
		log.Fatal("S3 server terminated", err)
	}
}
