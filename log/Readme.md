# Log

Este paquete basado en [zap](https://github.com/uber-go/zap)⚡ permite loguear los eventos y/o errores que se dan en nuestra aplicación.
Actualmente, estamos usando el stack de Elastic o ELK, con lo cual, nos basamos en el [estandar de arquitectura](https://github.com/architecture-it/elk-stack-logs-config/tree/master/filebeat).

## Console o JSON? 🔧

Usando la variable de entorno **LOG_CONFIG** podemos configurar si queremos que nuestros logs se visualicen en forma de consola con sus campos separados por pipes (|) o en json, y las keys de los diferentes campos de ese JSON.

### Estándares de configuracion

#### Console

```sh
export LOG_CONFIG = "{\"level\":\"debug\",\"encoding\":\"console\",\"development\":true,\"outputPaths\":[\"stderr\"],\"errorOutputPaths\":[\"stderr\"],\"encoderConfig\":{\"callerKey\":\"context\",\"timeKey\":\"timestamp\",\"messageKey\":\"message\",\"levelKey\":\"severity\",\"stacktraceKey\":\"\"}}"
```

#### JSON

```sh
export LOG_CONFIG = "{\"level\":\"debug\",\"encoding\":\"json\",\"development\":true,\"outputPaths\":[\"stderr\"],\"errorOutputPaths\":[\"stderr\"],\"encoderConfig\":{\"callerKey\":\"context\",\"timeKey\":\"timestamp\",\"messageKey\":\"message\",\"levelKey\":\"severity\",\"stacktraceKey\":\"\"}}"
```

* PD: Si nos olvidamos de setear dicha variable de entorno, el logger se configurará con una salida en formato **console**

## Acceso

El paquete dispone dos variables globales que permiten loguear:
* Logger
* SugarLogger

### Logger o SugarLogger? 🤔

**Logger** permite un logueo rápido recibiendo un parámetro de tipo String, en cualquiera de sus niveles de logueo.
En caso de que no queramos estar concatenando para obtener un parámetro único de String, podemos optar por **SugarLogger**, que si bien es un poco más lento que el **Logger** resuelve ciertas cosas por nosotros, es decir, imprime más al estilo de un printf(), por ejemplo:

```go
    log.SugarLogger.Infof("La conexión a la BD %s se realizó sin problemas!", "BD_PROD_01")
    // o bien
    log.SugarLogger.Error("Error al leer la respuesta desde la API.", err)
```

## Niveles de Logueo 🤓

Ambos Loggers admiten los siguientes niveles de Logueo:

* DEBUG: trazas de la aplicación al debuguear
* INFO: Información general
* WARN: Advertencia!
* ERROR: Errores en general, sean de conexión a una base de datos, a una API, un error al parsear un dato, etc.
* FATAL: Errores irrecuperables  

## Manos a la obra! 👨‍💻👩‍💻

### Console

```sh
export LOG_ENCODING="console"
```

```go
import "github.com/architecture-it/go-platform/log"
import "errors"

func main() {

    log.Logger.Info("Hola!")
    // genera esta salida
    // 2020-08-14 10:31:34.613	 | 0 | INFO | main.exe	| test-log/main.go:6 |	Hola!

    err := errors.New("No sos vos, soy yo")
    log.Error.JSON(err.Error())
    //2020-08-14 11:19:16.524	 | 0 | ERROR | main.exe	| test-log/main.go:11 |	No sos vos, soy yo
}
```

### JSON

```sh
export LOG_ENCODING="json"
```

```go
import "github.com/architecture-it/go-platform/log"
import "errors"

func main() {

    log.Logger.Info("Hola!")
    // genera esta salida
    // {"severity":"INFO","timestamp":"2020-08-17 11:34:15.064","context":"test-log/main.go:6","message":"Hola!","threadId":0,"applicationName":"main.exe"}

    err := errors.New("No sos vos, soy yo")
    log.Error.JSON(err.Error())
    // {"severity":"ERROR","timestamp":"2020-08-14 11:31:19.773","context":"test-log/main.go:11","message":"No sos vos, soy yo","threadId":0,"applicationName":"main.exe"}
}
```