// Package db предоставляет функциональность для подключения к базе данных
// и управления миграциями.
package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB инициализирует подключение к базе данных PostgreSQL через GORM.
func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

// RunMigrations выполняет миграции базы данных с помощью goose.
// Применяет все миграции из указанной директории к базе данных.
// Использует стандартный database/sql, а не GORM, так как goose
// работает напрямую с *sql.DB.
func RunMigrations(dsn string, migrationsDir string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open sql.DB: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Migrations applied successfully")
	return nil
}
