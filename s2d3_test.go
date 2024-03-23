/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3_test.go
 */

package s2d3

import (
	"testing"
)

func TestList(t *testing.T) {
	a := 1
	b := 2
	expected := a + b

	// if got := Add(a, b); got != expected {
	// 	t.Errorf("Add(%d, %d) = %d, didn't return %d", a, b, got, expected)
	// }
	print(expected)
}
