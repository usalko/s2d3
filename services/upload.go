/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */

package services

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

type UploadStart struct {
	Bucket   string `xml:"Bucket"`
	Key      string `xml:"Key"`
	UploadId string `xml:"UploadId"`
}

func Upload(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	switch request.Method {

	case "POST":

		response := &UploadStart{
			Bucket:   "",
			Key:      request.URL.Path,
			UploadId: base64.StdEncoding.EncodeToString([]byte(request.URL.Path)), // TODO: add timestamp
		}
		responseBytes, err := xml.Marshal(response)

		if err != nil {
			fmt.Printf("%s", err)
			return
		}

		writer.Write(responseBytes)

	case "PUT":
		parsedQuery, err := url.ParseQuery(request.URL.RawQuery)
		if err != nil {
			fmt.Printf("%s", err)
			return
		}

		uploadId, exists := parsedQuery["uploadId"]

		if exists {
			println(uploadId)
		} else {
			fmt.Printf("%s: got / request\n", ctx.Value(KeyServerAddr))
		}

	default:
		fmt.Printf("%s: got / request\n", ctx.Value(KeyServerAddr))
	}

}
