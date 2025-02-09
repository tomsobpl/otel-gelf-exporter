package gelfudpexporter

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"time"
)

type gelfUdpExporter struct {
	config *Config
	logger *zap.Logger
	writer *gelf.UDPWriter
}

func newGelfUdpExporter(cfg component.Config, set exporter.Settings) *gelfUdpExporter {
	c := cfg.(*Config)

	return &gelfUdpExporter{
		config: c,
		logger: set.Logger,
	}
}

func (e *gelfUdpExporter) initGelfWriter() bool {
	e.logger.Info(fmt.Sprintf("Initializing GELF writer for endpoint %s", e.config.Endpoint))
	w, _ := gelf.NewUDPWriter(e.config.Endpoint)
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

	m := &gelf.Message{
		Version:  "1.1",
		Host:     "hostname",
		Short:    "Some kind of error",
		Full:     "Optional full error message",
		TimeUnix: float64(time.Now().UnixNano()) / float64(time.Second),
		Level:    1,
		Facility: "",
		Extra:    nil,
		RawExtra: nil,
	}

	err := e.writer.WriteMessage(m)

	if err != nil {
		e.logger.Error(err.Error())
	}

	return nil
}
