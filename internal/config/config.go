package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	defaultHTTPPort            = "8080"
	defaultPostgresHost        = "localhost"
	defaultPostgresPort        = 5432
	defaultPostgresUser        = "postgres"
	defaultPostgresPassword    = "postgres"
	defaultPostgresDB          = "app"
	defaultPostgresSSLMode     = "disable"
	defaultJWTSecret           = "dev-secret-change-me"
	defaultJWTTTL              = 12 * time.Hour
	defaultHTTPShutdownTimeout = 10 * time.Second
)

type Config struct {
	HTTP HTTPConfig
	DB   DatabaseConfig
	JWT  JWTConfig
}

type HTTPConfig struct {
	Port            string
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	DSN      string
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

func Load() (Config, error) {
	httpShutdownTimeout, err := getDuration("HTTP_SHUTDOWN_TIMEOUT", defaultHTTPShutdownTimeout)
	if err != nil {
		return Config{}, fmt.Errorf("parse HTTP_SHUTDOWN_TIMEOUT: %w", err)
	}

	dbPort, err := getInt("POSTGRES_PORT", defaultPostgresPort)
	if err != nil {
		return Config{}, fmt.Errorf("parse POSTGRES_PORT: %w", err)
	}

	jwtTTL, err := getDuration("JWT_TTL", defaultJWTTTL)
	if err != nil {
		return Config{}, fmt.Errorf("parse JWT_TTL: %w", err)
	}

	cfg := Config{
		HTTP: HTTPConfig{
			Port:            getString("HTTP_PORT", defaultHTTPPort),
			ShutdownTimeout: httpShutdownTimeout,
		},
		DB: DatabaseConfig{
			Host:     getString("POSTGRES_HOST", defaultPostgresHost),
			Port:     dbPort,
			User:     getString("POSTGRES_USER", defaultPostgresUser),
			Password: getString("POSTGRES_PASSWORD", defaultPostgresPassword),
			Name:     getString("POSTGRES_DB", defaultPostgresDB),
			SSLMode:  getString("POSTGRES_SSLMODE", defaultPostgresSSLMode),
		},
		JWT: JWTConfig{
			Secret: getString("JWT_SECRET", defaultJWTSecret),
			TTL:    jwtTTL,
		},
	}

	cfg.DB.DSN = getString("DATABASE_DSN", cfg.DB.buildDSN())
	return cfg, nil
}

func (c DatabaseConfig) buildDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(c.User),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		url.PathEscape(c.Name),
		url.QueryEscape(c.SSLMode),
	)
}

func getString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	return strconv.Atoi(value)
}

func getDuration(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	return time.ParseDuration(value)
}
