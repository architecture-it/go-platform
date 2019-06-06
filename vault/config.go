package vault

import (
	"os"
)

type Config struct {

	Passphrase string
	VaultUrl string
	Namespace string
}

func ReadConfigFromEnv() *Config {
	return &Config{
		os.Getenv("VAULT_PASSPHRASE"),
		os.Getenv("VAULT_URL"),
		os.Getenv("APP_NAMESPACE")}
}