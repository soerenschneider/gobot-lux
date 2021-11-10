[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/gobot-lux)](https://goreportcard.com/report/github.com/soerenschneider/gobot-lux)

# Configuration
## Via Env Variables
| ENV                         | Default              | Description                                    |
|-----------------------------|----------------------|------------------------------------------------|
| GOBOT_LUX_PLACEMENT         | -                    | Location short name of this bot                |
| GOBOT_LUX_INTERVAL_S        | 30                   | Interval in seconds to dispatch message        |
| GOBOT_LUX_FIRMATA_PORT      | /dev/ttyUSB0         | Firmata port to use                            |
| GOBOT_LUX_AIO_PIN           | 0                    | AIO pin to use                                 |
| GOBOT_LUX_AIO_POLLING_INTERVAL_MS   | 500          | Interval in milliseconds to read from sensor   |
| GOBOT_LUX_MQTT_HOST         | -                    | Host of the MQTT broker, can be omitted        |
| GOBOT_LUX_MQTT_TOPIC        |                      | Topic to publish messages into                 |
| GOBOT_LUX_LOG_SENSOR        | false                | Log read sensor values                         |
| GOBOT_LUX_METRICS_ADDR      | :9194                | Prometheus http handler listen address         |

## Via Config File

```json
{
  "placement": "loc",
  "mqtt_host": "tcp://host:1883",
  "mqtt_topic": "sensors/%s/sub",
  "metrics_addr": ":1111",
  "log_sensor": true,
  "interval_s": 45,
  "firmata_port": "/dev/my-device",
  "aio_pin": "42",
  "aio_polling_interval_ms": 25
}
```

# Metrics

This project exposes the following metrics in Open Metrics format.

| Namespace | Subsystem | Name                           | Type    | Labels   | Help                                                               |
|-----------|-----------|--------------------------------|---------|----------|--------------------------------------------------------------------|
| gobot_lux |           | version                        | gauge   | version, commit | Version information of this robot                           |
| gobot_lux |           | heartbeat_seconds              | gauge   | placement | Continuous heartbeat of this bot                                  |
| gobot_lux | sensor    | brightness_level               | gauge   | placement | Current sensor reading of brightness level                        |
| gobot_lux | sensor    | read_errors_total              | gauge   | placement | The measured altitude in meters                                   |
| gobot_lux | mqtt      | messages_published_total       | counter | placement | The amount of published MQTT messages                             |
| gobot_lux | mqtt      | message_publish_errors_total   | counter | placement | Total amount of errors while trying to publish messages over MQTT |
