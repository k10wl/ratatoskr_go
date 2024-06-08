package main

import (
	"net/http"
	"ratatoskr/internal/webapp"
)

func main() {
	if err := webapp.Start(); err != http.ErrServerClosed {
		panic(err)
	}
}
