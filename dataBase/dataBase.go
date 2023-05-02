package database

import (
	"context"

	"github.com/architecture-it/go-platform/log"
	mysql "go.elastic.co/apm/module/apmgormv2/driver/mysql"
	sql "go.elastic.co/apm/module/apmgormv2/v2/driver/sqlserver"
	"gorm.io/gorm"
)

type DataRepository interface {
	GetDB(ctx context.Context) *gorm.DB
}

type dataRepository struct {
	dialect    string
	connection string
	db         *gorm.DB
}

func NewDataRepository(dialect, connection string) DataRepository {
	conn := createConnection(dialect, connection)
	return &dataRepository{
		db:         conn,
		dialect:    dialect,
		connection: connection,
	}
}

func createConnection(dialect, connectionString string) *gorm.DB {
	if dialect == "" || connectionString == "" {
		return nil
	}
	var dialector gorm.Dialector

	if dialect == "sql" || dialect == "mssql" {
		dialector = sql.Open(connectionString)
	} else if dialect == "mysql" {
		dialector = mysql.Open(connectionString)
	} else {
		return nil
	}

	connection, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		log.Logger.Error("Error al conectarse a la base de datos " + dialect + " Descripcion: " + err.Error())
		return nil
	}
	log.Logger.Info("Se conectó con éxito a " + dialect)
	return connection
}

// GetDB return the database connection
func (repo *dataRepository) GetDB(ctx context.Context) *gorm.DB {
	// reintento de conexion si algo fallo
	if repo.db == nil {
		repo.db = createConnection(repo.dialect, repo.connection)
	}
	if repo.db == nil {
		return nil
	}
	db := repo.db.WithContext(ctx)
	return db
}
