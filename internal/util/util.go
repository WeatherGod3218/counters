package util

import (
	"time"
	_ "time/tzdata"
)

func ValidateTime(inputTime int64) int64 {
	now := time.Now().Unix()

	if timestamp := inputTime; timestamp < now {
		return timestamp
	}

	return now
}
