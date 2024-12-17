package utils

import "strconv"

func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}
