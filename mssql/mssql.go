package mssql

import (
	"os"

	"github.com/eandreani/go-platform/log"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

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

// GetDB return the database connection
func GetDB() *gorm.DB {
	return db
}
