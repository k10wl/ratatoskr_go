package logger

import (
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
		name   string
		log    func(string)
		value  string
		stdout string
		stderr string
	}

	table := []tc{
		{
			name:  "log level only into stdout",
			value: "info level",
			log: func(s string) {
				logger.Info(s)
			},
			stdout: "TestLogger: [INFO] -- {{ info level }}\n",
			stderr: "",
		},

		{
			name:  "warning level only into stdout",
			value: "warning level",
			log: func(s string) {
				logger.Warning(s)
			},
			stdout: "TestLogger: [WARNING] -- {{ warning level }}\n",
			stderr: "",
		},

		{
			name:  "error level into stdout and stderr",
			value: "error level",
			log: func(s string) {
				logger.Error(s)
			},
			stdout: "TestLogger: [ERROR] -- {{ error level }}\n",
			stderr: "TestLogger: [ERROR] -- {{ error level }}\n",
		},
	}

	for _, test := range table {
		reset()
		test.log(test.value)
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
	}
}
