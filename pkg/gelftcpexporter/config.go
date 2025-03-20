package gelftcpexporter

import "github.com/tomsobpl/otel-gelf-exporter/pkg/gelfexporter"

type Config struct {
	gelfexporter.Config
	EndpointTLS *EndpointTLS `mapstructure:"endpoint_tls"`
}

type EndpointTLS struct {
	// Enabled is a flag that enables or disables TLS.
	// Default is true.
	Enabled bool `mapstructure:"enabled"`
}

func CreateDefaultConfig() *Config {
	return &Config{
		Config: *gelfexporter.CreateDefaultConfig().(*gelfexporter.Config),
		EndpointTLS: &EndpointTLS{
			Enabled: true,
		},
	}
}
