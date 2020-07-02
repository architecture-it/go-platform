# xsd
Este paquete exporta las estructuras del xsd, así como funciones para manejar los eventos correctamente. Se debe usar GetDate y GetTime para los campos de tipo Date y Time.
No se validan los campos obligatorios.
Las estructuras se actualizan desde Integraciones.Esquemas.

Ejemplo:

```go
import "github.com/architecture-it/go-platform/xsd"

func main() {

    ordenDeEnvio := xsd.OrdenDeEnvioSolicitada{Remitente: "APIs", Cuando: xsd.GetDate(time.Now())}
	str, err := xsd.GetEvento(&ordenDeEnvio, false)
	if err != nil {
		panic(err)
	}
}
```