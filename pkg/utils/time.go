package utils

import (
	"fmt"
	"time"
)

func GetUnixTimeStamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func GetNow() time.Time {
	return time.Now()
}

func GetNowInString() string {
	return GetNow().Format("2006-01-02 15:04:05")
}
