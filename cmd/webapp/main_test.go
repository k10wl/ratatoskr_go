package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"ratatoskr/internal/db"
	"ratatoskr/internal/models"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	go run(
		ctx,
		func(context.Context, string, string) (db.DB, error) { return dbMock{}, nil },
		getEnv,
		os.Stdout,
		os.Stderr,
	)
	err := waitForReady(ctx, time.Second, "http://127.0.0.1:8088/ping")
	if err != nil {
		t.Errorf("error upon waiting for server: %v", err)
	}

	res, err := http.Get("http://127.0.0.1:8088/")
	if err != nil {
		t.Errorf("error upon waiting for server: %v", err)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error upon waiting for server: %v", err)
		return
	}
	if res.StatusCode != http.StatusForbidden || len(data) > 0 {
		t.Error("did not block request without token")
	}

	res, err = http.Get("http://127.0.0.1:8088/ping")
	if err != nil {
		t.Errorf("error upon waiting for server: %v", err)
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error upon waiting for server: %v", err)
		return
	}
	if string(data) != "pong" {
		t.Error("failed to ping-pong")
	}

	res, err = http.Get("http://127.0.0.1:8088/TOKEN")
	if err != nil {
		t.Errorf("error upon waiting for server: %v", err)
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error upon waiting for server: %v", err)
		return
	}
	html := string(data)
	if !strings.Contains(html, "https://telegram.org/js/telegram-web-app.js") {
		t.Errorf("html does not contain telegram web app")
	}
	groups := []string{"group1", "group2"}
	if !containsMultipleSubstrings(html, groups) {
		t.Errorf("html does not contain some of groups: %v\n%v", groups, html)
	}
	tags := []string{"#tag1", "#tag2", "#tag3"}
	if !containsMultipleSubstrings(html, tags) {
		t.Errorf("html does not contain some of tags: %v\n%v", tags, html)
	}
}

func containsMultipleSubstrings(string string, substrings []string) bool {
	containst := true
	for _, substring := range substrings {
		if !strings.Contains(string, substring) {
			containst = false
			break
		}
	}
	return containst
}

func waitForReady(
	ctx context.Context,
	timeout time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %s\n", err.Error())
			continue
		}
		if resp.StatusCode == http.StatusOK {
			fmt.Println("Endpoint is ready!")
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func getEnv(s string) string {
	switch s {
	case "ADMIN_IDS":
		return "1234,7890"
	case "IP":
		return "127.0.0.1"
	case "PORT":
		return "8088"
	case "MONGO_URI":
		return "mongo://<name>:<pass>"
	case "MONGO_DB_NAME":
		return "database name"
	case "TOKEN":
		return "TOKEN"
	default:
		return ""
	}
}

type dbMock struct{}

func (_ dbMock) GetAllGroupsWithTags(context.Context) (*[]models.Group, error) {
	return &[]models.Group{
		{Name: "group1", Tags: []models.Tag{
			{Name: "#tag1"},
			{Name: "#tag2"},
			{Name: "#tag3"},
		}},
		{Name: "group2", Tags: []models.Tag{
			{Name: "#tag4"},
			{Name: "#tag5"},
			{Name: "#tag6"},
		}},
	}, nil
}
