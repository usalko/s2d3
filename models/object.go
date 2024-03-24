/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: object.go
 */
package models

import (
	"time"

	"github.com/usalko/s2d3/utils"
)

type Object struct {
	Key          string
	LastModified time.Time
	ETag         string
	Size         utils.SizeInBytes
	StorageClass string
	OwnerID      string
	OwnerName    string
}
