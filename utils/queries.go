package utils

import (
	"bytes"
	"net/url"
	"sort"
	"strings"
)

func V4QueryString(queryString string) []byte {
	//TODO: using standard library
	if queryString == "" {
		return []byte{}
	}

	queryParams := strings.Split(queryString, "&")
	sort.Strings(queryParams)

	queryParamsLengths := make([][]byte, len(queryParams))
	for i := range queryParams {
		keyValue := strings.SplitN(queryParams[i], "=", 2)
		key := url.QueryEscape(keyValue[0])
		var value []byte
		if len(keyValue) == 2 {
			value = []byte(url.QueryEscape(keyValue[1]))
		}

		queryParamsLengths[i] = make([]byte, len(key)+1+len(value))
		copy(queryParamsLengths[i], key)
		queryParamsLengths[i][len(key)] = 0x3d
		copy(queryParamsLengths[i][1+len(key):], value)
	}

	return bytes.Join(queryParamsLengths, []byte{0x26})
}
