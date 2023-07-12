package config

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

const (
	defaultFirmataPort          = "/dev/ttyUSB0"
	defaultAioPin               = "0"
	defaultAioPollingIntervalMs = 750
)

func defaultSensorConfig() SensorConfig {
	return SensorConfig{
		FirmAtaPort:          defaultFirmataPort,
		AioPin:               defaultAioPin,
		AioPollingIntervalMs: defaultAioPollingIntervalMs,
	}
}

type SensorConfig struct {
	FirmAtaPort          string `json:"firmata_port,omitempty" validate:"required"`
	AioPin               string `json:"aio_pin,omitempty" validate:"required,number"`
	AioPollingIntervalMs int    `json:"aio_polling_interval_ms,omitempty" validate:"required,min=1000,max=60000"`
}

func (conf *SensorConfig) Validate() error {
	parsedPin, err := strconv.Atoi(conf.AioPin)
	if err != nil {
		return fmt.Errorf("could not parse '%s' as pin: %v", conf.AioPin, err)
	}
	if parsedPin < 0 {
		return fmt.Errorf("invalid pin provided: %d", parsedPin)
	}

	if conf.AioPollingIntervalMs < 1000 {
		return fmt.Errorf("polling interval must not be smaller than 1000: %d", conf.AioPollingIntervalMs)
	}

	if conf.AioPollingIntervalMs > 60*1000 {
		return fmt.Errorf("polling interval too high: %d", conf.AioPollingIntervalMs)
	}

	if conf.FirmAtaPort == "" {
		return errors.New("missing firmAtaPort")
	}

	return nil
}

func (conf *SensorConfig) ConfigFromEnv() {
	firmataPort, err := fromEnv("FIRMATA_PORT")
	if err == nil {
		conf.FirmAtaPort = firmataPort
	}

	aioPin, err := fromEnv("AIO_PIN")
	if err == nil {
		conf.AioPin = aioPin
	}

	aioPollingInterval, err := fromEnvInt("AIO_POLLING_INTERVAL_MS")
	if err == nil {
		conf.AioPollingIntervalMs = aioPollingInterval
	}
}

func (conf *SensorConfig) Print() {
	log.Printf("AioPin=%s", conf.AioPin)
	log.Printf("AioPollingIntervalMs=%d", conf.AioPollingIntervalMs)
	log.Printf("FirmAtaPort=%s", conf.FirmAtaPort)
}
