package mq

import (
	"os"
)
type Config struct {
	HTTPMQAPIUrl string
	QueueName string
} 

//ReadConfigFromEnv lee la config de las var de entorno
//$HTTP_MQ_API_URL: el url a httpmqbridge, $QUEUE_NAME: la cola donde se quiere sacar y poner
func ReadConfigFromEnv() Config {
	return Config{os.Getenv("HTTP_MQ_API_URL"),os.Getenv("QUEUE_NAME")}
}