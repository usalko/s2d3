/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */

package services

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

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

func bucketNameAndObjectKey(path string) (string, string) {
	bucketName, objectKey, exists := strings.Cut(path, "/")
	if exists {
		return bucketName, objectKey
	}
	return "", path
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
			path, err := url.QueryUnescape(request.URL.Path)
			if err != nil {
				fmt.Printf("%s", err)
				return
			}
			bucketName, objectKey := bucketNameAndObjectKey(path)

			uploadId := strings.Join([]string{path, hex.EncodeToString(new(big.Int).SetInt64(time.Now().UnixMicro()).Bytes())}, ":")

			response := &UploadStart{
				Bucket:   bucketName,
				Key:      objectKey,
				UploadId: base64.StdEncoding.EncodeToString([]byte(uploadId)),
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
			for _, encodedUploadId := range uploadIds {
				uploadId, err := base64.StdEncoding.DecodeString(encodedUploadId)
				if err != nil {
					fmt.Printf("%s", err)
					return
				}

				path, suffix, found := strings.Cut(string(uploadId), ":")
				if !found {
					fmt.Printf("Invalid uploadId %s", err)
					return
				}
				bucketName, objectName := bucketNameAndObjectKey(path)

				storage := Storage{
					RootFolder: request.Context().Value(KeyDataFolder).(string),
				}
				err = storage.PushData(bucketName, objectName, suffix, request.Body)
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
			}
			return
		}

	}
	fmt.Printf("%s: [%s] /%s request\n", ctx.Value(KeyServerAddr), request.Method, request.URL.Path)

}
