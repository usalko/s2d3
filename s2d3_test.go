/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */

package s2d3

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/usalko/s2d3/client"
	"github.com/usalko/s2d3/services"
)

var serverAddr = ""

const TEST_SERVED_LOCAL_FOLDER = "./.s3data"
const TEST_OBJECT_CONTENT = "Test"
const TEST_OBJECT_PATH = "test123/test456"

func WithContextDecorator(handler http.HandlerFunc, dataFolder string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		ctx = context.WithValue(ctx, services.KeyServerAddr, serverAddr)
		ctx = context.WithValue(ctx, services.KeyDataFolder, dataFolder)
		handler(writer, request.WithContext(ctx))
	}
}

func TestList(t *testing.T) {
	InitStorage(TEST_SERVED_LOCAL_FOLDER)
	server := httptest.NewServer(WithContextDecorator(services.GetRoot, TEST_SERVED_LOCAL_FOLDER))
	// Close the server when test finishes
	defer server.Close()
	parsedUrl, _ := url.Parse(server.URL)
	serverAddr = parsedUrl.Host

	s3Client, err := client.NewClient(&client.Client{
		AccessKeyId: "",
		Domain:      parsedUrl.Host, //"localhost:3333",
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
	InitStorage(TEST_SERVED_LOCAL_FOLDER)
	server := httptest.NewServer(WithContextDecorator(services.GetRoot, TEST_SERVED_LOCAL_FOLDER))
	// Close the server when test finishes
	defer server.Close()
	parsedUrl, _ := url.Parse(server.URL)
	serverAddr = parsedUrl.Host

	s3Client, err := client.NewClient(&client.Client{
		AccessKeyId: "",
		Domain:      parsedUrl.Host, //"localhost:3333",
		Protocol:    "http",
	})
	if err != nil {
		t.Errorf("Error in attempt to create new client %d", err)
	}

	upload, err := s3Client.NewUpload("test123/test456", nil)
	if err != nil {
		t.Errorf("Error in attempt to upload object %d", err)
	}

	wroteBytesCount, err := upload.Stream(bytes.NewReader([]byte(TEST_OBJECT_CONTENT)), 5*1024*1024*1024)
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

func TestGet(t *testing.T) {
	InitStorage(TEST_SERVED_LOCAL_FOLDER)
	server := httptest.NewServer(WithContextDecorator(services.GetRoot, TEST_SERVED_LOCAL_FOLDER))
	// Close the server when test finishes
	defer server.Close()
	parsedUrl, _ := url.Parse(server.URL)
	serverAddr = parsedUrl.Host

	s3Client, err := client.NewClient(&client.Client{
		AccessKeyId: "",
		Domain:      parsedUrl.Host, //"localhost:3333",
		Protocol:    "http",
	})
	if err != nil {
		t.Errorf("Error in attempt to create new client %d", err)
	}

	reader, err := s3Client.Get(TEST_OBJECT_PATH)
	if err != nil {
		t.Errorf("Error in attempt to get object %d", err)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		t.Errorf("Error in attempt to write stream %d", err)
	}

	if len(body) == 0 {
		t.Errorf("Didn't read anything %d", err)
	}

	if string(body) != TEST_OBJECT_CONTENT {
		t.Errorf("Wrong content in object %d", err)
	}

}
