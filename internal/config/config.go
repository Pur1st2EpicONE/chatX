package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Logger  Logger  `mapstructure:"logger"`
	Server  Server  `mapstructure:"server"`
	Service Service `mapstructure:"service"`
	Cache   Cache   `mapstructure:"cache"`
	Storage Storage `mapstructure:"database"`
}

type Logger struct {
	Debug          bool   `mapstructure:"debug_mode"`
	LogDir         string `mapstructure:"log_directory"`
	RequestLogging bool   `mapstructure:"request_logging"`
}

type Server struct {
	Port            string        `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type Service struct {
	MaxMessageLength int `mapstructure:"max_message_length"`
	MaxTitleLength   int `mapstructure:"max_title_length"`
	GetLimitMax      int `mapstructure:"get_limit_max"`
	GetLimitDefault  int `mapstructure:"get_limit_default"`
}

type Storage struct {
	Dialect         string        `mapstructure:"goose_dialect"`
	MigrationsDir   string        `mapstructure:"goose_migrations_directory"`
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	RecoverLimit    int           `mapstructure:"recover_limit"`
}

type Cache struct {
	Capacity    int `mapstructure:"capacity"`
	MaxMessages int `mapstructure:"max_messages"`
}

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

func loggerConfig() Logger {
	return Logger{
		Debug:          viper.GetBool("logger.debug_mode"),
		LogDir:         viper.GetString("logger.log_directory"),
		RequestLogging: viper.GetBool("logger.request_logging"),
	}
}

func serverConfig() Server {
	return Server{
		Port:            viper.GetString("server.port"),
		ReadTimeout:     viper.GetDuration("server.read_timeout"),
		WriteTimeout:    viper.GetDuration("server.write_timeout"),
		MaxHeaderBytes:  viper.GetInt("server.max_header_bytes"),
		ShutdownTimeout: viper.GetDuration("server.shutdown_timeout"),
	}
}

func serviceConfig() Service {
	return Service{
		MaxMessageLength: viper.GetInt("service.max_message_length"),
		MaxTitleLength:   viper.GetInt("service.max_title_length"),
		GetLimitMax:      viper.GetInt("service.get_limit_max"),
		GetLimitDefault:  viper.GetInt("service.get_limit_default"),
	}
}

func cacheConfig() Cache {
	return Cache{
		Capacity:    viper.GetInt("cache.capacity"),
		MaxMessages: viper.GetInt("cache.max_messages"),
	}
}

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
		RecoverLimit:    viper.GetInt("database.recover_limit"),
	}
}

func loadEnvs(conf *Config) {
	conf.Storage.Username = os.Getenv("DB_USER")
	conf.Storage.Password = os.Getenv("DB_PASSWORD")
}
