package mssql

import (
	"os"

	"github.com/architecture-it/go-platform/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var db *gorm.DB

func init() {
	args := os.Getenv("SQL_CONNECTION")
	if args != "" {
		conn, err := gorm.Open("mssql", args)
		if err != nil {
			log.Logger.Error(err.Error())
		} else {
			db = conn
			log.Logger.Info("Se ha conectado exitosamente.")
		}
	}
}

// GetDB return the database connection
func GetDB() *gorm.DB {
	return db
}
