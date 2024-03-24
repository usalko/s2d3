/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: resources.go
 */
package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func V2Resource(bucket string, request *http.Request) []byte {
	resourceRequest := []byte(fmt.Sprintf("/%s%s", bucket, request.URL.Path))
	if request.URL.RawQuery == "" {
		return resourceRequest
	}

	requestParams := strings.Split(request.URL.RawQuery, "&")
	sort.Strings(requestParams)

	requestParts := make([][]byte, len(requestParams))
	for i := range requestParams {
		keyValue := strings.SplitN(requestParams[i], "=", 2)
		key := []byte(url.QueryEscape(keyValue[0]))
		if len(keyValue) == 2 {
			value := url.QueryEscape(keyValue[1])

			requestParts[i] = make([]byte, len(key)+1+len(value))
			copy(requestParts[i], key)
			requestParts[i][len(key)] = 0x3d
			copy(requestParts[i][1+len(key):], value)

		} else {
			requestParts[i] = key
		}
	}

	return bytes.Join([][]byte{
		resourceRequest,
		bytes.Join(requestParts, []byte{0x26}),
	}, []byte{0x3f})
}
