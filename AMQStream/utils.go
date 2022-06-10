package AMQStream

import (
	"os"
	"strconv"
)

func getOrDefaultString(key, _default string) string {
	value := os.Getenv(key)
	if value == "" {
		return _default
	}
	return value
}
func getOrDefaultInt(key string, _default int) int {
	value := os.Getenv(key)
	if value == "" {
		return _default
	}
	v, _ := strconv.Atoi(value)
	return v
}
func getOrDefaultBool(key string, _default bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return _default
	}
	v, _ := strconv.ParseBool(value)
	return v
}
