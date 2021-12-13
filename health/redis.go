package health

import (
	"fmt"
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

func RedisHealthChecker(fn RedisKeyCheck, key string) Checker {
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
	if fn != nil {
		result["dataOfKey"] = fn(key)
	}
	return Checker{Health: Health{Status: Status{Code: status, Description: ""}, Details: result}, Name: "redisHealthIndicator"}
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

type RedisKeyCheck func(clave string) interface{}

func CheckLenQueue(clave string) interface{} {
	result := client.LLen(clave)
	fmt.Println("RESULTADO == ", result)
	val, err := result.Result()
	info := map[string]interface{}{"key": clave,
		"len":   val,
		"extra": err,
	}
	return info
}
