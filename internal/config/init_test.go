package config

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	type tc struct {
		name        string
		shouldError bool
		getenv      func(string) string
		expected    *Config
	}

	table := []tc{
		{
			name:        "should fail if TOKEN was not provided",
			shouldError: true,
			getenv: func(s string) string {
				switch s {
				case "TOKEN":
					return ""
				case "ADMIN_IDS":
					return "1,2"
				default:
					return ""
				}
			},
			expected: nil,
		},

		{
			name:        "should fail if admin ID's are not provided",
			shouldError: true,
			getenv: func(s string) string {
				switch s {
				case "TOKEN":
					return "TOKEN"
				case "ADMIN_IDS":
					return ""
				default:
					return ""
				}
			},
			expected: nil,
		},

		{
			name:        "should get config",
			shouldError: false,
			getenv: func(s string) string {
				switch s {
				case "TOKEN":
					return "TOKEN"
				case "ADMIN_IDS":
					return "1,2"
				default:
					return ""
				}
			},
			expected: &Config{Token: "TOKEN", AdminIDs: []int64{1, 2}},
		},
	}

	for _, test := range table {
		c, err := Init(test.getenv)
		if test.shouldError {
			if err == nil {
				t.Errorf("Expected error, but did not fail: %s", test.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("Unexpected error in %s: \nError:%s\n", test.name, err)
			continue
		}
		if !reflect.DeepEqual(*c, *test.expected) {
			t.Errorf(
				"Did not create correct config instance in %s \nexpected: %+v\nactual:   %+v",
				test.name,
				*test.expected,
				*c,
			)
		}
	}
}

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
		actual, err := stringToIntSlice(test.input)
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
