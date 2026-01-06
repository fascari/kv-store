package config

import (
	"strconv"

	"github.com/felipeascari/kv-store/pkg/environment"
	"github.com/felipeascari/kv-store/pkg/storage"
)

type (
	Config struct {
		Storage StorageConfig
		Server  ServerConfig
	}

	StorageConfig struct {
		Type  storage.Type
		Redis RedisConfig
	}

	RedisConfig struct {
		Addr     string
		Password string
		DB       int
	}

	ServerConfig struct {
		Port string
	}
)

func Load() (*Config, error) {
	redisDB, _ := strconv.Atoi(environment.LoadEnv("REDIS_DB", "0"))

	return &Config{
		Storage: StorageConfig{
			Type: storage.Type(environment.LoadEnv("STORAGE_TYPE", storage.TypeRedis.String())),
			Redis: RedisConfig{
				Addr:     environment.LoadEnv("REDIS_ADDR", "localhost:6379"),
				Password: environment.LoadEnv("REDIS_PASSWORD", ""),
				DB:       redisDB,
			},
		},
		Server: ServerConfig{
			Port: environment.LoadEnv("SERVER_PORT", "8080"),
		},
	}, nil
}
