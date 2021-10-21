package internal

import (
	"github.com/soerenschneider/gobot-lux/internal/config"
	"gobot.io/x/gobot"
	"log"
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

func (m *BrightnessBot) publishMessage(msg []byte) {
	success := m.MqttAdaptor.Publish(m.Config.Topic, msg)
	if success {
		metricsMessagesPublished.WithLabelValues(m.Config.Location).Inc()
	} else {
		metricsMessagePublishErrors.WithLabelValues(m.Config.Location).Inc()
	}
}

func readValueAndDispatch(bot *BrightnessBot) {
	readValue, err := bot.Driver.Read()
	if err != nil {
		metricSensorError.WithLabelValues(bot.Config.Location).Inc()
	} else {
		metricBrightness.WithLabelValues(bot.Config.Location).Set(float64(readValue))
		bot.publishMessage([]byte(strconv.Itoa(readValue)))
	}

	if bot.Config.LogSensor {
		log.Printf("Read %d from sensor", readValue)
	}
}

func AssembleBot(bot *BrightnessBot) *gobot.Robot {
	work := func() {
		readValueAndDispatch(bot)
		gobot.Every(time.Duration(bot.Config.IntervalSecs)*time.Second, func() {
			readValueAndDispatch(bot)
		})
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
