package internal

import (
	"log"
	"sync"

	"gobot.io/x/gobot/v2"
)

type DummyAnalogSensorDriver struct {
	connection gobot.Connection
	value      float64
	gobot.Eventer
	gobot.Commander
}

func NewAnalogSensorDriver(conn gobot.Connection) *DummyAnalogSensorDriver {
	return &DummyAnalogSensorDriver{
		connection: conn,
	}
}

func (a *DummyAnalogSensorDriver) Start() (err error) {
	return nil
}

func (a *DummyAnalogSensorDriver) Halt() (err error) {
	return nil
}

func (a *DummyAnalogSensorDriver) Name() string { return "DummyAnalogSensorDriver" }

func (a *DummyAnalogSensorDriver) SetName(n string) {}

func (a *DummyAnalogSensorDriver) Pin() string { return "5" }

func (a *DummyAnalogSensorDriver) Connection() gobot.Connection {
	return a.connection
}

func (a *DummyAnalogSensorDriver) Read() (val float64, err error) {
	return a.value, nil
}

type FakeMqttAdapter struct {
	Msg   []byte
	Topic string
	mutex sync.Mutex
}

func (m *FakeMqttAdapter) Name() string {
	return "FakeMqttAdapter"
}
func (m *FakeMqttAdapter) SetName(n string) {

}
func (m *FakeMqttAdapter) Connect() error {
	return nil
}
func (m *FakeMqttAdapter) Message() string {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return string(m.Msg)
}

func (m *FakeMqttAdapter) Finalize() error {
	return nil
}

func (m *FakeMqttAdapter) Publish(topic string, msg []byte) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Topic = topic
	m.Msg = msg
	log.Printf("%s -> %v", topic, string(msg))
	return true
}

type FakeAdaptor struct {
}

func (driver *FakeAdaptor) Name() string {
	return "FakeAdaptor"
}

func (driver *FakeAdaptor) SetName(n string) {
}

func (driver *FakeAdaptor) Connect() error {
	return nil
}
func (driver *FakeAdaptor) Finalize() error {
	return nil
}
