package database

import (
	"context"
	"fmt"

	"github.com/architecture-it/go-platform/log"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository interface {
	GetDB(ctx context.Context) *mongo.Database
}

type mongoRepository struct {
	mongoURL       string
	mongoDB        string
	mongoMechanism string
	mongoUser      string
	mongoPass      string
	db             *mongo.Database
}

func NewMongoRepository(mongoURL, mongoDB, mongoMechanism, mongoUser, mongoPass string) MongoRepository {
	db := createConnectionMongo(mongoURL, mongoDB, mongoMechanism, mongoUser, mongoPass)
	return &mongoRepository{
		mongoURL:       mongoURL,
		mongoDB:        mongoDB,
		mongoMechanism: mongoMechanism,
		mongoUser:      mongoUser,
		mongoPass:      mongoPass,
		db:             db,
	}
}

func createConnectionMongo(mongoURL, mongoDB, mongoMechanism, mongoUser, mongoPass string) *mongo.Database {

	var clientOptionsAuth *options.ClientOptions

	clientOptions := options.Client().ApplyURI(mongoURL).
		SetMonitor(apmmongo.CommandMonitor())

	if mongoMechanism != "" && mongoUser != "" && mongoPass != "" {
		clientOptionsAuth = options.Client().SetAuth(options.Credential{
			AuthMechanism: mongoMechanism,
			AuthSource:    mongoDB,
			Username:      mongoUser,
			Password:      mongoPass,
		})
	}

	mongoClient, err := mongo.Connect(context.TODO(), clientOptions, clientOptionsAuth)
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

// // GetDB return the database connection
// func (repo *mongoRepository) GetDB(ctx context.Context) *mongo.Database {
// 	// reintento de conexion si algo fallo
// 	if repo.db == nil {
// 		log.Logger.Info("Se intenta reconectar a mongo.")
// 		repo.db = createConnectionMongo(repo.mongoURL, repo.mongoDB, repo.mongoMechanism, repo.mongoUser, repo.mongoPass)
// 	}
// 	if repo.db == nil {
// 		return nil
// 	}
// 	return repo.db
// }

const (
	maxRetries = 3
)

func (repo *mongoRepository) GetDB(ctx context.Context) *mongo.Database {
	for i := 0; i < maxRetries; i++ {
		if repo.db == nil {
			log.Logger.Info("Intentando reconectar a MongoDB...")
			repo.db = createConnectionMongo(repo.mongoURL, repo.mongoDB, repo.mongoMechanism, repo.mongoUser, repo.mongoPass)
		}

		if repo.db != nil {
			err := repo.db.Client().Ping(ctx, nil)
			if err == nil {
				return repo.db
			}

			log.Logger.Error(fmt.Sprintf("Error de conexión a MongoDB: %s", err))

			// Si hay un error de conexión, desconecta para forzar una reconexión en el próximo intento
			if err := repo.db.Client().Disconnect(ctx); err != nil {
				log.Logger.Info(fmt.Sprintf("Error de conexión a MongoDB: %s", err.Error()))
			}
			repo.db = nil
		}

		log.Logger.Info("Intentando nuevamente la conexión a MongoDB")

	}

	log.Logger.Error("No se pudo establecer la conexión después de varios intentos.")
	return nil
}
