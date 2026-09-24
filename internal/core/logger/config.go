package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Level  string `envconfig:"LEVEL" default:"DEBUG"`
	Folder string `envconfig:"FOLDER"`
}

func NewConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("LOG", &cfg); err != nil {
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
