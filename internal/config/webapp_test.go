package config

import (
	"reflect"
	"testing"
)

func TestGetWebAppConfig(t *testing.T) {
	type tc struct {
		name        string
		getenv      func(string) string
		shouldError bool
		expected    *WepAppConfig
	}

	table := []tc{
		{
			name: "should return error if ADMIN_IDS were not provided",
			getenv: func(s string) string {
				switch s {
				case "ADMIN_IDS":
					return ""
				default:
					return ""
				}
			},
			shouldError: true,
			expected:    nil,
		},

		{
			name: "should return expected config",
			getenv: func(s string) string {
				switch s {
				case "ADMIN_IDS":
					return "1234,7890"
				default:
					return ""
				}
			},
			shouldError: false,
			expected:    &WepAppConfig{AdminIDs: []int64{1234, 7890}},
		},
	}

	for _, test := range table {
		config, err := GetWebAppConfig(test.getenv)
		if test.shouldError {
			if err == nil {
				t.Errorf("%s did not error", test.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s unexpected error: %+v", test.name, err)
		}
		if !reflect.DeepEqual(test.expected, config) {
			t.Errorf(
				"%s returned unexpected value\nexpected: %+v\nactual:   %+v",
				test.name,
				test.expected,
				config,
			)
		}
	}
}
