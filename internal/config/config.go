// Package config provides configuration loading and parsing
// from environment variables, .env files, and YAML/JSON configs
// using Viper and godotenv.
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config aggregates all application configurations.
type Config struct {
	Logger  Logger  `mapstructure:"logger"`   // Logger configuration
	Server  Server  `mapstructure:"server"`   // HTTP server configuration
	Service Service `mapstructure:"service"`  // Application service limits
	Cache   Cache   `mapstructure:"cache"`    // In-memory cache configuration
	Storage Storage `mapstructure:"database"` // Database configuration
}

// Logger contains settings for logging behavior.
type Logger struct {
	Debug          bool   `mapstructure:"debug_mode"`      // Enable debug logging
	LogDir         string `mapstructure:"log_directory"`   // Directory to write logs
	RequestLogging bool   `mapstructure:"request_logging"` // Enable per-request logging
}

// Server contains HTTP server settings.
type Server struct {
	Port            string        `mapstructure:"port"`             // Listening port
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`     // Maximum read timeout
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`    // Maximum write timeout
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"` // Maximum size of request headers
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"` // Graceful shutdown timeout
}

// Service contains business logic constraints.
type Service struct {
	MaxMessageLength int `mapstructure:"max_message_length"` // Max length of a message
	MaxTitleLength   int `mapstructure:"max_title_length"`   // Max length of a title
	GetLimitMax      int `mapstructure:"get_limit_max"`      // Maximum GET limit
	GetLimitDefault  int `mapstructure:"get_limit_default"`  // Default GET limit
}

// Storage contains database connection settings.
type Storage struct {
	Dialect         string        `mapstructure:"goose_dialect"`              // Goose migration dialect
	MigrationsDir   string        `mapstructure:"goose_migrations_directory"` // Directory for Goose migrations
	Host            string        `mapstructure:"host"`                       // Database host
	Port            string        `mapstructure:"port"`                       // Database port
	Username        string        `mapstructure:"username"`                   // Database username
	Password        string        `mapstructure:"password"`                   // Database password
	DBName          string        `mapstructure:"dbname"`                     // Database name
	SSLMode         string        `mapstructure:"sslmode"`                    // SSL mode
	MaxOpenConns    int           `mapstructure:"max_open_conns"`             // Maximum open connections
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`             // Maximum idle connections
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`          // Connection max lifetime
}

// Cache contains in-memory caching settings.
type Cache struct {
	Capacity    int `mapstructure:"capacity"`     // Maximum number of chats to cache
	MaxMessages int `mapstructure:"max_messages"` // Maximum messages per cached chat
}

// Load reads configuration from Viper, .env, and environment variables.
// Returns a fully populated Config instance or an error.
func Load() (Config, error) {

	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("viper: %v", err)
	}

	if err := godotenv.Load(".env"); err != nil && !viper.GetBool("docker") {
		return Config{}, fmt.Errorf("godotenv: %v", err)
	}

	config := Config{
		Logger:  loggerConfig(),
		Server:  serverConfig(),
		Service: serviceConfig(),
		Cache:   cacheConfig(),
		Storage: storageConfig(),
	}

	loadEnvs(&config)

	return config, nil

}

// loggerConfig loads logger configuration from Viper.
func loggerConfig() Logger {
	return Logger{
		Debug:          viper.GetBool("logger.debug_mode"),
		LogDir:         viper.GetString("logger.log_directory"),
		RequestLogging: viper.GetBool("logger.request_logging"),
	}
}

// serverConfig loads server configuration from Viper.
func serverConfig() Server {
	return Server{
		Port:            viper.GetString("server.port"),
		ReadTimeout:     viper.GetDuration("server.read_timeout"),
		WriteTimeout:    viper.GetDuration("server.write_timeout"),
		MaxHeaderBytes:  viper.GetInt("server.max_header_bytes"),
		ShutdownTimeout: viper.GetDuration("server.shutdown_timeout"),
	}
}

// serviceConfig loads service constraints from Viper.
func serviceConfig() Service {
	return Service{
		MaxMessageLength: viper.GetInt("service.max_message_length"),
		MaxTitleLength:   viper.GetInt("service.max_title_length"),
		GetLimitMax:      viper.GetInt("service.get_limit_max"),
		GetLimitDefault:  viper.GetInt("service.get_limit_default"),
	}
}

// cacheConfig loads cache configuration from Viper.
func cacheConfig() Cache {
	return Cache{
		Capacity:    viper.GetInt("cache.capacity"),
		MaxMessages: viper.GetInt("cache.max_messages"),
	}
}

// storageConfig loads database configuration from Viper.
func storageConfig() Storage {
	return Storage{
		Dialect:         viper.GetString("database.goose_dialect"),
		MigrationsDir:   viper.GetString("database.goose_migrations_directory"),
		Host:            viper.GetString("database.host"),
		Port:            viper.GetString("database.port"),
		DBName:          viper.GetString("database.dbname"),
		SSLMode:         viper.GetString("database.sslmode"),
		MaxOpenConns:    viper.GetInt("database.max_open_conns"),
		MaxIdleConns:    viper.GetInt("database.max_idle_conns"),
		ConnMaxLifetime: viper.GetDuration("database.conn_max_lifetime"),
	}
}

// loadEnvs overrides specific configuration fields with environment variables.
func loadEnvs(conf *Config) {
	conf.Storage.Username = os.Getenv("DB_USER")
	conf.Storage.Password = os.Getenv("DB_PASSWORD")
}
