package mq

import (
	"os"
)
type Config struct {
	HTTPMQAPIUrl string
	QueueName string
} 

func ReadConfigFromEnv() Config {
	return Config{os.Getenv("HTTP_MQ_API_URL)"),os.Getenv("QUEUE_NAME")}
}