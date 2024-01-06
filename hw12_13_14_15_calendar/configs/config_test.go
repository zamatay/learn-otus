package configs

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Log"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("CONFIG_PATH", "../configs/config.yaml")
			if got := NewConfig(); got.DB.Driver == "" {
				t.Errorf("NewConfig() = %v, want %v", got, "лог пустой")
			}
		})
	}
}
