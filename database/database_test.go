package database

import (
	"ecommerce-api/config"
	"os"
	"testing"
)

// TestConnect tests successful and unsuccessful connections
func TestConnect(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		wantErr bool
	}{
		{
			name:    "successful connection",
			dsn:     "postgres://Barny:Yungvicky007@localhost:5432/ecommerce-api",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DB_CONNECTION_STRING", tt.dsn)

			config := config.Config()
			config.DB_CONNECTION_STRING = os.Getenv("DB_CONNECTION_STRING")

			err := Connect()

			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if Database != nil {
					t.Errorf("Database connection should be nil on error")
				}
			} else {
				if Database == nil {
					t.Errorf("Database connection should not be nil on success")
				} else {
					// Clean up connection
					sqlDB, _ := Database.DB()
					sqlDB.Close()
				}
			}
		})
	}
}

func TestConnectPanic(t *testing.T) {
	os.Setenv("DB_CONNECTION_STRING", "invalid_dsn")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Connect() did not panic as expected")
		}
	}()

	Connect()
}
