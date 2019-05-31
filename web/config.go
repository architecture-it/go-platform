package web

import (
	"os"
)
//Config para los parametros de config del server
type Config struct {

	Port string
	
}

//DefaultConfig crea una Config default que escucha en 8080
func DefaultConfig() Config {

	return Config{"8080"}

}

//ReadConfigFromEnv lee la config de las vars de entorno.
//$PORT: el puerto donde se expone el server
func ReadConfigFromEnv() Config {
	config := DefaultConfig()
	port := os.Getenv("PORT")
	if port != "" {
		config.Port = port
	}

	return config
}