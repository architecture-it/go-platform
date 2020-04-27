# vault
## Uso

```go get "github.com/architecture-it/go-platform/vault"```

Este paquete expone las apis para poder acceder al Vault de la configuracion

```go
import "github.com/architecture-it/go-platformm/vault"

func main() {
    //La configracion del Vault se puede leer del entorno. Deben estar definidas las siguientes variables: VAULT_PASSPHRASE, VAULT_URL y APP_NAMESPACE

    v = vault.GetVault(vault.ReadConfigFromEnv())
    value := v.Get("una-clave")

    //la config en el vault esta encriptada, para eso es el PASSPHRASE.
    //APP_NAMESPACE es el namespace de la app, tipicamente el nombre del repositorio de GitHub
    //Con el namspace mas la clave se accede al vault: go-platform.una-clave, por ejemplo.

}

```
Hay algunas apps de go-platform que se pueden configurar con el vault.

```go
import (
    "github.com/architecture-it/go-platform/vault"
    "github.com/architecture-it/go-platform/mq"
)

func main() {
      v = vault.GetVault(vault.ReadConfigFromEnv())
      queue:=mq.GetQueue(mq.ReadConfigFromVault(v))
}
```

## Roadmap

* Realizar un front-end para que se puedan agregar configuraciones al vault. Por el momento solo existe el api para hacerlo.
