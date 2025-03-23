package gelfudpexporter

import (
	"context"
	"fmt"
	ogc "github.com/tomsobpl/otel-gelf-converter/pkg"
	ogcfactory "github.com/tomsobpl/otel-gelf-converter/pkg/factory"
	"github.com/tomsobpl/otel-gelf-exporter/pkg/gelfexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"time"
)

type gelfUdpExporter struct {
	config                    *gelfexporter.Config
	logger                    *zap.Logger
	messageFactory            *ogcfactory.Factory
	writer                    *gelf.UDPWriter
	writerEndpoint            string
	writerEndpointRefreshTime int64
}

func newGelfUdpExporter(cfg component.Config, set exporter.Settings) *gelfUdpExporter {
	return &gelfUdpExporter{
		config:         cfg.(*gelfexporter.Config),
		logger:         set.Logger,
		messageFactory: ogc.CreateFactory(set.Logger),
	}
}

func (e *gelfUdpExporter) initGelfWriter() bool {
	e.logger.Info(fmt.Sprintf("initializing GELF writer for endpoint %s", e.config.Endpoint))
	var err error

	if err = e.resolveWriterEndpoint(); err != nil {
		e.logger.Error(fmt.Sprintf("failed to resolve IP address for %s", e.config.Endpoint), zap.Error(err))
		return false
	}

	if e.writer, err = gelf.NewUDPWriter(e.writerEndpoint); err != nil {
		e.logger.Error(fmt.Sprintf("failed to initialize GELF writer for endpoint %s", e.config.Endpoint), zap.Error(err))
		return false
	}

	return e.writer != nil
}

func (e *gelfUdpExporter) start(_ context.Context, _ component.Host) error {
	e.logger.Info("starting GELF UDP exporter")

	if !e.initGelfWriter() {
		e.logger.Error("failed to initialize GELF writer")
	}

	return nil
}

func (e *gelfUdpExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	e.logger.Info(fmt.Sprintf("processing %d resource log(s) with %d log record(s)", ld.ResourceLogs().Len(), ld.LogRecordCount()))

	if e.config.EndpointRefreshStrategy == gelfexporter.EndpointRefreshStrategyInterval && e.endpointRefreshIntervalExpired() {
		e.logger.Debug(fmt.Sprintf("refreshing writer endpoint due to '%s' strategy", e.config.EndpointRefreshStrategy))
		e.initGelfWriter()
	}

	for _, m := range e.messageFactory.FromOtelLogsData(ld) {
		if e.config.EndpointRefreshStrategy == gelfexporter.EndpointRefreshStrategyPerMessage {
			e.logger.Debug(fmt.Sprintf("refreshing writer endpoint due to '%s' strategy", e.config.EndpointRefreshStrategy))
			e.initGelfWriter()
		}

		if err := e.writer.WriteMessage(m.GetRawMessage()); err != nil {
			e.logger.Error("failed to write message", zap.Error(err))
		}
	}

	return nil
}

func (e *gelfUdpExporter) endpointRefreshIntervalExpired() bool {
	return time.Now().Unix()-e.writerEndpointRefreshTime > e.config.EndpointRefreshInterval
}

func (e *gelfUdpExporter) resolveWriterEndpoint() error {
	endpoint, err := gelfexporter.ResolveEndpoint(e.config.Endpoint)

	if err != nil {
		return err
	}

	e.writerEndpoint = endpoint
	e.writerEndpointRefreshTime = time.Now().Unix()

	e.logger.Debug(fmt.Sprintf("resolved Endpoint %s into %s", e.config.Endpoint, e.writerEndpoint))

	return nil
}
