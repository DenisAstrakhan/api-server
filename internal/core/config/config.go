package core_config

import (
	"os"
	"time"
)

type Config struct {
	TimeZone *time.Location
}

func NewConfig() *Config {
	timeZone := os.Getenv("TIME_ZONE")
	if timeZone == "" {
		timeZone = "UTC"
	}
	zone, err := time.LoadLocation(timeZone)
	if err != nil {
		zone = time.UTC
	}
	return &Config{
		TimeZone: zone,
	}
}
