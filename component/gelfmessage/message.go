package gelfmessage

import (
	"go.opentelemetry.io/collector/pdata/plog"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

type Message struct {
	rawmsg *gelf.Message
}

func (m *Message) GetRawMessage() *gelf.Message {
	return m.rawmsg
}

func (m *Message) SetBody(body string) {
	m.rawmsg.Short = body
}

func (m *Message) SetHost(host string) {
	m.rawmsg.Host = host
}

func (m *Message) SetSeverity(severity int32) {
	m.rawmsg.Level = severity
}

func (m *Message) SetTimestamp(timestamp float64) {
	m.rawmsg.TimeUnix = timestamp
}

func (m *Message) UpdateExtraFields(fields map[string]interface{}) {
	for k, v := range fields {
		m.rawmsg.Extra[k] = v
	}
}

func BuildFromOtelLogRecord(lr plog.LogRecord) *Message {
	m := createEmptyMessage()

	m.SetBody(lr.Body().AsString())
	m.SetSeverity(OtelSeverityToSyslogLevel(int32(lr.SeverityNumber())))
	m.SetTimestamp(OtelTimestampToGelfTimeUnix(lr.Timestamp(), lr.ObservedTimestamp()))

	m.UpdateExtraFields(map[string]interface{}{
		"otel_log_dropped_attributes_count": lr.DroppedAttributesCount(),
		"otel_log_event_name":               lr.EventName(),
		"otel_log_severity_number":          lr.SeverityNumber(),
		"otel_log_severity_text":            lr.SeverityText(),
		"otel_log_span_id":                  lr.SpanID().String(),
		"otel_log_trace_id":                 lr.TraceID().String(),
	})

	m.UpdateExtraFields(OtelAttributesToGelfExtra(lr.Attributes()))

	return m
}

func createEmptyMessage() *Message {
	return &Message{rawmsg: &gelf.Message{
		Version:  "1.1",
		Host:     "UNKNOWN",
		Extra:    map[string]interface{}{},
		RawExtra: nil,
	}}
}
