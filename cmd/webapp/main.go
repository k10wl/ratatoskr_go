package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"ratatoskr/internal/config"
	"ratatoskr/internal/db"
	"ratatoskr/internal/logger"
	"ratatoskr/internal/mongo_db"
	"ratatoskr/internal/webapp"
	"sync"
	"time"
)

func run(
	ctx context.Context,
	getDB func(ctx context.Context, URI string, database string) (db.DB, error),
	getenv func(string) string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	l := logger.NewLogger("WebApp", stdout, stderr)
	c, err := config.GetWebAppConfig(getenv)
	if err != nil {
		return l.Error(err.Error())
	}

	db, err := getDB(ctx, c.MongoURI, c.MongoDBName)
	if err != nil {
		return l.Error(err.Error())
	}

	svr, err := webapp.NewServer(db, l)
	if err != nil {
		return l.Error(err.Error())
	}

	httpServer := &http.Server{
		Handler: svr,
		Addr:    fmt.Sprintf("[%s]:%s", c.IP, c.Port),
	}

	go func() {
		l.Info(fmt.Sprintf("Starting server on %v", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			l.Error(fmt.Sprintf("error shutting down http server: %s\n", err))
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	if err := run(
		context.Background(),
		func(ctx context.Context, URI string, database string) (db.DB, error) {
			return mongo_db.NewMongoDB(ctx, URI, database)
		},
		os.Getenv,
		os.Stdout,
		os.Stderr,
	); err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
