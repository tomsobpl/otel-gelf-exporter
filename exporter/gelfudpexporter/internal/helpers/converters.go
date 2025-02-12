package helpers

// OtelSeverityToSyslogLevel maps OpenTelemetry severity number to Syslog level.
func OtelSeverityToSyslogLevel(severityNumber int32) int32 {
	/*
		OpenTelemetry severity ranges:
		1-4 	TRACE	A fine-grained debugging event. Typically disabled in default configurations.
		5-8		DEBUG	A debugging event.
		9-12	INFO	An informational event. Indicates that an event happened.
		13-16	WARN	A warning event. Not an error but is likely more important than an informational event.
		17-20	ERROR	An error event. Something went wrong.
		21-24	FATAL	A fatal error such as application or system crash.

		Syslog levels:
		0 EMERGENCY	System is unusable.
		1 ALERT		Action must be taken immediately.
		2 CRITICAL	Critical conditions.
		3 ERROR		Error conditions.
		4 WARNING	Warning conditions.
		5 NOTICE	Normal but significant condition.
		6 INFO		Informational messages.
		7 DEBUG		Debug-level messages.
	*/

	if severityNumber < 1 || severityNumber > 24 {
		panic("severity number out of range")
	}

	switch {
	case severityNumber >= 21:
		return 0 // EMERGENCY
	case severityNumber >= 17:
		return 3 // ERROR
	case severityNumber >= 13:
		return 4 // WARNING
	case severityNumber >= 9:
		return 6 // INFO
	}

	return 7 // DEBUG
}
