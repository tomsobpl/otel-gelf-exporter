package gelftcpexporter

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/tomsobpl/otel-gelf-exporter/pkg/gelftcpexporter/internal/metadata"
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		CreateDefaultConfig,
		exporter.WithLogs(createLogsExporter, metadata.ExporterStabilityLevel),
	)
}

func createLogsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config) (exporter.Logs, error) {
	e := newGelfTcpExporter(cfg, set)

	return exporterhelper.NewLogs(ctx, set, cfg, e.pushLogs, exporterhelper.WithStart(e.start))
}
