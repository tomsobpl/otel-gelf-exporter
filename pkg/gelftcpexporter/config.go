package gelftcpexporter

import (
	"errors"
	"go.opentelemetry.io/collector/component"
)

type Config struct {
	Endpoint string `mapstructure:"endpoint"`
}

func (cfg *Config) Validate() error {
	if cfg.Endpoint == "" {
		return errors.New("graylog UDP endpoint must be specified")
	}

	return nil
}

func createDefaultConfig() component.Config {
	return &Config{}
}
