package database

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/architecture-it/go-platform/log"
	"github.com/go-redis/redis"
	"go.elastic.co/apm/module/apmgoredis"
)

type RedisRepository interface {
	GetClient(ctx context.Context) *redis.Client
}

type redisRepository struct {
	client *redis.Client
	addr   string
	pass   string
	db     int
}

func NewRedisRepository(addr, pass, db string) RedisRepository {
	dbRedis, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Logger.Info("[REDIS] Se eligió una DB por default")
	}
	client := createConnectionRedis(addr, pass, dbRedis)
	return &redisRepository{
		client: client,
		addr:   addr,
		db:     dbRedis,
		pass:   pass,
	}
}

func createConnectionRedis(addr, pass string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    pass,
		DB:          db,
		DialTimeout: time.Millisecond * 50,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Logger.Info("[REDIS] Error en la prueba de conexión : " + err.Error())
	}
	log.Logger.Info("[REDIS] Se ha conectado exitosamente")
	return client
}

// GetDB return the database connection
func (repo *redisRepository) GetClient(ctx context.Context) *redis.Client {
	// reintento de conexion si algo fallo
	if repo.client == nil {
		log.Logger.Info("Se intenta reconectar a redis.")
		repo.client = createConnectionRedis(repo.addr, repo.pass, repo.db)
	}
	if repo.client == nil {
		return nil
	}
	clientAPM := apmgoredis.Wrap(repo.client).WithContext(ctx)
	return clientAPM.RedisClient()
}
