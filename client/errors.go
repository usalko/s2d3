package client

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func ResponseError(response *http.Response) error {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return ResponseErrorFrom(b)
}

func ResponseErrorFrom(body []byte) error {
	var payload struct {
		XMLName xml.Name `xml:"Error"`
		Code    string   `xml:"Code"`
		Message string   `xml:"Message"`
	}
	if err := xml.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unable to parse response xml: %s", err)
	}

	return fmt.Errorf("%s (%s) [raw %s]", payload.Message, payload.Code, string(body))
}
