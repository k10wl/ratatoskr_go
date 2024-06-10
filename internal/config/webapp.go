package config

import (
	"fmt"
	"ratatoskr/internal/utils"
)

type WepAppConfig struct {
	AdminIDs []int64
	IP       string
	PORT     string
}

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
	return &WepAppConfig{AdminIDs: adminIDs, IP: ip, PORT: port}, nil
}
