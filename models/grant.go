/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: grant.go
 */
package models

type Grant struct {
	GranteeID   string
	GranteeName string
	Group       string
	Permission  string
}
