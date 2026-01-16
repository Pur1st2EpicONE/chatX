// Package app contains application bootstrap logic, dependency wiring,
// lifecycle management, and graceful shutdown handling.
package app

import (
	"chatX/internal/cache"
	"chatX/internal/config"
	"chatX/internal/handler"
	"chatX/internal/logger"
	"chatX/internal/repository"
	"chatX/internal/server"
	"chatX/internal/service"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// App represents the main application container.
// It holds all core dependencies and controls the application lifecycle.
type App struct {
	logger  logger.Logger      // Application-wide logger
	logFile *os.File           // Log file handler
	server  server.Server      // HTTP server instance
	ctx     context.Context    // Root application context
	cancel  context.CancelFunc // Context cancellation function
	cache   cache.Cache        // Cache layer implementation
	storage repository.Storage // Persistent storage layer
}

// Boot initializes the application by loading configuration,
// setting up logging, bootstrapping the database, and wiring dependencies.
func Boot() *App {

	config, err := config.Load()
	if err != nil {
		log.Fatalf("app — failed to load configs: %v", err)
	}

	logger, logFile := logger.NewLogger(config.Logger)

	db, err := bootstrapDB(logger, config.Storage)
	if err != nil {
		logger.LogFatal("app — failed to bootstrap database", err, "layer", "app")
	}

	return wireApp(db, logger, logFile, config)

}

// bootstrapDB initializes the database connection,
// applies migrations, and returns a configured gorm.DB instance.
func bootstrapDB(logger logger.Logger, config config.Storage) (*gorm.DB, error) {

	gormDB, err := repository.ConnectDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB via gorm: %w", err)
	}

	logger.LogInfo("app — connected to database", "layer", "app")

	if err := goose.SetDialect(config.Dialect); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	if err := goose.Up(db, config.MigrationsDir); err != nil {
		return nil, fmt.Errorf("failed to apply goose migrations: %w", err)
	}

	logger.Debug("app — migrations applied", "layer", "app")

	return gormDB, nil

}

// wireApp constructs the App instance by wiring together
// all infrastructure, domain services, handlers, and the server.
func wireApp(db *gorm.DB, logger logger.Logger, logFile *os.File, config config.Config) *App {

	ctx, cancel := newContext(logger)
	storge := repository.NewStorage(logger, config.Storage, db)
	cache := cache.NewCache(logger, config.Cache)
	service := service.NewService(logger, config.Service, cache, storge)
	handler := handler.NewHandler(logger, config.Logger.RequestLogging, service)
	server := server.NewServer(logger, config.Server, handler)

	return &App{
		logger:  logger,
		logFile: logFile,
		server:  server,
		ctx:     ctx,
		cancel:  cancel,
		cache:   cache,
		storage: storge,
	}

}

// newContext creates a root application context that is cancelled
// when an OS termination signal is received.
func newContext(logger logger.Logger) (context.Context, context.CancelFunc) {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-sigCh
		sigString := sig.String()
		if sig == syscall.SIGTERM {
			sigString = "terminate" // sig.String() returns the SIGTERM string in past tense for some reason
		}
		logger.LogInfo("app — received signal "+sigString+", initiating graceful shutdown", "layer", "app")
		cancel()
	}()

	return ctx, cancel

}

// Run starts the HTTP server and blocks until
// the application context is cancelled.
func (a *App) Run() {

	go func() {
		if err := a.server.Run(); err != nil {
			a.logger.LogFatal("server run failed", err, "layer", "app")
		}
	}()

	<-a.ctx.Done()

	a.stop()

}

// stop gracefully shuts down all application resources.
func (a *App) stop() {

	a.server.Shutdown()

	a.cache.Close()
	a.storage.Close()

	if a.logFile != nil && a.logFile != os.Stdout {
		_ = a.logFile.Close()
	}

}
