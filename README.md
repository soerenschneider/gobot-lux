# gobot-lux
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/gobot-lux)](https://goreportcard.com/report/github.com/soerenschneider/gobot-lux)
![test-workflow](https://github.com/soerenschneider/gobot-lux/actions/workflows/test.yaml/badge.svg)
![release-workflow](https://github.com/soerenschneider/gobot-lux/actions/workflows/release.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/gobot-lux/actions/workflows/golangci-lint.yaml/badge.svg)

Detects and forwards brightness data using an analogous brightness sensor and a Raspberry PI

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
gobot-lux can be fully configured using either environment variables or a config file.

### Via Config File



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


## CHANGELOG
The changelog can be found [here](CHANGELOG.md)