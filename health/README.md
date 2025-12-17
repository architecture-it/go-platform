# Health 💊

Este paquete permite evaluar y determinar el estado de salud de nuestra aplicación.

## HealthCheckers disponibles

Por el momento, disponemos de 4 HealthCheckers:

* IBM MQ
* MS SQL
* MySQL
* Redis

En la siguiente tabla se pueden observar las variables de entorno necesarias para evaluar el correspondiente estado de cada una:

|HealthChecker | Variable/s de entorno |
|-|-|
| IBM MQ | HTTP_MQ_API_URL |
| MS SQL | SQL_CONNECTION |
| MySQL | MYSQL_CONNECTION |
| Redis | REDIS_DB, REDIS_ADDR y REDIS_PASS |

### Variables de timeout (opcionales)

| Variable | Default | Descripción |
|----------|---------|-------------|
| REDIS_DIAL_TIMEOUT_MS | 500 | Timeout de conexión en milisegundos |
| REDIS_READ_TIMEOUT_MS | 300 | Timeout de lectura en milisegundos |
| REDIS_WRITE_TIMEOUT_MS | 300 | Timeout de escritura en milisegundos |
| REDIS_SKIP_PING | false | Omitir PING inicial (true/false) |

## Manos a la obra! 👨‍💻👩‍💻

Para agregar los Checkers que se quieran monitorear, alcanza con agregarlos en la función AddHealth del WebServer

```go
import (
	"github.com/architecture-it/go-platform/health"
	"github.com/architecture-it/go-platform/web"
)

func main() {
	server := web.NewServer(web.ReadConfigFromEnv())
	server.AddHealth(health.IbmMQHealthChecker, health.MssqlHealthChecker, health.MysqlHealthChecker, health.RedisHealthChecker)
	api.SetupRouter(server.GetRouter())
	server.ListenAndServe()
}
```

El código anterior dispondrá en el endpoint **/health** una salida similar a la siguiente:

```json
{
	"status": {
		"code": "UP",
		"description": "AlwaysUpEndpoint"
	},
	"details": {
		"ibmQueueHealthIndicator": {
			"status": {
				"code": "UP",
				"description": ""
			},
			"details": ""
		},
		"mysqlHealthIndicator": {
			"status": {
				"code": "UP",
				"description": ""
			},
			"details": {
				"hostname": "localhost",
				"version": "10.0.29-MariaDB-0ubuntu0.16.04.1"
			}
		},
		"redisHealthIndicator": {
			"status": {
				"code": "UP",
				"description": ""
			},
			"details": {
				"address": "127.0.0.1:6379",
				"totalMemory": "3.70G",
				"usedMemory": "1.62G",
				"version": "3.2.3"
			}
		},
		"sqlServerHealthIndicator": {
			"status": {
				"code": "UP",
				"description": ""
			},
			"details": {
				"host": "localhost",
				"version": "15.00.2000"
			}
		}
	}
}
```