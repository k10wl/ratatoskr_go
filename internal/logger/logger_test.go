package logger

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder
	infoFlags = 0
	errFlags = 0
	reset := func() {
		stdout.Reset()
		stderr.Reset()
	}
	logger := NewLogger("TestLogger", &stdout, &stderr)

	type tc struct {
		name     string
		log      func(string) error
		value    string
		stdout   string
		stderr   string
		expected error
	}

	table := []tc{
		{
			name:  "log level only into stdout",
			value: "info level",
			log: func(s string) error {
				logger.Info(s)
				return nil
			},
			stdout:   "TestLogger: [INFO] -- {{ info level }}\n",
			stderr:   "",
			expected: nil,
		},

		{
			name:  "warning level only into stdout",
			value: "warning level",
			log: func(s string) error {
				logger.Warning(s)
				return nil
			},
			stdout:   "TestLogger: [WARNING] -- {{ warning level }}\n",
			stderr:   "",
			expected: nil,
		},

		{
			name:  "error level into stdout and stderr",
			value: "error level",
			log: func(s string) error {
				return logger.Error(s)
			},
			stdout:   "TestLogger: [ERROR] -- {{ error level }}\n",
			stderr:   "TestLogger: [ERROR] -- {{ error level }}\n",
			expected: fmt.Errorf("error level"),
		},
	}

	for _, test := range table {
		reset()
		err := test.log(test.value)
		if test.stdout != stdout.String() {
			t.Errorf(
				"Failed: %s - stdout\nexpected: %s\nactual:   %s",
				test.name,
				test.stdout,
				stdout.String(),
			)
		}
		if test.stderr != stderr.String() {
			t.Errorf(
				"Failed: %s - stderr\nexpected: %s\nactual:   %s",
				test.name,
				test.stderr,
				stderr.String(),
			)
		}
		if !reflect.DeepEqual(test.expected, err) {
			t.Errorf(
				"Failed error output creating %s\nexpected: %v\nactual:   %v",
				test.name,
				test.expected,
				err,
			)
		}
	}
}
