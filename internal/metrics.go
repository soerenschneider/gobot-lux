package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

const namespace = BotName

var (
	metricBrightness = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "brightness_total",
		Subsystem: "sensor",
		Help:      "Current sensor reading of brightness level",
	}, []string{"location"})

	metricSensorError = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "read_errors_total",
		Subsystem: "sensor",
		Help:      "Errors while reading the sensor",
	}, []string{"location"})

	metricsMessagesPublished = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "messages_published_total",
		Subsystem: "mqtt",
		Help:      "The assembleBot temperature in degrees Celsius",
	}, []string{"location"})

	metricsMessagePublishErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "message_publish_errors",
		Subsystem: "mqtt",
		Help:      "The assembleBot temperature in degrees Celsius",
	}, []string{"location"})
)

func StartMetricsServer(listenAddr string) {
	log.Printf("Starting metrics listener at %s", listenAddr)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("Could not start metrics listener: %v", err)
	}
}
