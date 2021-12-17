package health

import (
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err == nil && os.Getenv("REDIS_ADDR") != "" {
		client = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASS"),
			DB:       db,
		})
	}
}

func RedisHealthChecker(queue string) func() Checker {
	return func() Checker {
		status := UP
		_, err := client.Ping().Result()
		if err != nil {
			status = DOWN
		}
		infoResponse, err := client.Info().Result()
		info := parseInfo(infoResponse)
		result := make(map[string]interface{})
		result["address"] = client.Options().Addr
		result["version"] = info["redis_version"]
		result["usedMemory"] = info["used_memory_human"]
		result["totalMemory"] = info["total_system_memory_human"]
		if queue != "" {
			result["Queue "+queue] = RedisCheckLenQueue(queue)
		}
		return Checker{Health: Health{Status: Status{Code: status, Description: ""}, Details: result}, Name: "redisHealthIndicator"}
	}
}

func parseInfo(in string) map[string]string {
	info := map[string]string{}
	lines := strings.Split(in, "\r\n")

	for _, line := range lines {
		values := strings.Split(line, ":")

		if len(values) > 1 {
			info[values[0]] = values[1]
		}
	}
	return info
}

func RedisCheckLenQueue(key string) interface{} {
	result := client.LLen(key)
	val, err := result.Result()
	info := map[string]interface{}{"key": key,
		"len":   val,
		"extra": err,
	}
	return info
}
