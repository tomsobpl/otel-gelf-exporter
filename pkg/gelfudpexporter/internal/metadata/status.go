package metadata

import (
	"go.opentelemetry.io/collector/component"
)

var (
	Type      = component.MustNewType("gelfudp")
	ScopeName = "github.com/tomsobpl/otel-gelf-exporter/pkg/gelfudpexporter"
)

const (
	ExporterStabilityLevel = component.StabilityLevelDevelopment
)
