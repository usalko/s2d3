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
			return
		}
		Get(writer, request)
		return

	case "POST":
		_, exists := parsedQuery["uploads"]
		if exists {
			Upload(writer, request)
			return
		}

		_, exists = parsedQuery["uploadId"]
		if exists {
			Upload(writer, request)
			return
		}

	case "PUT":
		_, exists := parsedQuery["uploadId"]
		if exists {
			Upload(writer, request)
			return
		}

	}

	fmt.Printf("%s: [%s] %s request not processed\n", ctx.Value(KeyServerAddr), request.Method, request.URL.Path)
}

func GetHello(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(KeyServerAddr))
	io.WriteString(writer, "Hello, HTTP!\n")
}
