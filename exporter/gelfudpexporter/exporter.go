package gelfudpexporter

import (
	"context"
	"fmt"
	"github.com/tomsobpl/otel-collector-graylog/exporter/gelfudpexporter/internal/helpers"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
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

func (e *gelfUdpExporter) handleLogRecord(lr plog.LogRecord) *gelf.Message {
	//@TODO: Full message implementation

	m := &gelf.Message{
		Version:  "1.1",
		Host:     "UNKNOWN",
		Short:    lr.Body().AsString(),
		Full:     "TODO",
		TimeUnix: helpers.OtelTimestampToGelfTimeUnix(lr.Timestamp(), lr.ObservedTimestamp()),
		Level:    helpers.OtelSeverityToSyslogLevel(int32(lr.SeverityNumber())),
		Facility: "",
		Extra: map[string]interface{}{
			"otel_log_dropped_attributes_count": lr.DroppedAttributesCount(),
			"otel_log_event_name":               lr.EventName(),
			"otel_log_severity_number":          lr.SeverityNumber(),
			"otel_log_severity_text":            lr.SeverityText(),
			"otel_log_span_id":                  lr.SpanID().String(),
			"otel_log_trace_id":                 lr.TraceID().String(),
		},
		RawExtra: nil,
	}

	helpers.OtelAttributesToGelfExtra(lr.Attributes(), m)

	return m
}

func (e *gelfUdpExporter) handleScopeLog(sl plog.ScopeLogs) []*gelf.Message {
	msgs := make([]*gelf.Message, 0)

	for i := 0; i < sl.LogRecords().Len(); i++ {
		msgs = append(msgs, e.handleLogRecord(sl.LogRecords().At(i)))
	}

	for _, m := range msgs {
		m.Extra["otel_scope_dropped_attributes_count"] = sl.Scope().DroppedAttributesCount()
		m.Extra["otel_scope_name"] = sl.Scope().Name()
		m.Extra["otel_scope_version"] = sl.Scope().Version()
		helpers.OtelAttributesToGelfExtraWithPrefix(sl.Scope().Attributes(), m, "scope")
	}

	return msgs
}

func (e *gelfUdpExporter) handleResourceLog(rl plog.ResourceLogs) []*gelf.Message {
	msgs := make([]*gelf.Message, 0)

	for i := 0; i < rl.ScopeLogs().Len(); i++ {
		msgs = append(msgs, e.handleScopeLog(rl.ScopeLogs().At(i))...)
	}

	for _, m := range msgs {
		m.Extra["otel_resource_dropped_attributes_count"] = rl.Resource().DroppedAttributesCount()
		helpers.OtelAttributesToGelfExtraWithPrefix(rl.Resource().Attributes(), m, "resource")

		host, hostExist := rl.Resource().Attributes().Get("host.name")

		if hostExist {
			m.Host = host.AsString()
		}
	}

	return msgs
}

func (e *gelfUdpExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	e.logger.Info(fmt.Sprintf("Processing %d resource log(s) with %d log record(s)", ld.ResourceLogs(), ld.LogRecordCount()))

	msgs := make([]*gelf.Message, 0)

	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		msgs = append(msgs, e.handleResourceLog(ld.ResourceLogs().At(i))...)
	}

	for _, m := range msgs {
		//@TODO: Target should be refreshed if setup in config to allow for load balancing usage

		err := e.writer.WriteMessage(m)

		if err != nil {
			e.logger.Error(err.Error())
		}
	}

	return nil
}
