package mssql

import (
	"os"
	"sync"

	"github.com/architecture-it/go-platform/log"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func init() {
	once.Do(func() {
		args := os.Getenv("SQL_CONNECTION")
		if args != "" {
			conn, err := gorm.Open(sqlserver.Open(args), &gorm.Config{})
			if err != nil {
				log.Logger.Error(err.Error())
			} else {
				db = conn
				log.Logger.Info("Se ha conectado exitosamente a SQL Server.")
			}
		}
	})
}

func GetDB() *gorm.DB {
	return db
}
