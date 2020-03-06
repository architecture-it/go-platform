package mssql

import (
	"fmt"
	"os"
	"time"

	"github.com/eandreani/go-platform/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var db *gorm.DB

//Benchmarkf imprime el tiempo que transcurrio en el logger Trace.
//Ejemplo:
// defer Benchmarkf("paso el %s","tiempo")
// imprime: 2019/06/11 17:38:21 log.go:22: paso el tiempo: 1.2121ms
func Benchmark(fmtt string, args ...string) func() {
	started := time.Now()
	return func() {
		fmt.Printf("%s: %s", fmt.Sprintf(fmtt, args), time.Since(started))
	}
}

func init() {
	dialect := os.Getenv("SQL_DRIVER")
	args := os.Getenv("SQL_CONNECTION")
	conn, err := gorm.Open(dialect, args)
	if err != nil {
		log.Error.Println(err.Error())
	} else {
		db = conn
		log.Info.Println("Se ha conectado exitosamente.")
	}

}

// GetDB return the database
func GetDB() *gorm.DB {
	return db
}
