package config

import (
	"testing"
)

func TestInit(t *testing.T) {
	type tc struct {
		name        string
		shouldError bool
		getenv      func(string) string
	}

	table := []tc{
		{
			name:        "should fail if TOKEN was not provided",
			shouldError: true,
			getenv: func(s string) string {
				switch s {
				case "TOKEN":
					return ""
				case "URL":
					return "URL"
				case "WEBHOOK_SECRET":
					return "WEBHOOK_SECRET"
				default:
					return ""
				}
			},
		},
		{
			name:        "should fail if URL was not provided",
			shouldError: true,
			getenv: func(s string) string {
				switch s {
				case "TOKEN":
					return "TOKEN"
				case "URL":
					return ""
				case "WEBHOOK_SECRET":
					return "WEBHOOK_SECRET"
				default:
					return ""
				}
			},
		},
		{
			name:        "should fail if WEBHOOK_SECRET was not provided",
			shouldError: true,
			getenv: func(s string) string {
				switch s {
				case "TOKEN":
					return "TOKEN"
				case "URL":
					return "URL"
				case "WEBHOOK_SECRET":
					return ""
				default:
					return ""
				}
			},
		},
		{
			name:        "should get config",
			shouldError: false,
			getenv: func(s string) string {
				switch s {
				case "TOKEN":
					return "TOKEN"
				case "URL":
					return "URL"
				case "WEBHOOK_SECRET":
					return "WEBHOOK_SECRET"
				default:
					return ""
				}
			},
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
		if c == nil {
			t.Errorf("Did not create config instance: %s", test.name)
		}
	}
}
