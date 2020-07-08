# log

```go
import "github.com/architecture-it/go-platform/log"
import "errors"

func main() {

    log.Trace.Pipeline("Hola!")
    // lo imprime de esta forma en stdout:
    // 2020-06-12 14:13:23.041 | -1 | /tmp/go-build575633354/b001/exe/log | TRACE | log.go:18 | Hola!

    err := errors.New("No sos vos soy yo")
    log.Error.JSON(err.Error())
    // {"Date":"2020-06-12 14:13:23.041","Level":"3","Local":"/tmp/go-build575633354/b001/exe/log","Name":"ERROR","Path":"log.go:18","Message":"No sos vos soy yo"}
}
```