package utils

import (
	"strconv"
	"strings"
)

func StringToIntSlice(s string) ([]int64, error) {
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
