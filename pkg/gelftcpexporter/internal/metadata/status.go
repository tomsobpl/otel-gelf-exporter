package metadata

import (
	"go.opentelemetry.io/collector/component"
)

var (
	Type      = component.MustNewType("gelftcp")
	ScopeName = "github.com/tomsobpl/otel-gelf-exporter/pkg/gelftcpexporter"
)

const (
	ExporterStabilityLevel = component.StabilityLevelDevelopment
)
