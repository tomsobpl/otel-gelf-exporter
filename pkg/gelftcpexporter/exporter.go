package gelftcpexporter

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

type gelfTcpExporter struct {
	config                    *Config
	logger                    *zap.Logger
	messageFactory            *ogcfactory.Factory
	writer                    *gelf.TCPWriter
	writerEndpoint            string
	writerEndpointRefreshTime int64
	writerTLSGateway          *TLSGateway
}

func newGelfTcpExporter(cfg component.Config, set exporter.Settings) *gelfTcpExporter {
	return &gelfTcpExporter{
		config:         cfg.(*Config),
		logger:         set.Logger,
		messageFactory: ogc.CreateFactory(set.Logger),
	}
}

func (e *gelfTcpExporter) initGelfWriter() bool {
	e.logger.Info(fmt.Sprintf("Initializing GELF writer for endpoint %s", e.config.Endpoint))

	var err = e.resolveWriterEndpoint()

	if err != nil {
		e.logger.Error(fmt.Sprintf("Failed to resolve IP address for %s", e.config.Endpoint))
		e.logger.Error(err.Error())
		return false
	}

	writerEndpoint := e.writerEndpoint

	if e.config.EndpointTLS.Enabled {
		e.logger.Info("Starting GELF TCP exporter TLS Proxy")

		srcEndpoint := TLSGatewayEndpoint{Network: "tcp", Endpoint: "127.0.0.1:"}
		dstEndpoint := TLSGatewayEndpoint{Network: "tcp", Endpoint: writerEndpoint}
		gateway, err := NewTLSGateway(srcEndpoint, dstEndpoint)

		if err != nil {
			e.logger.Error("Failed to start local listener")
			e.logger.Error(err.Error())
			return false
		}

		writerEndpoint = gateway.Addr().String()
		e.logger.Debug(fmt.Sprintf("Started local listener on %s", writerEndpoint))

		//var wg sync.WaitGroup
		//
		//wg.Add(1)
		//go gateway.Run(&wg)
		gateway.Start()

		//e.logger.Debug("Waiting for TLSGateway to be ready")
		//wg.Wait()

		if e.writerTLSGateway != nil {
			e.logger.Debug("Shutting down previous TLSGateway")
			e.writerTLSGateway.Shutdown()
		}

		e.writerTLSGateway = gateway
	}

	e.writer, err = gelf.NewTCPWriter(writerEndpoint)

	if err != nil {
		e.logger.Error(fmt.Sprintf("Failed to initialize GELF writer for endpoint %s", e.config.Endpoint))
		e.logger.Error(err.Error())
		return false
	}

	return e.writer != nil
}

func (e *gelfTcpExporter) start(_ context.Context, _ component.Host) error {
	e.logger.Info("Starting GELF TCP exporter")

	if !e.initGelfWriter() {
		e.logger.Error("Failed to initialize GELF writer")
	}

	return nil
}

func (e *gelfTcpExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	e.logger.Info(fmt.Sprintf("Processing %d resource log(s) with %d log record(s)", ld.ResourceLogs().Len(), ld.LogRecordCount()))

	if e.config.EndpointRefreshStrategy == gelfexporter.EndpointRefreshStrategyInterval && e.endpointRefreshIntervalExpired() {
		e.logger.Debug(fmt.Sprintf("Refreshing writer endpoint due to '%s' strategy", e.config.EndpointRefreshStrategy))
		e.initGelfWriter()
	}

	for _, m := range e.messageFactory.FromOtelLogsData(ld) {
		if e.config.EndpointRefreshStrategy == gelfexporter.EndpointRefreshStrategyPerMessage {
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

	e.logger.Debug(fmt.Sprintf("Resolved Endpoint %s into %s", e.config.Endpoint, e.writerEndpoint))

	return nil
}
