package core_logger_repository_clickhouse

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	host     string        `envconfig:"HOST" default:"localhost:9000"`
	database string        `envconfig:"DB" default:"logs"`
	user     string        `envconfig:"USER" default:"default"`
	password string        `envconfig:"PASSWORD"`
	service  string        `envconfig:"SERVICE_NAME" default:"my-app"`
	timeout  time.Duration `envconfig:"TIMEOUT" default:"5s"`
}

func NewConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("CLICKHOUSE", &cfg); err != nil {
		return config{}, fmt.Errorf("failed process envconfig: %w", err)
	}

	return cfg, nil
}
