## errores


Este paquete exporta algunas funcionespara estandarizar los errores

Ejemplo:
```go
import "github.com/architecture-it/go-platform/errores"
import "context"

func main() {

  err := ctx.BindJSON(&bodyRequest)
	if err != nil {
   ctx.JSON(http.StatusBadRequest, errores.PedidoIncorrecto("detalle", listaDeErrores))
  }
  

}

```
