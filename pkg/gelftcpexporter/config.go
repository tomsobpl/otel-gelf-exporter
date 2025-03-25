package gelftcpexporter

import (
	"github.com/tomsobpl/otel-gelf-exporter/pkg/gelfexporter"
	"go.opentelemetry.io/collector/component"
)

const (
	DefaultEndpointTLSEnabled            = true
	DefaultEndpointTLSInsecureSkipVerify = false
)

type Config struct {
	gelfexporter.Config `mapstructure:",squash"`

	// EndpointTLS is a configuration of the TLS connection.
	EndpointTLS EndpointTLS `mapstructure:"endpoint_tls"`
}

type EndpointTLS struct {
	// Enabled is a flag that enables or disables TLS.
	// Default is true.
	Enabled bool `mapstructure:"enabled"`

	// InsecureSkipVerify is a flag that determines whether to skip verification of the server's certificate chain and host name.
	// Default is false.
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
}

func CreateDefaultConfig() component.Config {
	return &Config{
		Config: *gelfexporter.CreateDefaultConfig().(*gelfexporter.Config),
		EndpointTLS: EndpointTLS{
			Enabled:            DefaultEndpointTLSEnabled,
			InsecureSkipVerify: DefaultEndpointTLSInsecureSkipVerify,
		},
	}
}
