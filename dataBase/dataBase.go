package database

import (
	"context"

	"github.com/architecture-it/go-platform/log"
	"github.com/jinzhu/gorm"
	"go.elastic.co/apm/module/apmgorm"
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
	conn := createConnection("mssql", connection)
	return &dataRepository{
		db:         conn,
		dialect:    "mssql",
		connection: connection,
	}
}

func NewMysqlRepository(connection string) DataRepository {
	conn := createConnection("mysql", connection)
	return &dataRepository{
		db:         conn,
		dialect:    "mysql",
		connection: connection,
	}
}

func createConnection(dialect, connectionString string) *gorm.DB {
	if dialect == "" || connectionString == "" {
		log.Logger.Error("No se pudo conectar a la base de datos. Falta informacion!")
		return nil
	}

	connection, err := apmgorm.Open(dialect, connectionString)

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
		log.Logger.Info("Se intenta reconectar a la base de datos.")
		repo.db = createConnection(repo.dialect, repo.connection)
	}
	if repo.db == nil {
		log.Logger.Error("No se pudo conectar a la base de datos.")
		return nil
	}
	db := apmgorm.WithContext(ctx, repo.db)
	return db
}
