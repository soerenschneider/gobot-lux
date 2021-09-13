package internal

import (
	"strconv"
	"testing"
	"time"
)

func TestAssembleBot(t *testing.T) {
	conf := DefaultConfig()

	sensorValue := 75

	mqttAdaptor := &FakeMqttAdapter{}
	fakeAdaptor := &FakeAdaptor{}
	analogSensor := &DummyAnalogSensorDriver{value: sensorValue}
	station := &BrightnessBot{
		Driver:      analogSensor,
		Adaptor:     fakeAdaptor,
		MqttAdaptor: mqttAdaptor,
		Config:      conf,
	}

	bot := AssembleBot(station)
	go bot.Start()

	time.Sleep(2 * time.Second)
	err := bot.Stop()
	if err != nil {
		t.Error("Error while stopping bot")
	}

	if string(mqttAdaptor.Msg) != string([]byte(strconv.Itoa(analogSensor.value))) {
		t.Errorf("Expected to read value %d, got %s", sensorValue, mqttAdaptor.Msg)
	}
}
