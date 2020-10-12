package config

import (
	"os"
	"strconv"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/rs/zerolog"
)

// Config is the main configuration structure
type Config struct {
	httpPort int
	log      *zerolog.Logger
	statsd   *statsd.Client
}

// NewConfig is generating a default Config object
func NewConfig() *Config {
	httpPort, err := strconv.Atoi(GetEnv("PORT", "8080"))

	if err != nil {
		panic(err)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()

	statsd, err := statsd.New("127.0.0.1:8125", statsd.WithNamespace("survilleray."))

	if err != nil {
		panic(err)
	}

	return &Config{
		httpPort: httpPort,
		log:      &log,
		statsd:   statsd,
	}
}

// GetEnv return the current `key` value or `fallback`.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// HTTPPort returns the port our http server should listen for
func (c *Config) HTTPPort() int {
	return c.httpPort
}

// Logger returns the current configured logger
func (c *Config) Logger() *zerolog.Logger {
	return c.log
}
