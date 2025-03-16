package gelfudpexporter

import (
	"errors"
	"go.opentelemetry.io/collector/component"
)

const (
	EndpointRefreshStrategyDefault  string = "none"
	EndpointRefreshStrategyInterval string = "interval"
	EndpointRefreshStrategyPerchunk string = "perChunk"
)

type Config struct {
	Endpoint                string `mapstructure:"endpoint"`
	EndpointRefreshInterval int64  `mapstructure:"endpoint_refresh_interval"`
	EndpointRefreshStrategy string `mapstructure:"endpoint_refresh_strategy"`
}

//var _ component.Config = (*Config)(nil)

func (cfg *Config) Validate() error {
	if cfg.Endpoint == "" {
		return errors.New("graylog UDP endpoint must be specified")
	}

	switch cfg.EndpointRefreshStrategy {
	case EndpointRefreshStrategyDefault, EndpointRefreshStrategyInterval, EndpointRefreshStrategyPerchunk:
		break
	default:
		return errors.New("invalid endpoint refresh strategy")
	}

	return nil
}

func createDefaultConfig() component.Config {
	return &Config{
		EndpointRefreshInterval: 60,
		EndpointRefreshStrategy: EndpointRefreshStrategyDefault,
	}
}
