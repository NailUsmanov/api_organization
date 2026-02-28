// Package app отвечает за конфигурацию, инициализацию и запуск HTTP-приложения.
package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/NailUsmanov/api_organization/internal/config"
	"github.com/NailUsmanov/api_organization/internal/handlers"
	"github.com/NailUsmanov/api_organization/internal/middleware"
	"github.com/NailUsmanov/api_organization/internal/repository"
	"github.com/NailUsmanov/api_organization/internal/service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// App - состоит из маршуртизатора, логгера, базы данных, конфига.
type App struct {
	router *http.ServeMux
	logger *logrus.Logger
	db     *gorm.DB
	cfg    *config.Config
}

// NewApp создаёт экземпляр приложения, инициализирует зависимости и регистрирует маршруты.
func NewApp(logger *logrus.Logger, db *gorm.DB, cfg *config.Config) *App {
	router := http.NewServeMux()
	app := &App{
		router: router,
		logger: logger,
		db:     db,
		cfg:    cfg,
	}
	app.setupRoutes()
	return app
}

// setupRoutes регистрирует все обработчики и middleware.
func (a *App) setupRoutes() {
	deptRepo := repository.NewDepartmentRepository(a.db)
	empRepo := repository.NewEmployeeRepo(a.db)

	deptService := service.NewDepartmentService(deptRepo, empRepo)
	empService := service.NewEmpService(empRepo, deptRepo)

	deptHandler := handlers.NewDepartmentHandler(deptService)
	empHandler := handlers.NewEmployeeHandler(empService)

	a.router.HandleFunc("POST /departments", deptHandler.CreateDepartment)
	a.router.HandleFunc("GET /departments/{id}", deptHandler.GetDepartment)
	a.router.HandleFunc("PATCH /departments/{id}", deptHandler.UpdateDepartment)
	a.router.HandleFunc("DELETE /departments/{id}", deptHandler.DeleteDepartment)
	a.router.HandleFunc("POST /departments/{id}/employees", empHandler.CreateEmployee)
}

// Run запускает HTTP-сервер и корректно завершает его при получении сигнала.
func (a *App) Run(ctx context.Context, addr string) error {
	handler := middleware.Logger(a.logger)(a.router)

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		a.logger.Info("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	a.logger.Infof("Starting server on %s", addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
