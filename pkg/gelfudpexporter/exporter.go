package gelfudpexporter

import (
	"context"
	"fmt"
	ogc "github.com/tomsobpl/otel-gelf-converter/pkg"
	ogcfactory "github.com/tomsobpl/otel-gelf-converter/pkg/factory"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"net"
	"time"
)

type gelfUdpExporter struct {
	config                    *Config
	logger                    *zap.Logger
	messageFactory            *ogcfactory.Factory
	writer                    *gelf.UDPWriter
	writerEndpoint            string
	writerEndpointRefreshTime int64
}

func newGelfUdpExporter(cfg component.Config, set exporter.Settings) *gelfUdpExporter {
	c := cfg.(*Config)

	return &gelfUdpExporter{
		config:         c,
		logger:         set.Logger,
		messageFactory: ogc.CreateFactory(set.Logger),
	}
}

func (e *gelfUdpExporter) initGelfWriter() bool {
	e.logger.Info(fmt.Sprintf("Initializing GELF writer for endpoint %s", e.config.Endpoint))

	if err := e.resolveWriterEndpoint(); err != nil {
		e.logger.Error(err.Error())
		return false
	}

	w, _ := gelf.NewUDPWriter(e.writerEndpoint)
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
	e.logger.Info(fmt.Sprintf("Processing %d resource log(s) with %d log record(s)", ld.ResourceLogs(), ld.LogRecordCount()))

	if e.config.EndpointRefreshStrategy == EndpointRefreshStrategyInterval && e.endpointRefreshIntervalExpired() {
		e.logger.Debug(fmt.Sprintf("Refreshing writer endpoint due to '%s' strategy", e.config.EndpointRefreshStrategy))
		e.initGelfWriter()
	}

	for _, m := range e.messageFactory.FromOtelLogsData(ld) {
		if e.config.EndpointRefreshStrategy == EndpointRefreshStrategyPerchunk {
			e.logger.Debug(fmt.Sprintf("Refreshing writer endpoint due to '%s' strategy", e.config.EndpointRefreshStrategy))
			e.initGelfWriter()
		}

		err := e.writer.WriteMessage(m.GetRawMessage())

		if err != nil {
			e.logger.Error(err.Error())
		}
	}

	return nil
}

func (e *gelfUdpExporter) endpointRefreshIntervalExpired() bool {
	return time.Now().Unix()-e.writerEndpointRefreshTime > e.config.EndpointRefreshInterval
}

func (e *gelfUdpExporter) resolveWriterEndpoint() error {
	ips, err := net.LookupIP(e.config.Endpoint)

	if err != nil {
		e.logger.Error(err.Error())
	}

	e.writerEndpoint = ips[0].String()
	e.writerEndpointRefreshTime = time.Now().Unix()

	e.logger.Debug(fmt.Sprintf("Resolved Endpoint %s into %s", e.config.Endpoint, e.writerEndpoint))

	return nil
}
