package gelftcpexporter

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomsobpl/otel-gelf-exporter/pkg/gelfexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/confmap/xconfmap"
	"testing"
)

func TestConfigLoading(t *testing.T) {
	cm, err := confmaptest.LoadConf("testdata/config.yaml")
	require.NoError(t, err)

	tests := []struct {
		id       component.ID
		expected component.Config
	}{
		{
			id: component.NewIDWithName(component.MustNewType(gelfexporter.TcpExporterType), ""),
			expected: &Config{
				Config: gelfexporter.Config{
					Endpoint:                "localhost:12201",
					EndpointRefreshInterval: 60,
					EndpointRefreshStrategy: gelfexporter.EndpointRefreshStrategyDefault,
				},
				EndpointTLS: EndpointTLS{
					Enabled:            true,
					InsecureSkipVerify: false,
				},
			},
		},
		{
			id: component.NewIDWithName(component.MustNewType(gelfexporter.TcpExporterType), "insecure"),
			expected: &Config{
				Config: gelfexporter.Config{
					Endpoint:                "localhost:12201",
					EndpointRefreshInterval: 60,
					EndpointRefreshStrategy: gelfexporter.EndpointRefreshStrategyDefault,
				},
				EndpointTLS: EndpointTLS{
					Enabled:            false,
					InsecureSkipVerify: false,
				},
			},
		},
		{
			id: component.NewIDWithName(component.MustNewType(gelfexporter.TcpExporterType), "skipverify"),
			expected: &Config{
				Config: gelfexporter.Config{
					Endpoint:                "localhost:12201",
					EndpointRefreshInterval: 60,
					EndpointRefreshStrategy: gelfexporter.EndpointRefreshStrategyDefault,
				},
				EndpointTLS: EndpointTLS{
					Enabled:            true,
					InsecureSkipVerify: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			cfg := CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, sub.Unmarshal(cfg))

			assert.NoError(t, xconfmap.Validate(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr string
	}{
		{
			name: "NoEndpoint",
			cfg: func() *Config {
				cfg := CreateDefaultConfig().(*Config)
				return cfg
			}(),
			wantErr: "GELF input endpoint must be specified",
		},
		{
			name: "InvalidEndpointRefreshStrategy",
			cfg: func() *Config {
				cfg := CreateDefaultConfig().(*Config)
				cfg.Endpoint = "localhost:12201"
				cfg.EndpointRefreshStrategy = "invalid"
				return cfg
			}(),
			wantErr: "invalid endpoint refresh strategy",
		},
		{
			name: "Success",
			cfg: func() *Config {
				cfg := CreateDefaultConfig().(*Config)
				cfg.Endpoint = "localhost:12201"
				return cfg
			}(),
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
