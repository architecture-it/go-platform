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
	if err != nil {
		panic("No se encontró la variable de entorno REDIS_DB.")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       db,
	})
}

func RedisHealthChecker() Status {
	status := StatusResult(UP)
	_, err := client.Ping().Result()
	if err != nil {
		status = DOWN
	}
	infoResponse, err := client.Info().Result()
	info := parseInfo(infoResponse)
	result := make(map[string]interface{})
	result["Address"] = client.Options().Addr
	result["Version"] = info["redis_version"]
	result["Used Memory"] = info["used_memory_human"]
	result["Total Memory"] = info["total_system_memory_human"]
	return Status{"RedisHealthIndicator", status, result}
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
