package gelfudpexporter

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

type gelfUdpExporter struct {
	logger *zap.Logger
	writer *gelf.UDPWriter
}

func newGelfUdpExporter(set exporter.Settings) *gelfUdpExporter {
	return &gelfUdpExporter{
		logger: set.Logger,
	}
}

func (e *gelfUdpExporter) initGelfWriter() bool {
	e.logger.Info("Initializing GELF writer")
	w, _ := gelf.NewUDPWriter("localhost:12201")
	e.writer = w

	return e.writer != nil
}

func (e *gelfUdpExporter) start(ctx context.Context, host component.Host) error {
	e.logger.Info("Starting Graylog exporter")

	if !e.initGelfWriter() {
		e.logger.Error("Failed to initialize GELF writer")
	}

	return nil
}

func (e *gelfUdpExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	e.logger.Info("Logs",
		zap.Int("resource logs", ld.ResourceLogs().Len()),
		zap.Int("log records", ld.LogRecordCount()))

	return nil
}
