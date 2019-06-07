## mq


Este paquete exporta algunas apis para sacar y poner de las colas de MQ y para publicar. Para hacerlo no usa el protocolo de IBM/MQ si no que llama a http-mq-bridge.

```go
import "github.com/eandreani/go-platform/mq"
import "context"

func main() {

    //ReadConfigFromEnv() lee del environment HTTP_MQ_API_URL y QUEUE_NAME
    //Tambien se puede leer la config del Vault usando ReadConfigFromVault().
    q := mq.GetQueue(mq.ReadConfigFromEnv()) 

    err := q.Put("Hola")

    ctx,cancel := context.WithTimeout(context.Background(),1*time.Second)
    defer cancel()
    //Listen() ejecuta la goroutine del closure cada vez que llega un mensaje hasta que el contexto (ctx) se cancele.
	q.Listen(ctx,func(what string) {
        //en what deberia venir "Hola"
    })
    <-ctx.Done()

    topic := mq.GetTopic("topic/subtopic",mq.ReadConfigFromEnv())
    err = topic.Publish("mensaje")
}

```

## Roadmap
Ver si mejor no leer la configracion de forma global:
```go
mq.ReadConfigFromEnv()
mq.ReadConfigFromVault()
```
Y luego cambian las apis a:
```go
q:= mq.GetQueue("miQueue")
q.Put("mensaje")
t:=mq.GetTopic("untopic/otro")
t.Publish("un mensaje")
```