package config

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// This regex is not a very strict check, we don't validate hostname or ip (v4, v6) addresses...
	mqttHostRegex = regexp.MustCompile(`^\w{3,}://.{3,}:\d{2,5}$`)

	// We don't care that technically it's allowed to start with a slash
	mqttTopicRegex = regexp.MustCompile(`^([\w%]+)(/[\w%]+)*$`)
)

type MqttConfig struct {
	Host           string `json:"mqtt_host,omitempty" validate:"required,mqtt_broker"`
	Topic          string `json:"mqtt_topic,omitempty" validate:"required,mqtt_topic"`
	StatsTopic     string `json:"mqtt_stats_topic,omitempty" validate:"omitempty,mqtt_topic"`
	ClientKeyFile  string `json:"mqtt_ssl_key_file,omitempty" validate:"required_unless=ClientCertFile '',required_unless=ServerCaFile '',omitempty,file"`
	ClientCertFile string `json:"mqtt_ssl_cert_file,omitempty" validate:"required_unless=ClientKeyFile '',required_unless=ServerCaFile '',omitempty,file"`
	ServerCaFile   string `json:"mqtt_ssl_ca_file,omitempty" validate:"required_unless=ClientKeyFile '',required_unless=ClientCertFile '',omitempty,file"`
}

func (conf *MqttConfig) UsesSslCerts() bool {
	return len(conf.ClientCertFile) > 0 && len(conf.ClientKeyFile) > 0
}

func matchTopic(topic string) bool {
	return mqttTopicRegex.MatchString(topic)
}

func matchHost(host string) bool {
	return mqttHostRegex.Match([]byte(host))
}

func (conf *Config) FormatTopic() {
	if strings.Contains(conf.Topic, "%s") {
		conf.Topic = fmt.Sprintf(conf.Topic, conf.Placement)
	}
}
