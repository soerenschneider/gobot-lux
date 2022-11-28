package internal

import (
	"encoding/json"
	"fmt"
	"github.com/soerenschneider/gobot-lux/internal/config"
	"gobot.io/x/gobot"
	"log"
	"time"
)

type BrightnessDriver interface {
	Read() (val int, err error)
	Name() string
	SetName(name string)
	Start() error
	Halt() error
	Connection() gobot.Connection
	gobot.Eventer
	gobot.Commander
}

type MqttAdapter interface {
	gobot.Connection
	Publish(topic string, msg []byte) bool
}

type BrightnessBot struct {
	Driver      BrightnessDriver
	Adaptor     gobot.Connection
	MqttAdaptor MqttAdapter

	Config config.Config
}

func (m *BrightnessBot) publishMessage(topic string, msg []byte) {
	success := m.MqttAdaptor.Publish(topic, msg)
	if success {
		metricsMessagesPublished.WithLabelValues(m.Config.Placement).Inc()
	} else {
		metricsMessagePublishErrors.WithLabelValues(m.Config.Placement).Inc()
	}
}

func AssembleBot(bot *BrightnessBot) *gobot.Robot {
	metricVersionInfo.WithLabelValues(BuildVersion, CommitHash).Set(1)
	statsModule := NewSensorStats()
	var valuePercent float32
	work := func() {
		gobot.Every(60*time.Second, func() {
			metricHeartbeat.SetToCurrentTime()
		})

		gobot.Every(time.Duration(bot.Config.IntervalSecs)*time.Second, func() {
			if valuePercent >= 0 {
				msg := []byte(fmt.Sprintf("%f", valuePercent))
				bot.publishMessage(bot.Config.Topic, msg)
			}
		})

		gobot.Every(time.Duration(bot.Config.AioPollingIntervalMs)*time.Millisecond, func() {
			rawValue, err := bot.Driver.Read()
			if err != nil {
				metricSensorError.WithLabelValues(bot.Config.Placement).Inc()
				valuePercent = -1
			} else {
				valuePercent = (1024 - float32(rawValue)) * 100 / 1024
				statsModule.NewEvent(valuePercent)
				metricBrightness.WithLabelValues(bot.Config.Placement).Set(float64(valuePercent))
			}

			if bot.Config.LogSensor {
				log.Printf("Read %d from sensor (%f%%)", rawValue, valuePercent)
			}
		})

		if len(bot.Config.MqttConfig.StatsTopic) != 0 && len(bot.Config.StatIntervals) > 0 {
			min, _ := bot.Config.GetStatIntervalMin()
			max, _ := bot.Config.GetStatIntervalMax()

			gobot.Every(time.Duration(min)*time.Second, func() {
				statsDict := map[string]IntervalStatistics{}
				for _, stat := range bot.Config.StatIntervals {
					intervalStatistics, err := statsModule.GetIntervalStats(time.Duration(stat) * time.Second)
					if err != nil {
						continue
					}

					key := fmt.Sprintf("%ds", stat)
					statsDict[key] = intervalStatistics
					updateStatsIntervalMetrics(key, bot.Config.Placement, intervalStatistics)
					max = int(intervalStatistics.Max)
				}
				statsModule.PurgeStatsBefore(time.Now().Add(time.Duration(-max) * time.Second))
				metricsStatsSliceSize.WithLabelValues(bot.Config.Placement).Set(float64(statsModule.GetStatsSliceSize()))

				json, err := json.Marshal(statsDict)
				if err == nil {
					bot.publishMessage(bot.Config.StatsTopic, json)
				} else {
					log.Printf("Error while marshalling json: %v", err)
				}
			})
		}
	}

	adaptors := []gobot.Connection{bot.Adaptor}
	if bot.MqttAdaptor != nil {
		adaptors = append(adaptors, bot.MqttAdaptor)
	}

	return gobot.NewRobot(config.BotName,
		adaptors,
		[]gobot.Device{bot.Driver},
		work,
	)
}

func updateStatsIntervalMetrics(key, placement string, stats IntervalStatistics) {
	metricsStatsMin.WithLabelValues(key, placement).Set(float64(stats.Min))
	metricsStatsMax.WithLabelValues(key, placement).Set(float64(stats.Max))
	metricsStatsDelta.WithLabelValues(key, placement).Set(float64(stats.Delta))
	metricsStatsAvg.WithLabelValues(key, placement).Set(float64(stats.Avg))
}
