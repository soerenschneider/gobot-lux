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
		Name:      "message_publish_errors",
		Subsystem: "mqtt",
		Help:      "Total number of errors while trying to publish messages via MQTT",
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
