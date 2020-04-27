# log

```go
import "github.com/architecture-it/go-platform/log"
import "errors"

func main() {

    log.Info.Println("Hola!")
    // lo imprime de esta forma en stdout:
    // INFO: 2019/06/05 16:11:46 main.go:9: Hola!

    err := errors.New("No sos vos soy yo")
    log.Error.Printf("Hay un problema: %s",err)
    //etc.
}
```