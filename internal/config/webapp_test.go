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
				case "IP":
					return "127.0.0.1"
				case "PORT":
					return "8080"
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return "database name"
				case "TOKEN":
					return "TOKEN"
				default:
					return ""
				}
			},
			shouldError: true,
			expected:    nil,
		},

		{
			name: "should error if IP was not provided",
			getenv: func(s string) string {
				switch s {
				case "ADMIN_IDS":
					return "1234,7890"
				case "IP":
					return ""
				case "PORT":
					return "8080"
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return "database name"
				case "TOKEN":
					return "TOKEN"
				default:
					return ""
				}
			},
			shouldError: true,
			expected:    &WepAppConfig{AdminIDs: []int64{1234, 7890}},
		},

		{
			name: "should error if MONGO_URI was not provided",
			getenv: func(s string) string {
				switch s {
				case "ADMIN_IDS":
					return "1234,7890"
				case "IP":
					return "127.0.0.1"
				case "PORT":
					return "8080"
				case "MONGO_URI":
					return ""
				case "MONGO_DB_NAME":
					return "database name"
				case "TOKEN":
					return "TOKEN"
				default:
					return ""
				}
			},
			shouldError: true,
			expected:    &WepAppConfig{AdminIDs: []int64{1234, 7890}},
		},

		{
			name: "should error if MONGO_DB_NAME was not provided",
			getenv: func(s string) string {
				switch s {
				case "ADMIN_IDS":
					return "1234,7890"
				case "IP":
					return "127.0.0.1"
				case "PORT":
					return "8080"
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return ""
				case "TOKEN":
					return "TOKEN"
				default:
					return ""
				}
			},
			shouldError: true,
			expected:    &WepAppConfig{AdminIDs: []int64{1234, 7890}},
		},

		{
			name: "should error if PORT was not provided",
			getenv: func(s string) string {
				switch s {
				case "ADMIN_IDS":
					return "1234,7890"
				case "IP":
					return "127.0.0.1"
				case "PORT":
					return ""
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return "database name"
				case "TOKEN":
					return "TOKEN"
				default:
					return ""
				}
			},
			shouldError: true,
			expected:    &WepAppConfig{AdminIDs: []int64{1234, 7890}},
		},

		{
			name: "should error if TOKEN was not provided",
			getenv: func(s string) string {
				switch s {
				case "ADMIN_IDS":
					return "1234,7890"
				case "IP":
					return "127.0.0.1"
				case "PORT":
					return "8080"
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return "database name"
				case "TOKEN":
					return ""
				default:
					return ""
				}
			},
			shouldError: true,
			expected:    &WepAppConfig{AdminIDs: []int64{1234, 7890}},
		},

		{
			name: "should return expected config",
			getenv: func(s string) string {
				switch s {
				case "ADMIN_IDS":
					return "1234,7890"
				case "IP":
					return "127.0.0.1"
				case "PORT":
					return "8080"
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return "database name"
				case "TOKEN":
					return "TOKEN"
				default:
					return ""
				}
			},
			shouldError: false,
			expected: &WepAppConfig{
				Version:     WebAppVersion,
				AdminIDs:    []int64{1234, 7890},
				IP:          "127.0.0.1",
				Port:        "8080",
				MongoURI:    "mongo://<name>:<pass>",
				MongoDBName: "database name",
				Token:       "TOKEN",
			},
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
