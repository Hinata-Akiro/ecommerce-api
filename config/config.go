package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DB_CONNECTION_STRING string
	JWT_SECRET           string
	PORT                 string
	SWAGGER_SERVER_URL   string
}

var CONFIG *AppConfig

func Config() *AppConfig {
	godotenv.Load()
	appConfig := &AppConfig{
		DB_CONNECTION_STRING: os.Getenv("DB_CONNECTION_STRING"),
		JWT_SECRET:           os.Getenv("JWT_SECRET"),
		PORT:                 os.Getenv("PORT"),
		SWAGGER_SERVER_URL:   os.Getenv("SWAGGER_SERVER_URL"),
	}
	CONFIG = appConfig
	return appConfig
}
