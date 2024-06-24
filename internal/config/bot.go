package config

import (
	"fmt"
	"ratatoskr/internal/utils"
	"strconv"
)

type BotConfig struct {
	Token       string
	AdminIDs    []int64
	WebAppUrl   string
	ReceiverID  int64
	MongoURI    string
	MongoDBName string
}

func GetBotConfig(getenv func(string) string) (*BotConfig, error) {
	token := getenv("TOKEN")
	if token == "" {
		return nil, fmt.Errorf("bot token not provided")
	}
	adminIDs, err := utils.StringToIntSlice(getenv("ADMIN_IDS"))
	if err != nil {
		return nil, fmt.Errorf("error upon receiving admin IDs")
	}
	if len(adminIDs) == 0 {
		return nil, fmt.Errorf("admin IDs must be set")
	}
	webAppUrl := getenv("WEBAPP_URL")
	if webAppUrl == "" {
		return nil, fmt.Errorf("required WEBAPP_URL was not provided")
	}
	receiverString := getenv("RECEIVER_ID")
	if receiverString == "" {
		return nil, fmt.Errorf("required RECEIVER_ID was not provided")
	}
	receiverID, err := strconv.Atoi(receiverString)
	if err != nil {
		return nil, fmt.Errorf("required RECEIVER_ID could not be parsed to int")
	}
	mongoURI := getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("required PORT was not provided")
	}
	mongoDBName := getenv("MONGO_DB_NAME")
	if mongoDBName == "" {
		return nil, fmt.Errorf("required MONGO_DB_NAME was not provided")
	}
	return &BotConfig{
		Token:       token,
		AdminIDs:    adminIDs,
		WebAppUrl:   webAppUrl,
		ReceiverID:  int64(receiverID),
		MongoURI:    mongoURI,
		MongoDBName: mongoDBName,
	}, nil
}
