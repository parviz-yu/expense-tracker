package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	HTTPServer
	Database
}

type HTTPServer struct {
	Port        string        `env:"SERVER_PORT" env-default:"8080"`
	Timeout     time.Duration `env:"SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"SERVER_IDLETIMEOUT" env-default:"45s"`
}

type Database struct {
	PostgresHost     string `env:"POSTGRES_HOST" env-required:"true"`
	PostgresPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDatabase string `env:"POSTGRES_DATABASE" env-default:"postgres"`
}

func MustLoad() Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return cfg
}
