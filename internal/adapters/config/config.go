// Package config.
package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	Server   ServerConfig   `mapstructure:",squash"`
	Redis    RedisConfig    `mapstructure:",squash"`
}

type AppConfig struct {
	Name string `mapstructure:"APP_NAME"`
	Env  string `mapstructure:"APP_ENV"`
	URL  string `mapstructure:"APP_URL"`
}

type DatabaseConfig struct {
	DB              string `mapstructure:"POSTGRES_DB"`
	User            string `mapstructure:"POSTGRES_USER"`
	Password        string `mapstructure:"POSTGRES_PASSWORD"`
	Host            string `mapstructure:"POSTGRES_HOST"`
	Port            string `mapstructure:"POSTGRES_PORT"`
	SSLMode         string `mapstructure:"POSTGRES_SSL_MODE"`
	MaxOpenConns    int    `mapstructure:"POSTGRES_MAX_OPEN_CONNS"`
	MaxIdleConns    int    `mapstructure:"POSTGRES_MAX_IDLE_CONNS"`
	ConnMaxLifetime int    `mapstructure:"POSTGRES_CONN_MAX_LIFETIME"`
}

type ServerConfig struct {
	Env  string `mapstructure:"APP_ENV"`
	Port string `mapstructure:"APP_PORT"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	return &config
}

func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

func (cfg *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DB,
		cfg.SSLMode,
	)
}

func (cfg *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}
