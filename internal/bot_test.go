package internal

import (
	"fmt"
	"testing"
	"time"

	"github.com/soerenschneider/gobot-lux/internal/config"
)

func TestAssembleBot(t *testing.T) {
	conf := config.DefaultConfig()
	conf.AioPollingIntervalMs = 100
	conf.IntervalSecs = 1

	sensorValue := MaxSensorValue / 2.
	expected := 50. // 100% / 2 -> 50

	mqttAdaptor := &FakeMqttAdapter{}
	fakeAdaptor := &FakeAdaptor{}
	analogSensor := &DummyAnalogSensorDriver{value: sensorValue}
	station, err := NewBrightnessBot(analogSensor, fakeAdaptor, mqttAdaptor, conf)
	if err != nil {
		t.Error(err)
	}

	bot := AssembleBot(station)
	go func() {
		_ = bot.Start()
	}()

	time.Sleep(2 * time.Second)

	err = bot.Stop()
	if err != nil {
		t.Error("Error while stopping bot")
	}

	if mqttAdaptor.Message() != fmt.Sprintf("%f", expected) {
		t.Errorf("Expected to read value %f, got %s", sensorValue, mqttAdaptor.Msg)
	}
}
