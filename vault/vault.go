package vault

import (
	"gopkg.in/resty.v1"
	"fmt"
)
type Vault struct {
	config *Config
}

func GetVault(cfg *Config) Vault {
	return Vault{cfg}
}

func (v Vault) buildUrl(key string) string {
	return fmt.Sprintf("%s/keys/%s?passphrase=%s",v.config.VaultUrl,key,v.config.Passphrase)
}

func (v Vault) Get(key string) (string,error) {

	resp,err := resty.R().Get(v.buildUrl(key))
	return resp.String(),err	
}