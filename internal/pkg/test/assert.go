package test

import (
	"fmt"
	"testing"
)

func AssertEqual(t *testing.T, value interface{}, expected interface{}, message string) {
	if value == expected {
		return
	}

	t.Fatal(fmt.Sprintf("%v != %v: %s", value, expected, message))
}

func AssertNil(t *testing.T, value interface{}, message string) {
	if value == nil {
		return
	}

	t.Fatal(fmt.Sprintf("%v != nil: %s", value, message))
}
