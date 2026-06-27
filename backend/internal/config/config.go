package config

import (
	"fmt"
	"os"
)

// Config holds all runtime configuration sourced from environment variables.
type Config struct {
	HTTPPort string

	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	ServiceName    string
	ServiceVersion string
	OTLPEndpoint   string
}

func Load() Config {
	return Config{
		HTTPPort: getenv("HTTP_PORT", "8080"),

		DBHost: getenv("DB_HOST", "127.0.0.1"),
		DBPort: getenv("DB_PORT", "3306"),
		DBUser: getenv("DB_USER", "root"),
		DBPass: getenv("DB_PASS", ""),
		DBName: getenv("DB_NAME", "northwind"),

		ServiceName:    getenv("OTEL_SERVICE_NAME", "orders-api"),
		ServiceVersion: getenv("OTEL_SERVICE_VERSION", "0.1.0"),
		OTLPEndpoint:   getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
	}
}

// DSN builds the MySQL data source name used by database/sql.
func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
