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

type App struct {
	logger  logger.Logger
	logFile *os.File
	server  server.Server
	ctx     context.Context
	cancel  context.CancelFunc
	cache   cache.Cache
	storage repository.Storage
}

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

func (a *App) Run() {

	go func() {
		if err := a.server.Run(); err != nil {
			a.logger.LogFatal("server run failed", err, "layer", "app")
		}
	}()

	<-a.ctx.Done()

	a.stop()

}

func (a *App) stop() {

	a.server.Shutdown()

	a.cache.Close()
	a.storage.Close()

	if a.logFile != nil && a.logFile != os.Stdout {
		_ = a.logFile.Close()
	}

}
