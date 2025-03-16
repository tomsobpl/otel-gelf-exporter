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
	"strings"
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
	return &gelfUdpExporter{
		config:         cfg.(*Config),
		logger:         set.Logger,
		messageFactory: ogc.CreateFactory(set.Logger),
	}
}

func (e *gelfUdpExporter) initGelfWriter() bool {
	e.logger.Info(fmt.Sprintf("Initializing GELF writer for endpoint %s", e.config.Endpoint))

	var err = e.resolveWriterEndpoint()

	if err != nil {
		e.logger.Error(fmt.Sprintf("Failed to resolve IP address for %s", e.config.Endpoint))
		e.logger.Error(err.Error())
		return false
	}

	e.writer, err = gelf.NewUDPWriter(e.writerEndpoint)

	if err != nil {
		e.logger.Error(fmt.Sprintf("Failed to initialize GELF writer for endpoint %s", e.config.Endpoint))
		e.logger.Error(err.Error())
		return false
	}

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
	e.logger.Info(fmt.Sprintf("Processing %d resource log(s) with %d log record(s)", ld.ResourceLogs().Len(), ld.LogRecordCount()))

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
	host := e.config.Endpoint
	port := ""

	if strings.LastIndexByte(e.config.Endpoint, ':') != -1 {
		h, p, err := net.SplitHostPort(e.config.Endpoint)

		if err != nil {
			return err
		}

		host = h
		port = p
	}

	ips, err := net.LookupIP(host)

	if err != nil || ips == nil || len(ips) == 0 {
		return err
	}

	if port != "" {
		e.writerEndpoint = net.JoinHostPort(ips[0].String(), port)
	} else {
		e.writerEndpoint = ips[0].String()
	}

	e.writerEndpointRefreshTime = time.Now().Unix()
	e.logger.Debug(fmt.Sprintf("Resolved Endpoint %s into %s", e.config.Endpoint, e.writerEndpoint))

	return nil
}
