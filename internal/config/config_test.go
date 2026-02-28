// internal/config/config_test.go
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	// Очищаем переменные окружения перед тестом
	os.Clearenv()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "postgres", cfg.DBPassword)
	assert.Equal(t, "organization", cfg.DBName)
	assert.Equal(t, "8080", cfg.Port)
}

func TestLoad_WithEnvVars(t *testing.T) {
	// Устанавливаем тестовые переменные окружения
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("PORT", "9090")

	defer os.Clearenv() // очищаем после теста

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "testhost", cfg.DBHost)
	assert.Equal(t, "5433", cfg.DBPort)
	assert.Equal(t, "testuser", cfg.DBUser)
	assert.Equal(t, "testpass", cfg.DBPassword)
	assert.Equal(t, "testdb", cfg.DBName)
	assert.Equal(t, "9090", cfg.Port)
}

func TestLoad_PartialEnvVars(t *testing.T) {
	// Устанавливаем только некоторые переменные
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_USER", "testuser")
	defer os.Clearenv()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "testhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort) // должно быть значение по умолчанию
	assert.Equal(t, "testuser", cfg.DBUser)
	assert.Equal(t, "postgres", cfg.DBPassword) // значение по умолчанию
	assert.Equal(t, "organization", cfg.DBName) // значение по умолчанию
	assert.Equal(t, "8080", cfg.Port)           // значение по умолчанию
}

func TestDSN_Format(t *testing.T) {
	cfg := &Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "password",
		DBName:     "mydb",
	}

	expected := "host=localhost port=5432 user=postgres password=password dbname=mydb sslmode=disable"
	assert.Equal(t, expected, cfg.DSN())
}

func TestDSN_WithSpecialChars(t *testing.T) {
	cfg := &Config{
		DBHost:     "127.0.0.1",
		DBPort:     "5432",
		DBUser:     "user with spaces",
		DBPassword: "pass!@#$",
		DBName:     "db name",
	}

	// Проверяем, что строка формируется корректно (пробелы остаются как есть)
	expected := "host=127.0.0.1 port=5432 user=user with spaces password=pass!@#$ dbname=db name sslmode=disable"
	assert.Equal(t, expected, cfg.DSN())
}

func TestGetEnv_Default(t *testing.T) {
	os.Clearenv()
	result := getEnv("NONEXISTENT", "default")
	assert.Equal(t, "default", result)
}

func TestGetEnv_Existing(t *testing.T) {
	os.Setenv("EXISTENT", "value")
	defer os.Unsetenv("EXISTENT")

	result := getEnv("EXISTENT", "default")
	assert.Equal(t, "value", result)
}
