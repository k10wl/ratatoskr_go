package config

import "fmt"

type WepAppConfig struct {
	AdminIDs []int64
}

func GetWebAppConfig(getenv func(string) string) (*WepAppConfig, error) {
	stringAdminIDs := getenv("ADMIN_IDS")
	if stringAdminIDs == "" {
		return nil, fmt.Errorf("required ADMIN_IDS was not provided")
	}
	adminIDs, err := stringToIntSlice(stringAdminIDs)
	if err != nil {
		return nil, err
	}
	return &WepAppConfig{AdminIDs: adminIDs}, nil
}
