package config

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
	FirmAtaPort          string `json:"firmata_port,omitempty" env:"FIRMATA_PORT" validate:"required"`
	AioPin               string `json:"aio_pin,omitempty" env:"AIO_PIN" validate:"required,number"`
	AioPollingIntervalMs int    `json:"aio_polling_interval_ms,omitempty" env:"AIO_POLLING_MS" validate:"required,min=1000,max=60000"`
}
