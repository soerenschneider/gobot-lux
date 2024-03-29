package config

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_matchHost(t *testing.T) {
	tests := []struct {
		name string
		host string
		want bool
	}{
		{
			name: "no tld",
			host: "tcp://hostname:1883",
			want: true,
		},
		{
			name: "tld",
			host: "tcp://hostname.my.tld:1883",
			want: true,
		},
		{
			name: "ip",
			host: "tcp://192.168.0.1:1883",
			want: true,
		},
		{
			name: "no protocol",
			host: "192.168.0.1:1883",
			want: false,
		},
		{
			name: "no port",
			host: "tcp://host",
			want: false,
		},
		{
			name: "only host",
			host: "host",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchHost(tt.host); (got) != tt.want {
				t.Errorf("matchHost() error = %v, wantErr %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		Placement            string
		MetricConfig         string
		FirmAtaPort          string
		AioPin               string
		AioPollingIntervalMs int
		IntervalSecs         int
		LogValues            bool
		MqttConfig           MqttConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "all okay",
			fields: fields{
				Placement:            "placement",
				MetricConfig:         "0.0.0.0:9100",
				FirmAtaPort:          "/dev/ttyUSB0",
				AioPin:               "5",
				AioPollingIntervalMs: 7005,
				IntervalSecs:         30,
				LogValues:            false,
				MqttConfig: MqttConfig{
					Host:  "tcp://host:80",
					Topic: "topic/bla",
				},
			},
			wantErr: false,
		},
		{
			name: "missing loc",
			fields: fields{
				MetricConfig:         "127.0.0.1:9100",
				FirmAtaPort:          "/dev/ttyUSB0",
				AioPin:               "5",
				AioPollingIntervalMs: 7005,
				IntervalSecs:         30,
				LogValues:            false,
				MqttConfig: MqttConfig{
					Host:  "tcp://host:80",
					Topic: "topic/bla",
				},
			},
			wantErr: true,
		},
		{
			name: "missing firmata",
			fields: fields{
				Placement:            "loc",
				MetricConfig:         "127.0.0.1:9100",
				AioPin:               "5",
				AioPollingIntervalMs: 7005,
				IntervalSecs:         30,
				LogValues:            false,
				MqttConfig: MqttConfig{
					Host:  "tcp://host:80",
					Topic: "topic/bla",
				},
			},
			wantErr: true,
		},
		{
			name: "missing host",
			fields: fields{
				Placement:            "loc",
				MetricConfig:         "127.0.0.1:9100",
				FirmAtaPort:          "/dev/ttyUSB0",
				AioPin:               "5",
				AioPollingIntervalMs: 1000,
				IntervalSecs:         30,
				LogValues:            false,
				MqttConfig: MqttConfig{
					Topic: "topic/bla",
				},
			},
			wantErr: true,
		},
		{
			name: "all tls client certs opts set",
			fields: fields{
				Placement:            "loc",
				MetricConfig:         "127.0.0.1:9100",
				FirmAtaPort:          "/dev/ttyUSB0",
				AioPin:               "5",
				AioPollingIntervalMs: 1000,
				IntervalSecs:         30,
				LogValues:            false,
				MqttConfig: MqttConfig{
					Host:           "tcp://host:1883",
					Topic:          "topic/bla",
					ClientKeyFile:  "/etc/passwd",
					ClientCertFile: "/etc/passwd",
					ServerCaFile:   "/etc/passwd",
				},
			},
			wantErr: false,
		},
		{
			name: "single tls client cert missing",
			fields: fields{
				Placement:            "loc",
				MetricConfig:         "127.0.0.1:9100",
				FirmAtaPort:          "/dev/ttyUSB0",
				AioPin:               "5",
				AioPollingIntervalMs: 1000,
				IntervalSecs:         30,
				LogValues:            false,
				MqttConfig: MqttConfig{
					Host:           "tcp://host:1883",
					Topic:          "topic/bla",
					ClientKeyFile:  "/etc/passwd",
					ClientCertFile: "",
					ServerCaFile:   "/etc/passwd",
				},
			},
			wantErr: true,
		},
		{
			name: "all tls client opts set but one is not a file",
			fields: fields{
				Placement:            "loc",
				MetricConfig:         "127.0.0.1:9100",
				FirmAtaPort:          "/dev/ttyUSB0",
				AioPin:               "5",
				AioPollingIntervalMs: 1000,
				IntervalSecs:         30,
				LogValues:            false,
				MqttConfig: MqttConfig{
					Host:           "tcp://host:1883",
					Topic:          "topic/bla",
					ClientKeyFile:  "imnotafileman",
					ClientCertFile: "/etc/passwd",
					ServerCaFile:   "/etc/passwd",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Placement:    tt.fields.Placement,
				MetricConfig: tt.fields.MetricConfig,
				SensorConfig: SensorConfig{
					FirmAtaPort:          tt.fields.FirmAtaPort,
					AioPin:               tt.fields.AioPin,
					AioPollingIntervalMs: tt.fields.AioPollingIntervalMs,
				},
				IntervalSecs: tt.fields.IntervalSecs,
				LogSensor:    tt.fields.LogValues,
				MqttConfig:   tt.fields.MqttConfig,
			}
			if err := Validate(c); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadJsonConfig(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     *Config
		wantErr  bool
	}{
		{
			name:     "non-existent-file",
			filePath: "ihopethispathdoesntexist/somefile.json",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "example-config",
			filePath: "../../contrib/example-config-base.json",
			want: &Config{
				Placement:     "loc",
				MetricConfig:  ":1111",
				IntervalSecs:  45,
				LogSensor:     true,
				StatIntervals: []int{1, 2, 3},
				MqttConfig: MqttConfig{
					Host:       "tcp://host:1883",
					Topic:      "sensors/%s/sub",
					StatsTopic: "sensors/stats_topic",
				},
				SensorConfig: SensorConfig{
					FirmAtaPort:          "/dev/my-device",
					AioPin:               "42",
					AioPollingIntervalMs: 25,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readAndValidate(t *testing.T) {
	path := "../../contrib/real.json"
	conf, err := Read(path)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(conf)
	if err := Validate(conf); err != nil {
		t.Error(err)
	}
}

func Test_matchTopic(t *testing.T) {
	tests := []struct {
		name  string
		topic string
		want  bool
	}{
		{
			topic: "topicname",
			want:  true,
		},
		{
			topic: "more/complicated",
			want:  true,
		},
		{
			topic: "more/complicated/topic",
			want:  true,
		},
		{
			topic: "/leading",
			want:  false,
		},
		{
			topic: "trailing/",
			want:  false,
		},
		{
			topic: "replace/%s",
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchTopic(tt.topic); got != tt.want {
				t.Errorf("matchTopic() error = %v, wantErr %v", got, tt.want)
			}
		})
	}
}

func TestConfig_TemplateTopic(t *testing.T) {
	type fields struct {
		Placement  string
		MqttConfig MqttConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   *Config
	}{
		{
			name: "template",
			fields: fields{
				Placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix/%s",
				},
			},
			want: &Config{
				Placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix/loc",
				},
			},
		},
		{
			name: "no templating",
			fields: fields{
				Placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix",
				},
			},
			want: &Config{
				Placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Config{
				Placement:  tt.fields.Placement,
				MqttConfig: tt.fields.MqttConfig,
			}
			conf.FormatTopic()
			if !reflect.DeepEqual(conf, tt.want) {
				t.Fail()
			}
		})
	}
}

func TestConfig_Validate1(t *testing.T) {
	type fields struct {
		Placement     string
		MetricConfig  string
		IntervalSecs  int
		StatIntervals []int
		LogSensor     bool
		MqttConfig    MqttConfig
		SensorConfig  SensorConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Placement:     "kitchen",
				MetricConfig:  "127.0.0.1:9100",
				IntervalSecs:  60,
				StatIntervals: nil,
				LogSensor:     false,
				MqttConfig: MqttConfig{
					Host:           "tcp://broker:1883",
					Topic:          "lux",
					StatsTopic:     "stats_lux",
					ClientKeyFile:  "/etc/passwd",
					ClientCertFile: "/etc/passwd",
					ServerCaFile:   "/etc/passwd",
				},
				SensorConfig: SensorConfig{
					FirmAtaPort:          "port",
					AioPin:               "8",
					AioPollingIntervalMs: 1000,
				},
			},
			wantErr: false,
		},
		{
			name: "valid - no tls",
			fields: fields{
				Placement:     "kitchen",
				MetricConfig:  "127.0.0.1:9100",
				IntervalSecs:  60,
				StatIntervals: nil,
				LogSensor:     false,
				MqttConfig: MqttConfig{
					Host:       "tcp://broker:1883",
					Topic:      "lux",
					StatsTopic: "stats_lux",
				},
				SensorConfig: SensorConfig{
					FirmAtaPort:          "port",
					AioPin:               "8",
					AioPollingIntervalMs: 1000,
				},
			},
			wantErr: false,
		},
		{
			name: "valid - no tls, no metrics",
			fields: fields{
				Placement:     "kitchen",
				IntervalSecs:  60,
				StatIntervals: nil,
				LogSensor:     false,
				MqttConfig: MqttConfig{
					Host:       "tcp://broker:1883",
					Topic:      "lux",
					StatsTopic: "stats_lux",
				},
				SensorConfig: SensorConfig{
					FirmAtaPort:          "port",
					AioPin:               "8",
					AioPollingIntervalMs: 1000,
				},
			},
			wantErr: false,
		},
		{
			name: "valid - no tls, no metrics",
			fields: fields{
				Placement:     "kitchen",
				IntervalSecs:  60,
				StatIntervals: nil,
				LogSensor:     false,
				MqttConfig: MqttConfig{
					Host:       "ssl://broker:1883",
					Topic:      "lux/%s",
					StatsTopic: "stats_lux",
				},
				SensorConfig: SensorConfig{
					FirmAtaPort:          "port",
					AioPin:               "8",
					AioPollingIntervalMs: 1000,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Config{
				Placement:     tt.fields.Placement,
				MetricConfig:  tt.fields.MetricConfig,
				IntervalSecs:  tt.fields.IntervalSecs,
				StatIntervals: tt.fields.StatIntervals,
				LogSensor:     tt.fields.LogSensor,
				MqttConfig:    tt.fields.MqttConfig,
				SensorConfig:  tt.fields.SensorConfig,
			}
			if err := Validate(conf); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
