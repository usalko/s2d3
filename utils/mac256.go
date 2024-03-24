/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: mac256.go
 */
package utils

import (
	"crypto/hmac"
	"crypto/sha256"
)

func Mac256(key, message []byte) []byte {
	hmacHash := hmac.New(sha256.New, key)
	hmacHash.Write(message)
	return hmacHash.Sum(nil)
}
