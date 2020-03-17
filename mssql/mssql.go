package mssql

import (
	"os"

	"github.com/andreani-publico/go-platform/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var db *gorm.DB

func init() {
	args := os.Getenv("SQL_CONNECTION")
	conn, err := gorm.Open("mssql", args)
	if err != nil {
		log.Error.Println(err.Error())
	} else {
		db = conn
		log.Info.Println("Se ha conectado exitosamente.")
	}

}

// GetDB return the database connection
func GetDB() *gorm.DB {
	return db
}
