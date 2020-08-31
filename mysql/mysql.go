package mysql

import (
	"os"

	"github.com/architecture-it/go-platform/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	args := os.Getenv("MYSQL_CONNECTION")
	conn, err := gorm.Open("mysql", args)
	if err != nil {
		log.Logger.Error(err.Error())
	} else {
		db = conn
		log.Logger.Info("Se ha conectado exitosamente.")
	}
}

// GetDB return the database connection
func GetDB() *gorm.DB {
	return db
}
