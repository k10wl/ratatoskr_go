package config

import (
	"fmt"
)

type BotConfig struct {
	Token     string
	AdminIDs  []int64
	WebAppUrl string
}

func GetBotConfig(getenv func(string) string) (*BotConfig, error) {
	token := getenv("TOKEN")
	if token == "" {
		return nil, fmt.Errorf("bot token not provided")
	}
	adminIDs, err := stringToIntSlice(getenv("ADMIN_IDS"))
	if err != nil {
		return nil, fmt.Errorf("error upon receiving admin IDs")
	}
	if len(adminIDs) == 0 {
		return nil, fmt.Errorf("admin IDs must be set")
	}
	webAppUrl := getenv("WEBAPP_URL")
	return &BotConfig{
		Token:     token,
		AdminIDs:  adminIDs,
		WebAppUrl: webAppUrl,
	}, nil
}
