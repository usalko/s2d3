/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */

package models

import (
	"fmt"
)

type SizeInBytes int64

func (b SizeInBytes) String() string {
	if b < 1<<10 {
		return b.Bytes()
	}
	if b < 1<<20 {
		return b.Kilobytes()
	}
	if b < 1<<30 {
		return b.Megabytes()
	}
	if b < 1<<40 {
		return b.Gigabytes()
	}
	if b < 1<<50 {
		return b.Terabytes()
	}
	if b < 1<<60 {
		return b.Petabytes()
	}
	return b.Exabytes()
}

func (b SizeInBytes) Bytes() string {
	return fmt.Sprintf("%db", b)
}

func (b SizeInBytes) Kilobytes() string {
	return fmt.Sprintf("%dk", b/(1<<10))
}

func (b SizeInBytes) Megabytes() string {
	return fmt.Sprintf("%dm", b/(1<<20))
}

func (b SizeInBytes) Gigabytes() string {
	return fmt.Sprintf("%dg", b/(1<<30))
}

func (b SizeInBytes) Terabytes() string {
	return fmt.Sprintf("%dt", b/(1<<40))
}

func (b SizeInBytes) Petabytes() string {
	return fmt.Sprintf("%dp", b/(1<<50))
}

func (b SizeInBytes) Exabytes() string {
	return fmt.Sprintf("%dx", b/(1<<60))
}
