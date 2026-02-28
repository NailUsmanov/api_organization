// Package config предоставляет функциональность для загрузки и хранения
// конфигурации приложения из переменных окружения.
package config

import (
	"fmt"
	"os"
)

// Config хранит все настройки приложения, необходимые для его работы.
// Значения загружаются из переменных окружения или используются значения по умолчанию.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
}

// Load загружает конфигурацию из переменных окружения.
func Load() (*Config, error) {
	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "organization"),
		Port:       getEnv("PORT", "8080"),
	}
	return cfg, nil
}

// getEnv возвращает значение переменной окружения с заданным ключом.
func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}

// DSN формирует строку подключения к базе данных PostgreSQL (Data Source Name).
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
