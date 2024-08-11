/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */

package services

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type EntryOwner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

type Entry struct {
	Key          string     `xml:"Key"`
	LastModified string     `xml:"LastModified"`
	ETag         string     `xml:"ETag"`
	Size         int64      `xml:"Size"`
	StorageClass string     `xml:"StorageClass"`
	Owner        EntryOwner `xml:"Owner"`
}

type ListResponse struct {
	XMLName  xml.Name `xml:"ListBucketResult"`
	Next     string   `xml:"NextContinuationToken"`
	Contents []Entry  `xml:"Contents"`
}

func List(writer http.ResponseWriter, request *http.Request, listType any) {
	response := &ListResponse{}
	responseBytes, err := xml.Marshal(response)

	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	writer.Write(responseBytes)
}
