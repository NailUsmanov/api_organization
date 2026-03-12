package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/NailUsmanov/api_organization/internal/app"
	"github.com/NailUsmanov/api_organization/internal/config"
	"github.com/NailUsmanov/api_organization/internal/db"
	"github.com/sirupsen/logrus"
)

func main() {
	// Настройка логгера
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.InfoLevel)

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к БД
	gormDB, err := db.InitDB(cfg.DSN())
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Запуск миграций
	if err := db.RunMigrations(cfg.DSN(), "migrations"); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Создание приложения
	application := app.NewApp(logger, gormDB, cfg)

	// Контекст с отслеживанием сигналов для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Запуск сервера
	if err := application.Run(ctx, ":"+cfg.Port); err != nil {
		logger.Fatalf("Server run failed: %v", err)
	}
}
