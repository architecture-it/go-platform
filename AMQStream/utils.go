package AMQStream

import (
	"strconv"
)

func getOrDefaultString(config map[string]string, key, _default string) string {
	value := config[key]
	if value == "" {
		return _default
	}
	return value
}
func getOrDefaultInt(config map[string]string, key string, _default int) int {
	value := config[key]
	if value == "" {
		return _default
	}
	v, _ := strconv.Atoi(value)
	return v
}
func getOrDefaultBool(config map[string]string, key string, _default bool) bool {
	value := config[key]
	if value == "" {
		return _default
	}
	v, _ := strconv.ParseBool(value)
	return v
}
