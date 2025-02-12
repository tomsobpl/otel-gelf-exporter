import logging
import time
import signal

from opentelemetry._logs import set_logger_provider
from opentelemetry.exporter.otlp.proto.grpc._log_exporter import OTLPLogExporter
from opentelemetry.sdk._logs import LoggerProvider, LoggingHandler
from opentelemetry.sdk._logs.export import BatchLogRecordProcessor
from opentelemetry.sdk.resources import Resource

# Create and set the logger provider
logger_provider = LoggerProvider()
set_logger_provider(logger_provider)

# Create the OTLP log exporter that sends logs to configured destination
exporter = OTLPLogExporter()
logger_provider.add_log_record_processor(BatchLogRecordProcessor(exporter))

# Attach OTLP handler to root logger
handler = LoggingHandler(logger_provider=logger_provider)
logging.getLogger().addHandler(handler)

logger = logging.getLogger(__name__)

keep_going = True

def shutdown(signum, frame):
    global keep_going

    logging.info("Shutting down")
    keep_going = False

def main():
    global keep_going

    logger.info("Starting the application")

    log_id = 0

    while keep_going:
        logger.warning(f"Log incident #{log_id}")

        log_id += 1
        time.sleep(3)

    logger_provider.shutdown()

if __name__ == "__main__":
    signal.signal(signal.SIGINT, shutdown)
    main()
    logger.info("Bye!")
