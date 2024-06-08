package main

import (
	"fmt"
	"io"
	"os"
	"ratatoskr/internal/bot"
	"ratatoskr/internal/config"
	"ratatoskr/internal/logger"
	"time"
)

func run(
	getenv func(string) string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	c, err := config.Init(getenv)
	if err != nil {
		return err
	}

	err = bot.Run(logger.NewLogger("Telegram bot", stdout, stderr), c)
	if err != nil {
		return err
	}

	// TODO block exit forever
	time.Sleep(time.Second * 99999999)
	return nil
}

func main() {
	if err := run(os.Getenv, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
