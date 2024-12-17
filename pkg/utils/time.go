package utils

import (
	"fmt"
	"time"
)

func GetUnixTimeStamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
