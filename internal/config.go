package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const BotName = "gobot_brightness"

// This regex is not a very strict check, we don't validate hostname or ip (v4, v6) addresses...
var mqttHostRegex = regexp.MustCompile(`\w{3,}://.{3,}:\d{2,5}`)

type Config struct {
	Location             string `json:"location,omitempty"`
	MetricConfig         string `json:"metrics_addr,omitempty"`
	FirmAtaPort          string `json:"firmata_port,omitempty"`
	AioPin               string `json:"aio_pin,omitempty"`
	AioPollingIntervalMs int    `json:"aio_polling_interval_ms,omitempty"`
	IntervalSecs         int    `json:"interval_s,omitempty"`
	LogValues            bool   `json:"log_values,omitempty"`
	MqttConfig
}

type MqttConfig struct {
	Host     string `json:"mqtt_host,omitempty"`
	ClientId string `json:"mqtt_client_id,omitempty"`
	Topic    string `json:"mqtt_topic,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		LogValues:            false,
		FirmAtaPort:          "/dev/ttyUSB0",
		AioPin:               "7",
		AioPollingIntervalMs: 75,
		IntervalSecs:         30,
		MetricConfig:         ":9400",
	}
}

func ConfigFromEnv() Config {
	conf := DefaultConfig()

	location, err := fromEnv("LOCATION")
	if err == nil {
		conf.Location = location
	}

	firmataPort, err := fromEnv("AIO_PORT")
	if err == nil {
		conf.FirmAtaPort = firmataPort
	}

	aioPin, err := fromEnv("AIO_PIN")
	if err == nil {
		conf.AioPin = aioPin
	}

	logValues, err := fromEnvBool("LOG_VALUES")
	if err == nil {
		conf.LogValues = logValues
	}

	aioPollingInterval, err := fromEnvInt("AIO_POLLING_INTERVAL_MS")
	if err == nil {
		conf.AioPollingIntervalMs = aioPollingInterval
	}

	intervalSeconds, err := fromEnvInt("INTERVAL_S")
	if err == nil {
		conf.IntervalSecs = intervalSeconds
	}

	mqttHost, err := fromEnv("MQTT_HOST")
	if err == nil {
		conf.Host = mqttHost
	}

	mqttClientId, err := fromEnv("MQTT_CLIENT_ID")
	if err == nil {
		conf.ClientId = mqttClientId
	}

	mqttTopic, err := fromEnv("MQTT_TOPIC")
	if err == nil {
		conf.Topic = mqttTopic
	}

	metricConfig, err := fromEnv("METRICS_ADDR")
	if err == nil {
		conf.MetricConfig = metricConfig
	}

	return conf
}

func ReadJsonConfig(filePath string) (*Config, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read config from file: %v", err)
	}

	ret := DefaultConfig()
	err = json.Unmarshal(fileContent, &ret)
	return &ret, err
}

func (c *Config) Validate() error {
	if len(c.Location) == 0 {
		return fmt.Errorf("empty location provided")
	}

	parsedPin, err := strconv.Atoi(c.AioPin)
	if err != nil {
		return fmt.Errorf("could not parse '%s' as pin: %v", c.AioPin, err)
	}
	if parsedPin < 0 {
		return fmt.Errorf("invalid pin provided: %d", parsedPin)
	}

	if c.IntervalSecs < 30 {
		return fmt.Errorf("invalid interval: must not be lower than 30 but is %d", c.IntervalSecs)
	}
	if c.IntervalSecs > 300 {
		return fmt.Errorf("invalid interval: mut not be greater than 300 but is %d", c.IntervalSecs)
	}

	if c.AioPollingIntervalMs < 5 {
		return fmt.Errorf("polling interval must not be smaller than 5: %d", c.AioPollingIntervalMs)
	}

	if c.AioPollingIntervalMs > 500 {
		return fmt.Errorf("polling interval too high: %d", c.AioPollingIntervalMs)
	}

	// TODO: improve check
	if strings.Index(c.MqttConfig.Topic, " ") != -1 {
		return errors.New("invalid mqtt topic provided")
	}

	if c.Topic == "" {
		return errors.New("empty topic provided")
	}

	return matchHost(c.MqttConfig.Host)
}

func (c *Config) Print() {
	log.Printf("Location=%s", c.Location)
	log.Printf("LogValues=%t", c.LogValues)
	log.Printf("MetricConfig=%s", c.MetricConfig)
	log.Printf("AioPin=%s", c.AioPin)
	log.Printf("AioPollingIntervalMs=%d", c.AioPollingIntervalMs)
	log.Printf("FirmAtaPort=%s", c.FirmAtaPort)
	log.Printf("IntervalSecs=%s", c.IntervalSecs)
	log.Printf("Host=%s", c.Host)
	log.Printf("Topic=%s", c.Topic)
	log.Printf("ClientId=%s", c.ClientId)
}

func matchHost(host string) error {
	if !mqttHostRegex.Match([]byte(host)) {
		return fmt.Errorf("invalid host format used")
	}
	return nil
}

func computeEnvName(name string) string {
	return fmt.Sprintf("%s_%s", strings.ToUpper(BotName), strings.ToUpper(name))
}

func fromEnv(name string) (string, error) {
	name = computeEnvName(name)
	val := os.Getenv(name)
	if val == "" {
		return "", errors.New("not defined")
	}
	return val, nil
}

func fromEnvInt(name string) (int, error) {
	val, err := fromEnv(name)
	if err != nil {
		return -1, err
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return -1, err
	}
	return parsed, nil
}

func fromEnvBool(name string) (bool, error) {
	val, err := fromEnv(name)
	if err != nil {
		return false, err
	}

	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}
	return parsed, nil
}
