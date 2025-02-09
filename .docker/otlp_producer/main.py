import logging
import os
import random
import string
import time

from secrets import token_bytes

import requests

def otlp_log_sample():
    sample = {
        "resourceLogs": [
            {
                "resource": {
                    "attributes": [
                        {
                            "key": "service.name",
                            "value": { "stringValue": "my.service" }
                        }
                    ]
                },
                "scopeLogs": [
                    {
                        "scope": {
                            "name": "my.library",
                            "version": "1.0.0",
                            "attributes": [
                                {
                                    "key": "my.scope.attribute",
                                    "value": { "stringValue": "some scope attribute" }
                                }
                            ]
                        },
                        "logRecords": [
                            {
                                "timeUnixNano": time.time_ns(),
                                "observedTimeUnixNano": time.time_ns() - 10_000,
                                "severityNumber": random.randint(0, 9),
                                "severityText": "Information",
                                "traceId": token_bytes(16).hex(),
                                "spanId": token_bytes(8).hex(),
                                "body": {
                                    "stringValue": ''.join(random.choices(string.ascii_uppercase + string.digits, k=20))
                                },
                                "attributes": [
                                    {
                                        "key": "string.attribute",
                                        "value": { "stringValue": ''.join(random.choices(string.ascii_uppercase + string.digits, k=20)) }
                                    },
                                    {
                                        "key": "boolean.attribute",
                                        "value": { "boolValue": random.choice([True, False]) }
                                    },
                                    {
                                        "key": "int.attribute",
                                        "value": { "intValue": random.randint(0, 100) }
                                    },
                                    {
                                        "key": "double.attribute",
                                        "value": { "doubleValue": random.uniform(0, 100) }
                                    },
                                    {
                                        "key": "array.attribute",
                                        "value": {
                                            "arrayValue": {
                                                "values": [
                                                    { "stringValue": "many" }, { "stringValue": "values" }
                                                ]
                                            }
                                        }
                                    },
                                    {
                                        "key": "map.attribute",
                                        "value": {
                                            "kvlistValue": {
                                                "values": [
                                                    {
                                                        "key": "some.map.key",
                                                        "value": { "stringValue": "some value" }
                                                    }
                                                ]
                                            }
                                        }
                                    }
                                ]
                            }
                        ]
                    }
                ]
            }
        ]
    }

    return sample

def main():
    logger = logging.getLogger()

    while True:
        try:
            logger.info("Sending log sample ...")
            requests.post(os.environ.get('OTLP_HTTP_ENDPOINT'), json=otlp_log_sample())
        except Exception as e:
            logger.error(f"Failed to send log sample: {e}")

        time.sleep(3)

if __name__ == "__main__":
    main()
