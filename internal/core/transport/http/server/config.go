package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Addr            string        `envconfig:"ADDR" required:"true"` //адрес на котором запускаем сервер
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
}

func NewConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("HTTP", &cfg); err != nil {
		return config{}, fmt.Errorf("failed process envconfig: %w", err)
	}
	return cfg, nil
}

func NewConfigMust() config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("Failed get logger config: %w", err)
		panic(err)
	}
	return config
}
