package config

import "testing"

func TestSensorConfig_Validate(t *testing.T) {
	type fields struct {
		FirmAtaPort          string
		AioPin               string
		AioPollingIntervalMs int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name: "valid",
			fields: fields{
				FirmAtaPort:          "/dev/ttyUSB0",
				AioPin:               "5",
				AioPollingIntervalMs: 75,
			},
			wantErr: false,
		},
		{
			name: "missing pin",
			fields: fields{
				FirmAtaPort:          "/dev/ttyUSB0",
				AioPollingIntervalMs: 0,
			},
			wantErr: true,
		},
		{
			name: "missing polling",
			fields: fields{
				FirmAtaPort: "/dev/ttyUSB0",
				AioPin:      "5",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &SensorConfig{
				FirmAtaPort:          tt.fields.FirmAtaPort,
				AioPin:               tt.fields.AioPin,
				AioPollingIntervalMs: tt.fields.AioPollingIntervalMs,
			}
			if err := conf.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
