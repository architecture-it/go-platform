
## mssql


Este paquete se permite la conexion a una base de datos MSSql. Se deben configurar la variable de entorno: SQL_CONNECTION


```go
import "github.com/andreani-publico/go-platform/mssql"
import "context"

func FindByCondition() struct string {
	var struct []string
	table := os.Getenv("TABLE_TEST")
    GetDB().Table(table).Select(table).Where("column1 IS NOT NULL").Where("column2 IS NOT NULL").Find(&struct)
    return struct
}
