/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */

package s2d3

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/usalko/s2d3/client"
	"github.com/usalko/s2d3/services"
)

func WithContext(handler http.HandlerFunc, dataFolder string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		ctx = context.WithValue(ctx, services.KeyDataFolder, dataFolder)
		handler(writer, request.WithContext(ctx))
	}
}

func TestList(t *testing.T) {
	localFolder := "./s3data"
	InitStorage(localFolder)
	server := httptest.NewServer(WithContext(services.GetRoot, localFolder))
	// Close the server when test finishes
	defer server.Close()
	u, _ := url.Parse(server.URL)

	s3Client, err := client.NewClient(&client.Client{
		AccessKeyId: "",
		Domain:      u.Host, //"localhost:3333",
		Protocol:    "http",
	})
	if err != nil {
		t.Errorf("Error in attempt to create new client %d", err)
	}

	result, err := s3Client.List()
	if err != nil {
		t.Errorf("Error in attempt to list objects %d", err)
	}

	println(result)
}

func TestUpload(t *testing.T) {
	localFolder := "./s3data"
	InitStorage(localFolder)
	server := httptest.NewServer(WithContext(services.GetRoot, localFolder))
	// Close the server when test finishes
	defer server.Close()
	u, _ := url.Parse(server.URL)

	s3Client, err := client.NewClient(&client.Client{
		AccessKeyId: "",
		Domain:      u.Host, //"localhost:3333",
		Protocol:    "http",
	})
	if err != nil {
		t.Errorf("Error in attempt to create new client %d", err)
	}

	upload, err := s3Client.NewUpload("test123/test456", nil)
	if err != nil {
		t.Errorf("Error in attempt to upload object %d", err)
	}

	wroteBytesCount, err := upload.Stream(bytes.NewReader([]byte("Test")), 5*1024*1024*1024)
	if err != nil {
		t.Errorf("Error in attempt to write stream %d", err)
	}

	if wroteBytesCount <= 0 {
		t.Errorf("Didn't write anything %d", err)
	}

	err = upload.Done()
	if err != nil {
		t.Errorf("Error in attempt to finish upload %d", err)
	}

}
