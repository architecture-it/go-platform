# web
Helpers para montar servicios de APIs usando (go-gonic)[]https://github.com/gin-gonic/gin

## Uso

``` go get github.com/eandreani/go-platform/web ``` 

## Ejemplos

```go
import (
    "github.com/eandreani/go-platform/web"
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    
    server:= web.NewServer(web.ReadConfigFromEnv()) 
	server.AddMetrics() // /metrics de los requests en formato prometheus 
    server.AddCorsAllOrigins()
    
    //AddHealth(health.NewRedisHealthChecker(redisHealthChecker.Config{}),
    //			health.NewMySqlHealthChecker(mySqlHealthChecker.Config{}),
    //			...func())
    server.AddHealth(web.HealthAlwaysUp) // pone un endpoint en /health

    // apidocs con la documentacion de openApi que se especifique
    server.AddApiDocs("https://raw.githubusercontent.com/eandreani/proyecto/openapi.yaml")

    r := server.GetRouter()
    r.GET("/ping", func (c *gin.Context) {
        c.String(http.StatusOK,"pong")
    })

    server.ListenAndServe()
}
``` 