package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/soerenschneider/gobot-lux/internal/config"
	"gobot.io/x/gobot/v2"
)

const MaxSensorValue = 1024.

var deviationPercent float64 = 5

type BrightnessDriver interface {
	Read() (val float64, err error)
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
	Driver           BrightnessDriver
	Adaptor          gobot.Connection
	MqttAdaptor      MqttAdapter
	statsModule      *SensorStats
	mutex            *sync.RWMutex
	valuePercent     float64
	valuePercentSent float64

	Config config.Config
}

func NewBrightnessBot(driver BrightnessDriver, adaptor gobot.Connection, mqtt MqttAdapter, conf config.Config) (*BrightnessBot, error) {
	if driver == nil {
		return nil, errors.New("nil driver passed")
	}
	if adaptor == nil {
		return nil, errors.New("nil connection passed")
	}

	return &BrightnessBot{
		Driver:      driver,
		Adaptor:     adaptor,
		MqttAdaptor: mqtt,
		statsModule: NewSensorStats(),
		mutex:       &sync.RWMutex{},
		Config:      conf,
	}, nil
}

func (bot *BrightnessBot) publishMessage(topic string, msg []byte) {
	success := bot.MqttAdaptor.Publish(topic, msg)
	if success {
		metricsMessagesPublished.WithLabelValues(bot.Config.Placement).Inc()
	} else {
		metricsMessagePublishErrors.WithLabelValues(bot.Config.Placement).Inc()
	}
}

// exceedsDeviation compares the newest sensor value with the last sent
// and the last read value and returns whether they exceed a deviation threshold.
// Using this function as decision on whether to dispatch a sensor value immediately
// can be used to set a higher IntervalSecs value but still get deviated values
// immediately.
func exceedsDeviation(prevReading, prevSent, now float64) bool {
	if math.Abs(prevReading-now) >= deviationPercent {
		return true
	}

	if math.Abs(prevSent-now) >= deviationPercent {
		return true
	}

	return false
}

func (bot *BrightnessBot) updateStats() {
	bot.mutex.RLock()
	defer bot.mutex.RUnlock()

	if bot.valuePercent >= 0 {
		bot.valuePercentSent = bot.valuePercent
		msg := []byte(fmt.Sprintf("%f", bot.valuePercent))
		bot.publishMessage(bot.Config.Topic, msg)
	}
}

func (bot *BrightnessBot) updateValue() {
	bot.mutex.Lock()
	defer bot.mutex.Unlock()

	rawValue, err := bot.Driver.Read()
	if err != nil {
		metricSensorError.WithLabelValues(bot.Config.Placement).Inc()
		bot.valuePercent = -1
	} else {
		prevValue := bot.valuePercent
		bot.valuePercent = (MaxSensorValue - rawValue) * 100 / MaxSensorValue
		if exceedsDeviation(prevValue, bot.valuePercentSent, bot.valuePercent) {
			msg := []byte(fmt.Sprintf("%f", bot.valuePercent))
			bot.publishMessage(bot.Config.Topic, msg)
		}
		bot.statsModule.NewEvent(float32(bot.valuePercent))
		metricBrightness.WithLabelValues(bot.Config.Placement).Set(bot.valuePercent)
	}

	if bot.Config.LogSensor {
		log.Printf("Read %f from sensor (%f%%)", rawValue, bot.valuePercent)
	}
}

func (bot *BrightnessBot) sendStats() {
	max, _ := bot.Config.GetStatIntervalMax()
	statsDict := map[string]IntervalStatistics{}
	for _, stat := range bot.Config.StatIntervals {
		intervalStatistics, err := bot.statsModule.GetIntervalStats(time.Duration(stat) * time.Second)
		if err != nil {
			continue
		}

		key := fmt.Sprintf("%ds", stat)
		statsDict[key] = intervalStatistics
		updateStatsIntervalMetrics(key, bot.Config.Placement, intervalStatistics)
		max = int(intervalStatistics.Max)
	}
	bot.statsModule.PurgeStatsBefore(time.Now().Add(time.Duration(-max) * time.Second))
	metricsStatsSliceSize.WithLabelValues(bot.Config.Placement).Set(float64(bot.statsModule.GetStatsSliceSize()))

	json, err := json.Marshal(statsDict)
	if err == nil {
		bot.publishMessage(bot.Config.StatsTopic, json)
	} else {
		log.Printf("Error while marshalling json: %v", err)
	}
}

func AssembleBot(bot *BrightnessBot) *gobot.Robot {
	metricVersionInfo.WithLabelValues(BuildVersion, CommitHash).Set(1)

	work := func() {
		gobot.Every(60*time.Second, func() {
			metricHeartbeat.SetToCurrentTime()
		})

		gobot.Every(time.Duration(bot.Config.IntervalSecs)*time.Second, bot.updateStats)

		gobot.Every(time.Duration(bot.Config.AioPollingIntervalMs)*time.Millisecond, bot.updateValue)

		if len(bot.Config.MqttConfig.StatsTopic) != 0 && len(bot.Config.StatIntervals) > 0 {
			min, _ := bot.Config.GetStatIntervalMin()
			gobot.Every(time.Duration(min)*time.Second, bot.sendStats)
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
