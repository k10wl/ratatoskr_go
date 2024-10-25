package config

import (
	"fmt"
	"ratatoskr/internal/utils"
)

type WepAppConfig struct {
	Version     string
	AdminIDs    []int64
	IP          string
	Port        string
	MongoURI    string
	MongoDBName string
	Token       string
}

const WebAppVersion = "1.1.4"

func GetWebAppConfig(getenv func(string) string) (*WepAppConfig, error) {
	stringAdminIDs := getenv("ADMIN_IDS")
	if stringAdminIDs == "" {
		return nil, fmt.Errorf("required ADMIN_IDS was not provided")
	}
	adminIDs, err := utils.StringToIntSlice(stringAdminIDs)
	if err != nil {
		return nil, err
	}
	ip := getenv("IP")
	if ip == "" {
		return nil, fmt.Errorf("required IP was not provided")
	}
	port := getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("required PORT was not provided")
	}
	mongoURI := getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("required PORT was not provided")
	}
	mongoDBName := getenv("MONGO_DB_NAME")
	if mongoDBName == "" {
		return nil, fmt.Errorf("required MONGO_DB_NAME was not provided")
	}
	token := getenv("TOKEN")
	if token == "" {
		return nil, fmt.Errorf("bot token not provided")
	}
	return &WepAppConfig{
		Version:     WebAppVersion,
		AdminIDs:    adminIDs,
		IP:          ip,
		Port:        port,
		MongoURI:    mongoURI,
		MongoDBName: mongoDBName,
		Token:       token,
	}, nil
}
