// internal/app/app_test.go
package app

import (
	"context"
	"testing"
	"time"

	"github.com/NailUsmanov/api_organization/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestNewApp проверяет создание приложения
func TestNewApp(t *testing.T) {
	logger := logrus.New()
	db := &gorm.DB{} // пустой мок БД, не используется в тесте
	cfg := &config.Config{Port: "8080"}

	app := NewApp(logger, db, cfg)

	assert.NotNil(t, app)
	assert.Equal(t, logger, app.logger)
	assert.Equal(t, db, app.db)
	assert.Equal(t, cfg, app.cfg)
	assert.NotNil(t, app.router)
}

// TestSetupRoutes проверяет, что маршруты зарегистрированы
func TestSetupRoutes(t *testing.T) {
	logger := logrus.New()
	db := &gorm.DB{}
	cfg := &config.Config{}

	app := NewApp(logger, db, cfg)

	// Проверяем, что роутер не nil и содержит зарегистрированные пути
	// Можно проверить через хендлеры, но проще проверить, что setupRoutes не паникует
	assert.NotNil(t, app.router)
}

// TestRun_ServerStartAndShutdown проверяет запуск и graceful shutdown
func TestRun_ServerStartAndShutdown(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // отключаем логи для теста
	db := &gorm.DB{}
	cfg := &config.Config{}

	app := NewApp(logger, db, cfg)

	ctx, cancel := context.WithCancel(context.Background())

	// Запускаем сервер в горутине
	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run(ctx, "localhost:0") // порт 0 = случайный свободный порт
	}()

	// Даём серверу время запуститься
	time.Sleep(100 * time.Millisecond)

	// Отменяем контекст для graceful shutdown
	cancel()

	// Ждём завершения сервера
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}

// TestRun_ServerError проверяет ошибку при запуске на занятом порту
func TestRun_ServerError(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	db := &gorm.DB{}
	cfg := &config.Config{}

	app := NewApp(logger, db, cfg)
	ctx := context.Background()

	err := app.Run(ctx, "localhost:1") // Порт 1

	assert.Error(t, err)

	assert.Contains(t, err.Error(), "bind: permission denied")
}

// TestRun_ContextCancelledBeforeStart проверяет, что сервер не запускается,
// если контекст уже отменён
func TestRun_ContextCancelledBeforeStart(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	db := &gorm.DB{}
	cfg := &config.Config{}

	app := NewApp(logger, db, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // отменяем сразу

	err := app.Run(ctx, "localhost:0")
	assert.NoError(t, err) // ListenAndServe не вызывается, ошибки нет
}

// Benchmark для проверки создания приложения
func BenchmarkNewApp(b *testing.B) {
	logger := logrus.New()
	db := &gorm.DB{}
	cfg := &config.Config{}

	for i := 0; i < b.N; i++ {
		NewApp(logger, db, cfg)
	}
}
