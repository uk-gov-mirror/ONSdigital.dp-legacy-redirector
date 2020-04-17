package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-legacy-redirector
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	HealthckeckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	HealthckeckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
}

var cfg *Config

// Get returns the default config with any modifications made through environment variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg := &Config{
		BindAddr:                   ":24600",
		HealthckeckCriticalTimeout: time.Minute,
		HealthckeckInterval:        time.Second * 10,
	}

	return cfg, envconfig.Process("", cfg)
}
