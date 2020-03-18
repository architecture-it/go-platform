package mq

import (
	"os"
	"github.com/andreani-publico/go-platform/log"
	"github.com/andreani-publico/go-platform/vault"
)
type Config struct {
	HTTPMQAPIUrl string
	QueueName string
} 

const (
	MQ_API_CONFIG_KEY = "HTTP_MQ_API_URL"
	QUEUE_NAME_KEY = "QUEUE_NAME"
)

//ReadConfigFromEnv lee la config de las var de entorno
//$HTTP_MQ_API_URL: el url a httpmqbridge, $QUEUE_NAME: la cola donde se quiere sacar y poner
func ReadConfigFromEnv() Config {
	return Config{os.Getenv(MQ_API_CONFIG_KEY),os.Getenv(QUEUE_NAME_KEY)}
}

//ReadConfigFromVault lee la config desde el ConfigVault.
func ReadConfigFromVault(v vault.Vault) Config {
	mq,err := v.Get(MQ_API_CONFIG_KEY)
	if err != nil {
		log.Fatal.Printf("can't read vault. key:%s. error:%s. Trying to read form env..",MQ_API_CONFIG_KEY,err)
		return ReadConfigFromEnv()
	}
	queue,err := v.Get(QUEUE_NAME_KEY)
	if err != nil {
		log.Fatal.Printf("can't read vault. key:%s. error:%s. Trying to read form env..",QUEUE_NAME_KEY,err)
		return ReadConfigFromEnv()
	}
	return Config{
		mq,
		queue}
}