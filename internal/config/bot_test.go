package config

import (
	"reflect"
	"testing"
)

func TestGetBotConfig(t *testing.T) {
	type tc struct {
		name        string
		shouldError bool
		getenv      func(string) string
		expected    *BotConfig
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
				case "WEBAPP_URL":
					return "https:// link is required"
				case "RECEIVER_ID":
					return "1234"
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return "database name"
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
				case "WEBAPP_URL":
					return "https:// link is required"
				case "RECEIVER_ID":
					return "1234"
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return "database name"
				default:
					return ""
				}
			},
			expected: nil,
		},

		{
			name:        "should fail if webapp url is not provided",
			shouldError: true,
			getenv: func(s string) string {
				switch s {
				case "TOKEN":
					return "TOKEN"
				case "ADMIN_IDS":
					return "1234,7890"
				case "WEBAPP_URL":
					return ""
				case "RECEIVER_ID":
					return "1234"
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return "database name"
				default:
					return ""
				}
			},
			expected: nil,
		},

		{
			name:        "should fail if MONGO_URI is not provided",
			shouldError: true,
			getenv: func(s string) string {
				switch s {
				case "TOKEN":
					return "TOKEN"
				case "ADMIN_IDS":
					return "1234,7890"
				case "WEBAPP_URL":
					return "https:// link is required"
				case "RECEIVER_ID":
					return "1234"
				case "MONGO_URI":
					return ""
				case "MONGO_DB_NAME":
					return "database name"
				default:
					return ""
				}
			},
			expected: nil,
		},

		{
			name:        "should fail if MONGO_DB_NAME is not provided",
			shouldError: true,
			getenv: func(s string) string {
				switch s {
				case "TOKEN":
					return "TOKEN"
				case "ADMIN_IDS":
					return "1234,7890"
				case "WEBAPP_URL":
					return "https:// link is required"
				case "RECEIVER_ID":
					return "1234"
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return ""
				default:
					return ""
				}
			},
			expected: nil,
		},

		{
			name:        "should fail if RECEIVER_ID is not provided",
			shouldError: true,
			getenv: func(s string) string {
				switch s {
				case "TOKEN":
					return "TOKEN"
				case "ADMIN_IDS":
					return "1234,7890"
				case "WEBAPP_URL":
					return "https:// link is required"
				case "RECEIVER_ID":
					return ""
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return "database name"
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
				case "WEBAPP_URL":
					return "https:// link is required"
				case "RECEIVER_ID":
					return "1234"
				case "MONGO_URI":
					return "mongo://<name>:<pass>"
				case "MONGO_DB_NAME":
					return "database name"
				default:
					return ""
				}
			},
			expected: &BotConfig{
				Version:     BotVersion,
				Token:       "TOKEN",
				AdminIDs:    []int64{1, 2},
				WebAppUrl:   "https:// link is required",
				ReceiverID:  1234,
				MongoURI:    "mongo://<name>:<pass>",
				MongoDBName: "database name",
			},
		},
	}

	for _, test := range table {
		c, err := GetBotConfig(test.getenv)
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
