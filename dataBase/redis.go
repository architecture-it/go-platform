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
	// Leer timeouts desde ENV (en milisegundos)
	dialTimeoutMs := 500 // default 500ms
	if timeoutEnv := os.Getenv("REDIS_DIAL_TIMEOUT_MS"); timeoutEnv != "" {
		if parsed, err := strconv.Atoi(timeoutEnv); err == nil {
			dialTimeoutMs = parsed
		}
	}

	readTimeoutMs := 300 // default 300ms
	if timeoutEnv := os.Getenv("REDIS_READ_TIMEOUT_MS"); timeoutEnv != "" {
		if parsed, err := strconv.Atoi(timeoutEnv); err == nil {
			readTimeoutMs = parsed
		}
	}

	writeTimeoutMs := 300 // default 300ms
	if timeoutEnv := os.Getenv("REDIS_WRITE_TIMEOUT_MS"); timeoutEnv != "" {
		if parsed, err := strconv.Atoi(timeoutEnv); err == nil {
			writeTimeoutMs = parsed
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     pass,
		DB:           db,
		DialTimeout:  time.Millisecond * time.Duration(dialTimeoutMs),
		ReadTimeout:  time.Millisecond * time.Duration(readTimeoutMs),
		WriteTimeout: time.Millisecond * time.Duration(writeTimeoutMs),
	})

	// Opcional: PING configurable
	skipPing := os.Getenv("REDIS_SKIP_PING") == "true"
	if !skipPing {
		_, err := client.Ping().Result()
		if err != nil {
			log.Logger.Info("[REDIS] Error en la conexión : " + err.Error())
			return nil
		}
		log.Logger.Info("[REDIS] Se ha conectado exitosamente")
	} else {
		log.Logger.Info("[REDIS] Cliente Redis creado sin validación de PING")
	}

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
