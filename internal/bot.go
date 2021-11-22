package internal

import (
	"encoding/json"
	"fmt"
	"github.com/soerenschneider/gobot-lux/internal/config"
	"gobot.io/x/gobot"
	"log"
	"math"
	"strconv"
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
	stats := NewSensorStats()
	readValue := math.MinInt16
	work := func() {
		gobot.Every(60*time.Second, func() {
			metricHeartbeat.Inc()
		})

		gobot.Every(time.Duration(bot.Config.IntervalSecs)*time.Second, func() {
			if readValue >= 0 {
				msg := []byte(strconv.Itoa(readValue))
				bot.publishMessage(bot.Config.Topic, msg)
			}
		})

		gobot.Every(time.Duration(bot.Config.AioPollingIntervalMs)*time.Millisecond, func() {
			var err error
			readValue, err = bot.Driver.Read()
			stats.NewEvent(readValue)
			if err != nil {
				metricSensorError.WithLabelValues(bot.Config.Placement).Inc()
			} else {
				metricBrightness.WithLabelValues(bot.Config.Placement).Set(float64(readValue))
			}

			if bot.Config.LogSensor {
				log.Printf("Read %d from sensor", readValue)
			}
		})

		if len(bot.Config.MqttConfig.StatsTopic) != 0 && len(bot.Config.StatIntervals) > 0 {
			min, _ := bot.Config.GetStatIntervalMin()
			max, _ := bot.Config.GetStatIntervalMax()

			gobot.Every(time.Duration(min)*time.Second, func() {
				statsDict := map[string]Measure{}
				for _, stat := range bot.Config.StatIntervals {
					min, max, err := stats.GetEventCountNewerThan(time.Duration(stat) * time.Second)
					if err != nil {
						continue
					}

					key := fmt.Sprintf("%ds", stat)
					statsDict[key] = Measure{
						Min:   min,
						Max:   max,
						Delta: max - min,
					}
					metricsStatsMin.WithLabelValues(bot.Config.Placement).Set(float64(min))
					metricsStatsMax.WithLabelValues(bot.Config.Placement).Set(float64(max))
					metricsStatsDelta.WithLabelValues(bot.Config.Placement).Set(float64(max-min))
				}
				stats.PurgeStatsBefore(time.Now().Add(time.Duration(-max) * time.Second))
				metricsStatsSliceSize.WithLabelValues(bot.Config.Placement).Set(float64(stats.GetStatsSliceSize()))

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

type Measure struct {
	Min   int16 `json:"min"`
	Max   int16 `json:"max"`
	Delta int16 `json:"delta"`
}