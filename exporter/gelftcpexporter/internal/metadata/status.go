package metadata

import (
	"go.opentelemetry.io/collector/component"
)

var (
	Type      = component.MustNewType("gelftcp")
	ScopeName = "github.com/tomsobpl/otel-collector-graylog/exporter/gelftcpexporter"
)

const (
	ExporterStabilityLevel = component.StabilityLevelDevelopment
)
