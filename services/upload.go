/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */

package services

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/usalko/s2d3/models"
)

type UploadStart struct {
	Bucket   string `xml:"Bucket"`
	Key      string `xml:"Key"`
	UploadId string `xml:"UploadId"`
}

type UploadPart struct {
	XMLName xml.Name         `xml:"CompleteMultipartUpload"`
	Parts   []models.XmlPart `xml:"Part"`
}

type UploadDone struct {
	XMLName xml.Name         `xml:"CompleteMultipartUpload"`
	Parts   []models.XmlPart `xml:"Part"`
}

func Upload(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	switch request.Method {

	case "POST":

		parsedQuery, err := url.ParseQuery(request.URL.RawQuery)
		if err != nil {
			fmt.Printf("%s", err)
			return
		}

		uploadIds, exists := parsedQuery["uploadId"]
		if exists {
			for _, uploadId := range uploadIds {
				path, err := base64.StdEncoding.DecodeString(uploadId)
				if err != nil {
					fmt.Printf("%s", err)
					return
				}

				body, err := io.ReadAll(request.Body)
				if err != nil {
					fmt.Printf("%s", err)
					return
				}

				payload := UploadDone{}
				err = xml.Unmarshal(body, &payload)
				if err != nil {
					fmt.Printf("%s", err)
					return
				}
				println(string(path[:]))
			}
			return
		} else {
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
			return

		}
	case "PUT":
		parsedQuery, err := url.ParseQuery(request.URL.RawQuery)
		if err != nil {
			fmt.Printf("%s", err)
			return
		}

		uploadIds, exists := parsedQuery["uploadId"]

		if exists {
			for _, uploadId := range uploadIds {
				path, err := base64.StdEncoding.DecodeString(uploadId)
				if err != nil {
					fmt.Printf("%s", err)
					return
				}

				body, err := io.ReadAll(request.Body)
				if err != nil {
					fmt.Printf("%s", err)
					return
				}

				// payload := UploadPart{}
				// err = xml.Unmarshal(body, &payload)
				// if err != nil {
				// 	fmt.Printf("%s", err)
				// 	return
				// }
				println(string(body[:]))
				println(string(path[:]))
			}
			return
		}

	}
	fmt.Printf("%s: [%s] /%s request\n", ctx.Value(KeyServerAddr), request.Method, request.URL.Path)

}
