package handlers

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

var lock = &sync.Mutex{}

type Config struct {
	Http struct {
		Port int `env:"HTTP_PORT" env-default:"8080"`
	}
	Database struct {
		Host       string `env:"DATABASE_HOST" env-default:"db"`
		Port       int    `env:"DATABASE_PORT" env-default:"5432"`
		Database   string `env:"DATABASE_DBNAME" env-default:"postgres"`
		Username   string `env:"DATABASE_USERNAME" env-default:"postgres"`
		Password   string `env:"DATABASE_PASSWORD" env-default:"postgres"`
		DisableTLS bool   `env:"DATABASE_DISABLE_TLS" env-default:"true"`
		CetrPath   string `env:"DATABASE_CERT_PATH"`
	}
	Minio struct {
		Endpoint        string `env:"MINIO_ENDPOINT" env-default:"storage"`
		AccessTokenID   string `env:"MINIO_ACCESS_TOKEN_ID" env-default:"storage-access-key"`
		SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY" env-default:"storage-secret-key"`
		Region          string `env:"MINIO_REGION"`
		UseSSL          bool   `env:"MINIO_USE_SSL" env-default:"false"`
	}
}

var config *Config

//GetConfig load config from ENVs or .env file. If config is already loaded, return it.
func getConfig() (*Config, error) {
	if config == nil {
		lock.Lock()
		defer lock.Unlock()

		if config == nil {
			var cfg Config
			if err := cleanenv.ReadEnv(&cfg); err != nil {
				return nil, errors.Wrap(err, "api: Read ENVs error")
			}
			config = &cfg
		}
	}
	return config, nil
}
