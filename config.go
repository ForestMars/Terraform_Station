package TerraformStation

import (
	"time"
)

type Config struct {
	// OpenTofu configuration
	OpenTofuPath    string        `json:"opentofu_path" yaml:"opentofu_path"`
	WorkingDirectory string        `json:"working_directory" yaml:"working_directory"`
	Timeout          time.Duration `json:"timeout" yaml:"timeout"`
	
	// Database configuration
	Database DatabaseConfig `json:"database" yaml:"database"`
	
	// Logging configuration
	LogLevel string `json:"log_level" yaml:"log_level"`
	
	// API configuration
	Port         string `json:"port" yaml:"port"`
	Host         string `json:"host" yaml:"host"`
	EnableCORS   bool   `json:"enable_cors" yaml:"enable_cors"`
}

type DatabaseConfig struct {
	Driver   string `json:"driver" yaml:"driver"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
	SSLMode  string `json:"ssl_mode" yaml:"ssl_mode"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		OpenTofuPath:    "tofu",
		WorkingDirectory: "./tofu",
		Timeout:          30 * time.Minute,
		LogLevel:         "info",
		Port:             "8080",
		Host:             "localhost",
		EnableCORS:       true,
		Database: DatabaseConfig{
			Driver:   "sqlite",
			Host:     "localhost",
			Port:     5432,
			SSLMode:  "disable",
		},
	}
}
