package common

import "strconv"

func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return value
}
