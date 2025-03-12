package gelfmessage

import (
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type Factory struct {
	logger *zap.Logger
}

func NewFactory(logger *zap.Logger) *Factory {
	return &Factory{
		logger: logger,
	}
}

func (f *Factory) BuildMessagesFromOtelLogsData(logs plog.Logs) []*Message {
	messages := make([]*Message, 0)

	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		messages = append(messages, f.handleOtelResourceLogs(logs.ResourceLogs().At(i))...)
	}

	return messages
}

func (f *Factory) handleOtelResourceLogs(rl plog.ResourceLogs) []*Message {
	messages := make([]*Message, 0)

	for i := 0; i < rl.ScopeLogs().Len(); i++ {
		messages = append(messages, f.handleOtelScopeLogs(rl.ScopeLogs().At(i))...)
	}

	for _, m := range messages {
		m.UpdateExtraFields(map[string]interface{}{
			"otel_resource_dropped_attributes_count": rl.Resource().DroppedAttributesCount(),
		})

		m.UpdateExtraFields(OtelAttributesToGelfExtraWithPrefix(rl.Resource().Attributes(), "resource"))

		host, hostExist := rl.Resource().Attributes().Get("host.name")

		if hostExist {
			m.SetHost(host.AsString())
		}
	}

	return messages
}

func (f *Factory) handleOtelScopeLogs(sl plog.ScopeLogs) []*Message {
	messages := make([]*Message, 0)

	for i := 0; i < sl.LogRecords().Len(); i++ {
		messages = append(messages, f.handleOtelLogRecord(sl.LogRecords().At(i)))
	}

	for _, m := range messages {
		m.UpdateExtraFields(map[string]interface{}{
			"otel_scope_dropped_attributes_count": sl.Scope().DroppedAttributesCount(),
			"otel_scope_name":                     sl.Scope().Name(),
			"otel_scope_version":                  sl.Scope().Version(),
		})

		m.UpdateExtraFields(OtelAttributesToGelfExtraWithPrefix(sl.Scope().Attributes(), "scope"))
	}

	return messages
}

func (f *Factory) handleOtelLogRecord(lr plog.LogRecord) *Message {
	return BuildFromOtelLogRecord(lr)
}
