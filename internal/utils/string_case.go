package utils

import (
	"strings"
)

func CamelCaseToTitleCase(s string) string {
	if len(s) == 0 {
		return ""
	}

	var result string
	for i, c := range s {
		if i == 0 {
			result += string(c)
			continue
		}
		if c >= 'A' && c <= 'Z' {
			result += " " + string(c)
		} else {
			result += string(c)
		}
	}
	return strings.ToTitle(result)
}
