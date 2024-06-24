package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"ratatoskr/internal/bot"
	"ratatoskr/internal/config"
	"ratatoskr/internal/logger"
	"ratatoskr/internal/mongo_db"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := mongo_db.NewMongoDB(ctx, c.MongoURI, c.MongoDBName)
	if err != nil {
		return err
	}

	err = bot.Run(db, logger.NewLogger("Telegram bot", stdout, stderr), c)
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
