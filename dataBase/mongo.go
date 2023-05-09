package database

import (
	"context"

	"github.com/architecture-it/go-platform/log"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository interface {
	GetDB(ctx context.Context) *mongo.Database
}

type mongoRepository struct {
	mongoURL string
	mongoDB  string
	db       *mongo.Database
}

func NewMongoRepository(mongoURL, mongoDB string) MongoRepository {
	db := createConnectionMongo(mongoURL, mongoDB)
	return &mongoRepository{
		mongoURL: mongoURL,
		mongoDB:  mongoDB,
		db:       db,
	}
}

func createConnectionMongo(mongoURL, mongoDB string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(mongoURL).
		SetMonitor(apmmongo.CommandMonitor())
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil || mongoDB == "" {
		log.Logger.Error("Error Fatal NO se conectó a MongoDB, la url es: " + mongoURL +
			" , la BBDD es: " + mongoDB + " . Verifique que sean correctos")
	}

	err = mongoClient.Ping(context.TODO(), nil)

	if err != nil {
		log.Logger.Error("PING a MongoDB error: " + err.Error())
		return nil
	}
	log.Logger.Info("Se conectó con éxito a MongoDB")
	db := mongoClient.Database(mongoDB)

	return db
}

// GetDB return the database connection
func (repo *mongoRepository) GetDB(ctx context.Context) *mongo.Database {
	// reintento de conexion si algo fallo
	if repo.db == nil {
		log.Logger.Info("Se intenta reconectar a mongo.")
		repo.db = createConnectionMongo(repo.mongoURL, repo.mongoDB)
	}
	if repo.db == nil {
		return nil
	}
	return repo.db
}
