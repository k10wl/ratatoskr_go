package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"ratatoskr/internal/webapp"
)

func run(
	getenv func(string) string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	if err := webapp.Start(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func main() {
	if err := run(
		os.Getenv,
		os.Stdout,
		os.Stderr,
	); err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
