package config

import (
	"fmt"
)

type Config struct {
	Token         string
	WebAppURL     string
	WebhookSecret string
}

func Init(getenv func(string) string) (*Config, error) {
	token := getenv("TOKEN")
	if token == "" {
		return nil, fmt.Errorf("bot token not provided")
	}
	webAppURL := getenv("URL")
	if webAppURL == "" {
		return nil, fmt.Errorf("URL environment variable is empty")
	}
	webhookSecret := getenv("WEBHOOK_SECRET")
	if webhookSecret == "" {
		return nil, fmt.Errorf("WEBHOOK_SECRET environment variable is empty")
	}
	return &Config{
		Token:         token,
		WebhookSecret: webhookSecret,
		WebAppURL:     webAppURL,
	}, nil
}
