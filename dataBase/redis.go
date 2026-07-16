package database

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/architecture-it/go-platform/log"
	"github.com/go-redis/redis"
	"go.elastic.co/apm/module/apmgoredis"
)

const (
	// Default Redis timeout values in milliseconds
	defaultDialTimeoutMs  = 500
	defaultReadTimeoutMs  = 300
	defaultWriteTimeoutMs = 300
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
	dialTimeoutMs := getTimeoutFromEnv("REDIS_DIAL_TIMEOUT_MS", defaultDialTimeoutMs)
	readTimeoutMs := getTimeoutFromEnv("REDIS_READ_TIMEOUT_MS", defaultReadTimeoutMs)
	writeTimeoutMs := getTimeoutFromEnv("REDIS_WRITE_TIMEOUT_MS", defaultWriteTimeoutMs)

	options := &redis.Options{
		Addr:         addr,
		Password:     pass,
		DB:           db,
		DialTimeout:  time.Millisecond * time.Duration(dialTimeoutMs),
		ReadTimeout:  time.Millisecond * time.Duration(readTimeoutMs),
		WriteTimeout: time.Millisecond * time.Duration(writeTimeoutMs),
	}

	tlsEnabled, err := strconv.ParseBool(
		strings.TrimSpace(os.Getenv("REDIS_TLS_ENABLED")),
	)
	if err != nil && os.Getenv("REDIS_TLS_ENABLED") != "" {
		log.Logger.Info(
			"[REDIS] REDIS_TLS_ENABLED tiene un valor inválido. Se utilizará false",
		)
		tlsEnabled = false
	}

	if tlsEnabled {
		serverName := strings.TrimSpace(
			os.Getenv("REDIS_TLS_SERVER_NAME"),
		)

		// Si no se informó explícitamente, se obtiene desde REDIS_ADDR.
		if serverName == "" {
			host, _, splitErr := net.SplitHostPort(addr)
			if splitErr != nil {
				log.Logger.Info(
					"[REDIS] No se pudo obtener el hostname TLS desde REDIS_ADDR: " +
						splitErr.Error(),
				)
			} else {
				serverName = host
			}
		}

		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: serverName,
		}

		log.Logger.Info(
			"[REDIS] TLS habilitado. ServerName=" + serverName,
		)
	} else {
		log.Logger.Info("[REDIS] TLS deshabilitado")
	}

	client := redis.NewClient(options)

	skipPing := strings.EqualFold(
		os.Getenv("REDIS_SKIP_PING"),
		"true",
	)

	if !skipPing {
		_, err := client.Ping().Result()
		if err != nil {
			log.Logger.Info(
				"[REDIS] Error en la conexión: " + err.Error(),
			)

			_ = client.Close()
			return nil
		}

		log.Logger.Info("[REDIS] Se ha conectado exitosamente")
	} else {
		log.Logger.Info(
			"[REDIS] Cliente Redis creado sin validación de PING",
		)
	}

	return client
}

// getTimeoutFromEnv retrieves timeout value from environment variable or returns default
func getTimeoutFromEnv(envVar string, defaultValue int) int {
	if timeoutEnv := os.Getenv(envVar); timeoutEnv != "" {
		if parsed, err := strconv.Atoi(timeoutEnv); err == nil {
			return parsed
		}
	}
	return defaultValue
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
