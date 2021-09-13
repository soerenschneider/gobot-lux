package internal

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func Test_fromEnvBool1(t *testing.T) {
	key := "asdjnasdogsagsadgjsdgsdgasdgjsdg"
	resulting_key := fmt.Sprintf("%s_%s", strings.ToUpper(BotName), strings.ToUpper(key))
	os.Setenv(resulting_key, "true")
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				name: "asdjfasdighasgasdgasdg",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "default",
			args: args{
				name: "asdjfasdighasgasdgasdg",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "test",
			args: args{
				name: key,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromEnvBool(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromEnvBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fromEnvBool() got = %v, want %v", got, tt.want)
			}
		})
	}
}
