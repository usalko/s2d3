/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: size_in_bytes.go
 */

package utils

import (
	"fmt"
)

type SizeInBytes int64

func (sizeInBytes SizeInBytes) String() string {
	if sizeInBytes < 1<<10 {
		return sizeInBytes.Bytes()
	}
	if sizeInBytes < 1<<20 {
		return sizeInBytes.Kilobytes()
	}
	if sizeInBytes < 1<<30 {
		return sizeInBytes.Megabytes()
	}
	if sizeInBytes < 1<<40 {
		return sizeInBytes.Gigabytes()
	}
	if sizeInBytes < 1<<50 {
		return sizeInBytes.Terabytes()
	}
	if sizeInBytes < 1<<60 {
		return sizeInBytes.Petabytes()
	}
	return sizeInBytes.Exabytes()
}

func (sizeInBytes SizeInBytes) Bytes() string {
	return fmt.Sprintf("%db", sizeInBytes)
}

func (sizeInBytes SizeInBytes) Kilobytes() string {
	return fmt.Sprintf("%dk", sizeInBytes/(1<<10))
}

func (sizeInBytes SizeInBytes) Megabytes() string {
	return fmt.Sprintf("%dm", sizeInBytes/(1<<20))
}

func (sizeInBytes SizeInBytes) Gigabytes() string {
	return fmt.Sprintf("%dg", sizeInBytes/(1<<30))
}

func (sizeInBytes SizeInBytes) Terabytes() string {
	return fmt.Sprintf("%dt", sizeInBytes/(1<<40))
}

func (sizeInBytes SizeInBytes) Petabytes() string {
	return fmt.Sprintf("%dp", sizeInBytes/(1<<50))
}

func (sizeInBytes SizeInBytes) Exabytes() string {
	return fmt.Sprintf("%dx", sizeInBytes/(1<<60))
}
