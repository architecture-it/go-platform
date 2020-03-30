
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

Sino se puede hacer de la siguiente forma:

type Ejemplo struct {
	Dato1 string `sql:"column:nombreDelDatoEnBBDD"` 
	// Si no se agrega la etiqueta sql:"column:nombre" toma por defecto el nombre del struct con minuscula
}

//TableName renombro struct por el el de la base de datos
func (Ejemplo) TableName() string {
	return os.Getenv("VariableConElNombreDeLaTablaEjemplo")
}

func FindByConditionOpcion2() struct string {
	var struct []Ejemplo 
	//En el find buscara el nombre de la tabla Ejemplo por la Funcion TableName() y devuelve todos los datos del struct
    	GetDB().Where("column1 IS NOT NULL").Where("column2 IS NOT NULL").Find(&struct)
   	return struct
}
