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
	// Endpoint is the address of the GELF input.
	Endpoint string `mapstructure:"endpoint"`

	// EndpointRefreshInterval is the interval in seconds between endpoint refreshes.
	// Default value is 60.
	// This is only used when EndpointRefreshStrategy is set to "interval".
	EndpointRefreshInterval int64 `mapstructure:"endpoint_refresh_interval"`

	// EndpointRefreshStrategy is the strategy used to refresh the endpoint.
	// Possible values are "none", "interval" and "perChunk".
	// Default value is "none".
	// "none" means that the endpoint is not refreshed.
	// "interval" means that the endpoint is refreshed every EndpointRefreshInterval seconds.
	// "perChunk" means that the endpoint is refreshed for every chunk of logs.
	EndpointRefreshStrategy string `mapstructure:"endpoint_refresh_strategy"`
}

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
