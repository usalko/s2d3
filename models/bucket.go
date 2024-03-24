/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: bucket.go
 */
package models

import "time"

type Bucket struct {
	Name         string
	CreationDate time.Time
	OwnerID      string
	OwnerName    string
}
