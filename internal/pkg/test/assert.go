package test

import (
	"fmt"
	"testing"
)

// AssertEqual fails the test if value and expected are not equal, otherwise does nothing
func AssertEqual(t *testing.T, value interface{}, expected interface{}, message string) {
	if value == expected {
		return
	}

	t.Fatal(fmt.Sprintf("%v != %v: %s", value, expected, message))
}

// AssertNil fails the test if value is not nil, otherwise does nothing
func AssertNil(t *testing.T, value interface{}, message string) {
	if value == nil {
		return
	}

	t.Fatal(fmt.Sprintf("%v != nil: %s", value, message))
}
