## errores


Este paquete exporta algunas funcionespara estandarizar los errores

Ejemplo:
```go
import "github.com/architecture-it/go-platform/errores"
import "context"

func main() {

  err := context.BindJSON(&bodyRequest)
	if err != nil {
   context.JSON(http.StatusBadRequest, errores.ErrorResponse.Default("Detalle", err))
  }
  

}

```
