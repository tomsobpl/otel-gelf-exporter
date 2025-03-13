package metadata

import (
	"go.opentelemetry.io/collector/component"
)

var (
	Type      = component.MustNewType("gelfudp")
	ScopeName = "github.com/tomsobpl/otel-collector-graylog/exporter/gelfudpexporter"
)

const (
	ExporterStabilityLevel = component.StabilityLevelDevelopment
)
