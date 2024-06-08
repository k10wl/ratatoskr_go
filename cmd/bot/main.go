package main

import (
	"fmt"
	"io"
	"os"
	"ratatoskr/internal/bot"
	"ratatoskr/internal/config"
	"ratatoskr/internal/logger"
)

func run(
	getenv func(string) string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	c, err := config.GetBotConfig(getenv)
	if err != nil {
		return err
	}

	err = bot.Run(logger.NewLogger("Telegram bot", stdout, stderr), c)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(os.Getenv, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
