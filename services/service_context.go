/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */

package services

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ServiceContext struct {
}

type ServiceContextKey string

const KeyServerAddr ServiceContextKey = "serverAddr"
const KeyDataFolder ServiceContextKey = "dataFolder"

func GetRoot(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	fmt.Printf("%s: got / request\n", ctx.Value(KeyServerAddr))

	parsedQuery, err := url.ParseQuery(request.URL.RawQuery)

	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	switch request.Method {

	case "GET":
		listType, exists := parsedQuery["list-type"]
		if exists {
			List(writer, request, listType)
		} else {
			fmt.Printf("%s: got / request\n", ctx.Value(KeyServerAddr))
		}

	case "POST":
		_, exists := parsedQuery["uploads"]
		if exists {
			Upload(writer, request)
		} else {
			fmt.Printf("%s: got / request\n", ctx.Value(KeyServerAddr))
		}

	case "PUT":
		_, exists := parsedQuery["uploadId"]
		if exists {
			Upload(writer, request)
		} else {
			fmt.Printf("%s: got / request\n", ctx.Value(KeyServerAddr))
		}

	default:
		fmt.Printf("%s: got / request\n", ctx.Value(KeyServerAddr))
	}

}

func GetHello(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(KeyServerAddr))
	io.WriteString(writer, "Hello, HTTP!\n")
}
