package config

import (
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	Token    string
	AdminIDs []int64
}

func Init(getenv func(string) string) (*Config, error) {
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
	return &Config{
		Token:    token,
		AdminIDs: adminIDs,
	}, nil
}

func stringToIntSlice(s string) ([]int64, error) {
	splitted := strings.Split(s, ",")
	ids := []int64{}
	var err error
	for _, val := range splitted {
		id, e := strconv.ParseInt(val, 10, 64)
		if id == 0 {
			continue
		}
		if e != nil {
			err = e
			break
		}
		ids = append(ids, id)
	}
	return ids, err
}
