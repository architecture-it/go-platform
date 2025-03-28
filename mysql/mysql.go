package mysql

import (
	"os"
	"sync"

	"github.com/architecture-it/go-platform/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func init() {
	once.Do(func() {
		args := os.Getenv("MYSQL_CONNECTION")
		if args != "" {
			conn, err := gorm.Open(mysql.Open(args), &gorm.Config{})
			if err != nil {
				log.Logger.Error(err.Error())
			} else {
				db = conn
				log.Logger.Info("Se ha conectado exitosamente a MySQL.")
			}
		}
	})
}

func GetDB() *gorm.DB {
	return db
}
