package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateUID() string {
	uuid := GenerateUUID()
	uuidParts := strings.Split(uuid, "-")
	return uuidParts[len(uuidParts)-1]
}
