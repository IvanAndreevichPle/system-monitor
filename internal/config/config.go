package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Logging  LoggingConfig  `yaml:"logging"`
	Collector CollectorConfig `yaml:"collector"`
}

// ServerConfig contains server settings
type ServerConfig struct {
	Port int `yaml:"port"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level  string `yaml:"level"`  // debug, info, warn, error
	Format string `yaml:"format"` // json, text
}

// CollectorConfig contains metrics collector settings
type CollectorConfig struct {
	DefaultIntervalSeconds       int32 `yaml:"default_interval_seconds"`
	DefaultAveragingWindowSeconds int32 `yaml:"default_averaging_window_seconds"`
	CollectionIntervalSeconds     int32 `yaml:"collection_interval_seconds"`
}

// LoadConfig loads configuration from file and CLI flags
// CLI flags override config file values
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		Collector: CollectorConfig{
			DefaultIntervalSeconds:       5,
			DefaultAveragingWindowSeconds: 15,
			CollectionIntervalSeconds:     1,
		},
	}

	// Define CLI flags
	var configPath string
	port := flag.Int("port", 0, "Port to listen on (overrides config file)")
	logLevel := flag.String("log-level", "", "Log level: debug, info, warn, error (overrides config file)")
	logFormat := flag.String("log-format", "", "Log format: json, text (overrides config file)")
	flag.StringVar(&configPath, "config", "", "Path to config file")

	// Parse flags
	flag.Parse()

	// Load from file if provided
	if configPath != "" {
		if err := cfg.LoadFromFile(configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with CLI flags (only if flags were explicitly set)
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "port":
			cfg.Server.Port = *port
		case "log-level":
			cfg.Logging.Level = *logLevel
		case "log-format":
			cfg.Logging.Format = *logFormat
		}
	})

	return cfg, nil
}

// LoadFromFile loads configuration from YAML file
func (c *Config) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// SaveToFile saves configuration to YAML file
func (c *Config) SaveToFile(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
