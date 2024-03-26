package client

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/usalko/s2d3/models"
)

func ResponseError(response *http.Response) error {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return ResponseErrorFrom(b)
}

func ResponseErrorFrom(body []byte) error {
	payload := models.Error{}
	if err := xml.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unable to parse response xml: %s", err)
	}

	return fmt.Errorf("%s (%s) [raw %s]", payload.Message, payload.Code, string(body))
}
