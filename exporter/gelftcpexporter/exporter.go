package gelftcpexporter

import (
	"context"
	"fmt"
	"github.com/tomsobpl/otel-collector-graylog/component/gelfmessage"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

type gelfTcpExporter struct {
	config         *Config
	logger         *zap.Logger
	messageFactory *gelfmessage.Factory
	writer         *gelf.TCPWriter
}

func newGelfTcpExporter(cfg component.Config, set exporter.Settings) *gelfTcpExporter {
	c := cfg.(*Config)

	return &gelfTcpExporter{
		config:         c,
		logger:         set.Logger,
		messageFactory: gelfmessage.NewFactory(set.Logger),
	}
}

func (e *gelfTcpExporter) initGelfWriter() bool {
	e.logger.Info(fmt.Sprintf("Initializing GELF writer for endpoint %s", e.config.Endpoint))
	w, _ := gelf.NewTCPWriter(e.config.Endpoint)
	e.writer = w

	return e.writer != nil
}

func (e *gelfTcpExporter) start(ctx context.Context, host component.Host) error {
	e.logger.Info("Starting Graylog exporter")

	if !e.initGelfWriter() {
		e.logger.Error("Failed to initialize GELF writer")
	}

	return nil
}

func (e *gelfTcpExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	e.logger.Info(fmt.Sprintf("Processing %d resource log(s) with %d log record(s)", ld.ResourceLogs(), ld.LogRecordCount()))

	for _, m := range e.messageFactory.BuildMessagesFromOtelLogsData(ld) {
		//@TODO: Target connection should be checked before writing
		//@TODO: Target should be refreshed if setup in config to allow for load balancing usage

		err := e.writer.WriteMessage(m.GetRawMessage())

		if err != nil {
			e.logger.Error(err.Error())
		}
	}

	return nil
}
