package gelfexporter

import (
	"errors"
	"go.opentelemetry.io/collector/component"
)

const (
	DefaultEndpointInitBackoff        int    = 10
	DefaultEndpointInitRetries        int    = 5
	DefaultEndpointRefreshInterval    int64  = 60
	EndpointRefreshStrategyNone       string = "none"
	EndpointRefreshStrategyInterval   string = "interval"
	EndpointRefreshStrategyPerMessage string = "perMessage"
	TcpExporterType                   string = "gelftcp"
	UdpExporterType                   string = "gelfudp"
)

type Config struct {
	// Endpoint is the address of the GELF input.
	Endpoint string `mapstructure:"endpoint"`

	// EndpointInitBackoff is a delay between retries to initialize the endpoint.
	// Default is 10.
	EndpointInitBackoff int `mapstructure:"endpoint_init_backoff"`

	// EndpointInitRetries is a number of retries to initialize the endpoint.
	// Default is 5.
	EndpointInitRetries int `mapstructure:"endpoint_init_retries"`

	// EndpointRefreshInterval is the interval in seconds between endpoint refreshes.
	// Default value is 60.
	// This is only used when EndpointRefreshStrategy is set to "interval".
	EndpointRefreshInterval int64 `mapstructure:"endpoint_refresh_interval"`

	// EndpointRefreshStrategy is the strategy used to refresh the endpoint.
	// Possible values are "none", "interval" and "perMessage".
	// Default value is "none".
	// "none" means that the endpoint is not refreshed.
	// "interval" means that the endpoint is refreshed every EndpointRefreshInterval seconds.
	// "perMessage" means that the endpoint is refreshed for every log message.
	EndpointRefreshStrategy string `mapstructure:"endpoint_refresh_strategy"`
}

func (cfg *Config) Validate() error {
	if cfg.Endpoint == "" {
		return errors.New("GELF input endpoint must be specified")
	}

	switch cfg.EndpointRefreshStrategy {
	case EndpointRefreshStrategyNone, EndpointRefreshStrategyInterval, EndpointRefreshStrategyPerMessage:
		break
	default:
		return errors.New("invalid endpoint refresh strategy")
	}

	return nil
}

// CreateDefaultConfig creates the default configuration for the exporter.
func CreateDefaultConfig() component.Config {
	return &Config{
		EndpointInitBackoff:     DefaultEndpointInitBackoff,
		EndpointInitRetries:     DefaultEndpointInitRetries,
		EndpointRefreshInterval: DefaultEndpointRefreshInterval,
		EndpointRefreshStrategy: EndpointRefreshStrategyNone,
	}
}
