/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */

package services

import (
	"fmt"
	"io"
	"net/http"

	"github.com/usalko/s2d3/models"
)

type ServiceContext struct {
}

type ServiceContextKey string

const KeyServerAddr ServiceContextKey = "serverAddr"

func GetRoot(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	bucket := models.Bucket{
		Name: "",
	}
	println(&bucket)

	fmt.Printf("%s: got / request\n", ctx.Value(KeyServerAddr))
	io.WriteString(writer, "This is my website!\n")
}

func GetHello(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(KeyServerAddr))
	io.WriteString(writer, "Hello, HTTP!\n")
}
