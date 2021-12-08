package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soerenschneider/gobot-lux/internal/config"
	"log"
	"net/http"
)

const namespace = config.BotName

var (
	metricVersionInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "version",
		Help:      "Version information of this robot",
	}, []string{"version", "commit"})

	metricHeartbeat = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "heartbeat_seconds",
		Help:      "Continuous heartbeat of this bot",
	})

	metricBrightness = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "brightness_level",
		Subsystem: "sensor",
		Help:      "Current sensor reading of brightness level",
	}, []string{"placement"})

	metricSensorError = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "read_errors_total",
		Subsystem: "sensor",
		Help:      "Errors while reading the sensor",
	}, []string{"placement"})

	metricsMessagesPublished = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "messages_published_total",
		Subsystem: "mqtt",
		Help:      "Total number of published messages via MQTT",
	}, []string{"placement"})

	metricsMessagePublishErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "message_publish_errors_total",
		Subsystem: "mqtt",
		Help:      "Total number of errors while trying to publish messages via MQTT",
	}, []string{"placement"})

	metricsStatsMin = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "min_per_interval",
		Subsystem: "stats",
		Help:      "Minimum sensor value during given intervals",
	}, []string{"interval", "placement"})

	metricsStatsMax = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "max_per_interval",
		Subsystem: "stats",
		Help:      "Maximum sensor value during given intervals",
	}, []string{"interval", "placement"})

	metricsStatsDelta = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "delta_per_interval",
		Subsystem: "stats",
		Help:      "Delta sensor value during given intervals",
	}, []string{"interval", "placement"})

	metricsStatsAvg = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "avg_per_interval",
		Subsystem: "stats",
		Help:      "Avg sensor value during given intervals",
	}, []string{"interval", "placement"})

	metricsStatsSliceSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "slice_entries_total",
		Subsystem: "stats",
		Help:      "The amount of entries in the stats slice",
	}, []string{"placement"})
)

func StartMetricsServer(listenAddr string) {
	log.Printf("Starting metrics listener at %s", listenAddr)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("Could not start metrics listener: %v", err)
	}
}
