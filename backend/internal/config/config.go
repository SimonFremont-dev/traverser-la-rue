package config

import (
	"encoding/json"
	"os"
)

// Config holds application configuration
type Config struct {
	Port          string `json:"port"`
	SQLitePath    string `json:"sqlite_path"`
	SQLDSN        string `json:"sql_dsn,omitempty"`
	KeywordCategories []string `json:"keyword_categories"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Port:           "8080",
		SQLitePath:     "data/personal.db",
		SQLDSN:        "", // Will be configured via env for SQL database
		KeywordCategories: []string{
			"programming",
			"tools",
			"frameworks",
			"databases",
			"cloud",
			"devops",
			"network",
			"security",
			"methodologies",
			"testing",
		},
	}
}

// LoadConfig loads configuration from file or uses defaults
func LoadConfig(path string) (*Config, error) {
	config := DefaultConfig()
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, nil
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	
	if err := json.Unmarshal(data, config); err != nil {
		return config, err
	}
	
	return config, nil
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	config := DefaultConfig()
	
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}
	
	if sqlitePath := os.Getenv("SQLITE_PATH"); sqlitePath != "" {
		config.SQLitePath = sqlitePath
	}
	
	if sqlDSN := os.Getenv("SQL_DSN"); sqlDSN != "" {
		config.SQLDSN = sqlDSN
	}
	
	return config
}
