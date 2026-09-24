package core_pgx_pool

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"5s"`
}

func NewConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("POSTGRES", &cfg); err != nil {
		return config{}, fmt.Errorf("failed process envconfig: %w", err)
	}

	return cfg, nil
}

func NewConfigMast() config {
	cfg, err := NewConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[bootstrap] get Postgres connection pool config: %v\n return default config", err)
		return config{
			Host:     "localhost:5432",
			User:     "User",
			Password: "pass",
			Database: "db",
			Timeout:  5 * time.Second,
		}
	}
	return cfg
}
