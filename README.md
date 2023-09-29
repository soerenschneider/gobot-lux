# gobot-lux
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/gobot-lux)](https://goreportcard.com/report/github.com/soerenschneider/gobot-lux)
![test-workflow](https://github.com/soerenschneider/gobot-lux/actions/workflows/test.yaml/badge.svg)
![release-workflow](https://github.com/soerenschneider/gobot-lux/actions/workflows/release.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/gobot-lux/actions/workflows/golangci-lint.yaml/badge.svg)

Reads and forwards brightness data using an analogous brightness sensor and a Raspberry PI

## Features

🤖 Integrates with Home-Assistant<br/>
📊 Calculates statistics about brightness data over time windows, accessible via MQTT and metrics<br/>
🔐 Allows connecting to secure MQTT brokers using TLS client certificates<br/>
🔭 Expose brightness data as metrics to enable alerting and Grafana dashboards<br/>

## Installation

### Binaries
Download a prebuilt binary from the [releases section](https://github.com/soerenschneider/gobot-lux/releases) for your system.

### From Source
As a prerequisite, you need to have [Golang SDK](https://go.dev/dl/) installed. Then you can install gobot-lux from source by invoking:
```shell
$ go install github.com/soerenschneider/gobot-lux@latest
```

## Configuration

gobot-lux can be fully configured using either environment variables or a config file. To supply a config file, the `-config` parameter is used.

### Configuration Reference

| Field                  | JSON Key                | Environment Variable               | Validation Rules                                        |
|------------------------|-------------------------|------------------------------------|---------------------------------------------------------|
| Placement              | placement               | GOBOT_LUX_PLACEMENT                | required                                                |
| MetricConfig           | metrics_addr            | GOBOT_LUX_METRICS_LISTEN_ADDR      | optional, must be a valid TCP address                   |
| IntervalSecs           | interval_s              | GOBOT_LUX_INTERVAL_S               | minimum: 1, maximum: 300                                |
| StatIntervals          | stat_intervals          | GOBOT_LUX_STAT_INTERVALS           | each value: minimum: 10, maximum: 3600                  |
| LogSensor              | log_sensor              | GOBOT_LUX_LOG_SENSOR_READINGS      |                                                         |
| MqttConfig             |                         |                                    |                                                         |
| - Host                 | mqtt_host               | GOBOT_LUX_MQTT_BROKER              | required, must be a valid MQTT broker                   |
| - Topic                | mqtt_topic              | GOBOT_LUX_MQTT_TOPIC               | required, must be a valid MQTT topic                    |
| - StatsTopic           | mqtt_stats_topic        | GOBOT_LUX_MQTT_STATS_TOPIC         | optional, must be a valid MQTT topic                    |
| - ClientKeyFile        | mqtt_ssl_key_file       | GOBOT_LUX_MQTT_TLS_CLIENT_KEY_FILE | required unless ClientCertFile is empty, must be a file |
| - ClientCertFile       | mqtt_ssl_cert_file      | GOBOT_LUX_MQTT_TLS_CLIENT_CRT_FILE | required unless ClientKeyFile is empty, must be a file  |
| - ServerCaFile         | mqtt_ssl_ca_file        | GOBOT_LUX_MQTT_TLS_SERVER_CA_FILE  | optional, must be a file                                |
| SensorConfig           |                         |                                    |                                                         |
| - FirmAtaPort          | firmata_port            | GOBOT_LUX_FIRMATA_PORT             | required                                                |
| - AioPin               | aio_pin                 | GOBOT_LUX_AIO_PIN                  | required, must be a number                              |
| - AioPollingIntervalMs | aio_polling_interval_ms | GOBOT_LUX_AIO_POLLING_MS           | required, minimum: 1000, maximum: 60000                 |

### Example Config

```json
{
  "aio_pin": "0",
  "aio_polling_interval_ms": 1000,
  "interval_s": 1,
  "mqtt_host": "ssl://mqtt.eclipse.org:8883",
  "mqtt_ssl_cert_file": "/etc/passwd",
  "mqtt_ssl_key_file": "/etc/passwd",
  "mqtt_stats_topic": "lux/corridor_stats",
  "mqtt_topic": "lux/corridor",
  "placement": "corridor",
  "metrics_addr": "0.0.0.0:9192"
}
```

## Metrics
This project exposes the following metrics in Open Metrics format under the prefix `gobot_lux`

| Metric Name                  | Metric Type  | Description                                                      | Labels              |
|------------------------------|--------------|------------------------------------------------------------------|---------------------|
| version                      | GaugeVec     | Version information of this robot                                | version, commit     |
| heartbeat_seconds            | Gauge        | Continuous heartbeat of this bot                                 |                     |
| brightness_level_percent     | GaugeVec     | Current sensor reading of brightness level                       | placement           |
| read_errors_total            | CounterVec   | Errors while reading the sensor                                  | placement           |
| messages_published_total     | CounterVec   | Total number of published messages via MQTT                      | placement           |
| message_publish_errors_total | CounterVec   | Total number of errors while trying to publish messages via MQTT | placement           |
| min_per_interval_percent     | GaugeVec     | Minimum sensor value during given intervals                      | interval, placement |
| max_per_interval_percent     | GaugeVec     | Maximum sensor value during given intervals                      | interval, placement |
| delta_per_interval_percent   | GaugeVec     | Delta sensor value during given intervals                        | interval, placement |
| avg_per_interval_percent     | GaugeVec     | Avg sensor value during given intervals                          | interval, placement |
| slice_entries_total          | GaugeVec     | The amount of entries in the stats slice                         | placement           |


## Changelog
The full changelog can be found [here](CHANGELOG.md)
