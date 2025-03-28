package database

import (
	"context"
	"fmt"

	"github.com/architecture-it/go-platform/log"
	apmmysql "go.elastic.co/apm/module/apmgormv2/v2/driver/mysql"
	apmssql "go.elastic.co/apm/module/apmgormv2/v2/driver/sqlserver"
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

func NewSQLRepository(connection string) DataRepository {
	return NewDataRepository("mssql", connection)
}

func NewMysqlRepository(connection string) DataRepository {
	return NewDataRepository("mysql", connection)
}

func createConnection(dialect, connectionString string) *gorm.DB {
	if dialect == "" || connectionString == "" {
		log.Logger.Error("No se pudo conectar a la base de datos. Falta información!")
		return nil
	}

	var dialector gorm.Dialector
	switch dialect {
	case "mssql":
		dialector = apmssql.Open(connectionString)
	case "mysql":
		dialector = apmmysql.Open(connectionString)
	default:
		log.Logger.Error(fmt.Sprintf("Dialect '%s' no soportado", dialect))
		return nil
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Logger.Error("Error al conectarse a la base de datos: " + err.Error())
		return nil
	}

	log.Logger.Info("Se conectó con éxito a " + dialect)
	return db
}

func (repo *dataRepository) GetDB(ctx context.Context) *gorm.DB {
	if repo.db == nil {
		log.Logger.Info("Intentando reconectar a la base de datos.")
		repo.db = createConnection(repo.dialect, repo.connection)
	}
	if repo.db == nil {
		log.Logger.Error("No se pudo conectar a la base de datos.")
		return nil
	}
	return repo.db.WithContext(ctx)
}
