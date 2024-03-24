/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: headers.go
 */
package utils

import (
	"bytes"
	"net/http"
	"sort"
	"strings"
)

func V2Headers(request *http.Request) []byte {
	//TODO: using standard library
	subset := make(map[string]string)
	names := make([]string, 0)

	for header := range request.Header {
		LowerCaseHeader := strings.ToLower(header)
		if strings.HasPrefix(LowerCaseHeader, "x-amz-") {
			names = append(names, LowerCaseHeader)
			subset[LowerCaseHeader] = strings.Trim(request.Header.Get(header), " \t\r\n\f") + "\n"
		}
	}
	sort.Strings(names)

	headersLengths := make([][]byte, len(names))
	for i, header := range names {
		headersLengths[i] = bytes.Join([][]byte{[]byte(header), []byte(subset[header])}, []byte{0x3a})
	}
	return bytes.Join(headersLengths, []byte{})
}

func V4Headers(request *http.Request) ([]byte, []byte) {
	//TODO: using standard library
	subset := make(map[string]string)
	names := make([]string, 0)

	for header := range request.Header {
		lowerCaseHeader := strings.ToLower(header)
		if lowerCaseHeader == "host" || strings.HasPrefix(lowerCaseHeader, "x-amz-") {
			names = append(names, lowerCaseHeader)
			subset[lowerCaseHeader] = strings.Trim(request.Header.Get(header), " \t\r\n\f")
		}
	}
	sort.Strings(names)

	headersLengths := make([][]byte, len(names))
	headersValues := make([][]byte, len(names))
	for i, header := range names {
		headersValues[i] = []byte(header)
		headersLengths[i] = bytes.Join([][]byte{headersValues[i], []byte(subset[header])}, []byte{0x3a})
	}

	signed := bytes.Join(headersValues, []byte{0x3b})
	return signed, bytes.Join([][]byte{
		bytes.Join(headersLengths, []byte{0x0a}),
		nil, /* force an empty line */
		signed,
	}, []byte{0x0a})
}
