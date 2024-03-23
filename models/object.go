/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */
package models

import (
	"time"
)

type Object struct {
	Key          string
	LastModified time.Time
	ETag         string
	Size         SizeInBytes
	StorageClass string
	OwnerID      string
	OwnerName    string
}
