package utils

import (
	"reflect"
	"testing"
)

func TestStringToIntSlice(t *testing.T) {
	type tc struct {
		name        string
		input       string
		expected    []int64
		shouldError bool
	}

	table := []tc{
		{
			name:     "should return error if wrong values were provided",
			input:    "this is not int",
			expected: []int64{},
		},

		{
			name:     "should return error if wrong values were provided",
			input:    ",",
			expected: []int64{},
		},

		{
			name:     "should return empty array",
			input:    "",
			expected: []int64{},
		},

		{
			name:     "should return slice of ints",
			input:    "123,456",
			expected: []int64{123, 456},
		},
	}

	for _, test := range table {
		actual, err := StringToIntSlice(test.input)
		if err != nil {
			t.Errorf("unexpected error in %s: %+v", test.name, err)
		}

		if !reflect.DeepEqual(test.expected, actual) {
			t.Errorf(
				"failed to give expected result %s\nexpected: %v\nactual:   %v",
				test.name,
				test.expected,
				actual,
			)
		}
	}
}
