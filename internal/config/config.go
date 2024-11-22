package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	HttpServer     HttpServerConfig
	PostgresConfig PostgresConfig
}

func MustLoad() *AppConfig {
	var cfg AppConfig

	_ = godotenv.Load()

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
