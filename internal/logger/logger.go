package logger

import (
	"fmt"
	"io"
	"log"
)

type Logger struct {
	infoLogger *log.Logger
	errLogger  *log.Logger
}

const delimeterStart = "{{ "
const delimeterEnd = " }}"

var infoFlags = (log.Ldate | log.Ltime | log.Lmicroseconds)
var errFlags = (log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

func NewLogger(
	name string,
	stdout io.Writer,
	stderr io.Writer,
) *Logger {
	combined := io.MultiWriter(stdout, stderr)
	return &Logger{
		infoLogger: log.New(
			stdout,
			fmt.Sprintf("%s: ", name),
			infoFlags,
		),
		errLogger: log.New(
			combined,
			fmt.Sprintf("%s: ", name),
			errFlags,
		),
	}
}

func (l Logger) Info(s string) {
	l.infoLogger.Printf(format("INFO", s))
}

func (l Logger) Warning(s string) {
	l.infoLogger.Printf(format("WARNING", s))
}

func (l Logger) Error(s string) error {
	l.errLogger.Printf(format("ERROR", s))
	return fmt.Errorf(s)
}

func format(level string, message string) string {
	return fmt.Sprintf("[%s] -- %s%s%s\n", level, delimeterStart, message, delimeterEnd)
}
