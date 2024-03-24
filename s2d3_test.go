/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */

package s2d3

import (
	"io/fs"
	"os"
	"testing"

	"github.com/usalko/s2d3/client"
)

func TestList(t *testing.T) {
	os.Mkdir("s3data", fs.ModeAppend)
	_, cancelFunc := Serve("./s3data")
	defer cancelFunc()

	s3Client, err := client.NewClient(&client.Client{
		AccessKeyId: "",
		Domain:      "localhost:3333",
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

	// a := 1
	// b := 2
	// expected := a + b

	// if got := Add(a, b); got != expected {
	// 	t.Errorf("Add(%d, %d) = %d, didn't return %d", a, b, got, expected)
	// }
}
