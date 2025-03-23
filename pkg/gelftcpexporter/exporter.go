package gelftcpexporter

import (
	"context"
	"crypto/tls"
	"fmt"
	ogc "github.com/tomsobpl/otel-gelf-converter/pkg"
	ogcfactory "github.com/tomsobpl/otel-gelf-converter/pkg/factory"
	"github.com/tomsobpl/otel-gelf-exporter/pkg/gelfexporter"
	"github.com/tomsobpl/otel-gelf-exporter/pkg/gelftcpexporter/internal/tlsgateway"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"time"
)

type gelfTcpExporter struct {
	config                    *Config
	logger                    *zap.Logger
	messageFactory            *ogcfactory.Factory
	writer                    *gelf.TCPWriter
	writerEndpoint            string
	writerEndpointRefreshTime int64
	writerTLSGateway          *tlsgateway.TLSGateway
}

func newGelfTcpExporter(cfg component.Config, set exporter.Settings) *gelfTcpExporter {
	return &gelfTcpExporter{
		config:         cfg.(*Config),
		logger:         set.Logger,
		messageFactory: ogc.CreateFactory(set.Logger),
	}
}

func (e *gelfTcpExporter) initGelfWriter() bool {
	e.logger.Info(fmt.Sprintf("initializing GELF writer for endpoint %s", e.config.Endpoint))

	var err = e.resolveWriterEndpoint()

	if err != nil {
		e.logger.Error(fmt.Sprintf("failed to resolve IP address for %s", e.config.Endpoint), zap.Error(err))
		return false
	}

	writerEndpoint := e.writerEndpoint

	if e.config.EndpointTLS.Enabled {
		e.logger.Info("starting GELF TCP exporter TLS Proxy")

		srcEndpoint := tlsgateway.Endpoint{Network: "tcp", Endpoint: "127.0.0.1:"}
		dstEndpoint := tlsgateway.Endpoint{Network: "tcp", Endpoint: writerEndpoint}
		gateway, err := tlsgateway.NewTLSGateway(srcEndpoint, dstEndpoint, e.logger)

		if err != nil {
			e.logger.Error("failed to start local listener", zap.Error(err))
			return false
		}

		writerEndpoint = gateway.Addr().String()
		e.logger.Debug(fmt.Sprintf("started local listener on %s", writerEndpoint))

		tlsConfig := &tls.Config{
			InsecureSkipVerify: e.config.EndpointTLS.InsecureSkipVerify,
		}

		if err := gateway.Start(tlsConfig); err != nil {
			e.logger.Error("failed to start TLS gateway", zap.Error(err))
			return false
		}

		if e.writerTLSGateway != nil {
			e.logger.Debug("shutting down previous TLSGateway")
			if err := e.writerTLSGateway.Shutdown(); err != nil {
				e.logger.Error("failed to shutdown previous TLSGateway", zap.Error(err))
			}
		}

		e.writerTLSGateway = gateway
	}

	e.writer, err = gelf.NewTCPWriter(writerEndpoint)

	if err != nil {
		e.logger.Error(fmt.Sprintf("failed to initialize GELF writer for endpoint %s", e.config.Endpoint), zap.Error(err))
		return false
	}

	return e.writer != nil
}

func (e *gelfTcpExporter) start(_ context.Context, _ component.Host) error {
	e.logger.Info("starting GELF TCP exporter")

	if !e.initGelfWriter() {
		e.logger.Error("failed to initialize GELF writer")
	}

	return nil
}

func (e *gelfTcpExporter) pushLogs(_ context.Context, ld plog.Logs) error {
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

func (e *gelfTcpExporter) endpointRefreshIntervalExpired() bool {
	return time.Now().Unix()-e.writerEndpointRefreshTime > e.config.EndpointRefreshInterval
}

func (e *gelfTcpExporter) resolveWriterEndpoint() error {
	endpoint, err := gelfexporter.ResolveEndpoint(e.config.Endpoint)

	if err != nil {
		return err
	}

	e.writerEndpoint = endpoint
	e.writerEndpointRefreshTime = time.Now().Unix()

	e.logger.Debug(fmt.Sprintf("resolved Endpoint %s into %s", e.config.Endpoint, e.writerEndpoint))

	return nil
}
