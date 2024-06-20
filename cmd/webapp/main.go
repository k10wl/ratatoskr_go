package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"ratatoskr/internal/config"
	"ratatoskr/internal/logger"
	"ratatoskr/internal/webapp"
)

func run(
	getenv func(string) string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	l := logger.NewLogger("WebApp", stdout, stderr)
	c, err := config.GetWebAppConfig(getenv)
	if err != nil {
		return l.Error(err.Error())
	}

	svr, err := webapp.NewServer(l)
	if err != nil {
		return l.Error(err.Error())
	}

	httpServer := &http.Server{
		Handler: svr,
		Addr:    fmt.Sprintf("[%s]:%s", c.IP, c.PORT),
	}

	l.Info(fmt.Sprintf("Starting server on %v", httpServer.Addr))
	return httpServer.ListenAndServe()

}

func main() {
	if err := run(
		os.Getenv,
		os.Stdout,
		os.Stderr,
	); err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
