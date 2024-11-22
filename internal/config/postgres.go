package config

import (
	"fmt"
	"github.com/orungrau/em_song_library/internal/repository/storage/song"
)

type PostgresConfig struct {
	Host       string `env:"POSTGRES_HOST" env-required:"true"`
	Port       int    `env:"POSTGRES_PORT" env-required:"true"`
	User       string `env:"POSTGRES_USER" env-required:"true"`
	Password   string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DisableSSL bool   `env:"POSTGRES_DISABLE_SSL" env-required:"true"`
	Database   string `env:"POSTGRES_DATABASE" env-required:"true"`

	MigrationSource string `env:"POSTGRES_MIGRATION_SOURCE" env-required:"true"`
}

func NewPostgresConfig() song.PostgresStorageConfig {
	return &PostgresConfig{}
}

func (p *PostgresConfig) GetConnection() string {
	sslMode := "require"
	if p.DisableSSL {
		sslMode = "disable"
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.Database, sslMode)

	return connectionString
}

func (p *PostgresConfig) GetDatabase() string {
	return p.Database
}

func (p *PostgresConfig) GetMigrationSource() string {
	return p.MigrationSource
}
