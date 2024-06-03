package main

import (
	"fmt"
	"io"
	"os"
)

func run(stdout io.Writer, stderr io.Writer) error {
	return nil
}

func main() {
	if err := run(os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
