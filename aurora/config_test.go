package aurora

import "testing"

func TestConfig_CreateAuroraClient(t *testing.T) {
	tests := []struct {
		name    string
		c       *Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.CreateAuroraClient(); (err != nil) != tt.wantErr {
				t.Errorf("Config.CreateAuroraClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
