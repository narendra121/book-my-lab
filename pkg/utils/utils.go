package utils

import (
	"time"
)

func GetExpTime(exp int64) int64 {
	return time.Now().Add(time.Duration(exp) * time.Minute).Unix()
}
func IsExpired(exp int64) bool {
	return exp < time.Now().Unix()
}
