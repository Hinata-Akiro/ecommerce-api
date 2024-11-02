package config

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setTestEnvVariables() {
	_ = os.Setenv("DB_CONNECTION_STRING", "test_db_connection_string")
	_ = os.Setenv("JWT_SECRET", "test_jwt_secret")
	_ = os.Setenv("PORT", "8080")
	_ = os.Setenv("SWAGGER_SERVER_URL", "http://localhost:8080")
}

func clearTestEnvVariables() {
	_ = os.Unsetenv("DB_CONNECTION_STRING")
	_ = os.Unsetenv("JWT_SECRET")
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("SWAGGER_SERVER_URL")
}

func TestConfig(t *testing.T) {
	setTestEnvVariables()
	defer clearTestEnvVariables()

	godotenv.Load()

	appConfig := Config()

	// Assertions
	assert.NotNil(t, appConfig)
	assert.Equal(t, "test_db_connection_string", appConfig.DB_CONNECTION_STRING)
	assert.Equal(t, "test_jwt_secret", appConfig.JWT_SECRET)
	assert.Equal(t, "8080", appConfig.PORT)
	assert.Equal(t, "http://localhost:8080", appConfig.SWAGGER_SERVER_URL)

	assert.Equal(t, CONFIG, appConfig)
}
