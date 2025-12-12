package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config holds application configuration
type Config struct {
	// Refresh intervals
	ServiceRefreshInterval time.Duration `json:"service_refresh_interval"`
	UIRefreshInterval      time.Duration `json:"ui_refresh_interval"`
	
	// Metrics settings
	MetricsHistorySize int `json:"metrics_history_size"`
	
	// UI settings
	ShowEmptyServices bool `json:"show_empty_services"`
	ColorScheme       string `json:"color_scheme"`
	
	// Detection settings
	EnableDocker     bool `json:"enable_docker"`
	EnableKubernetes bool `json:"enable_kubernetes"`
	EnableSystemd    bool `json:"enable_systemd"`
	EnableProcess    bool `json:"enable_process"`
	
	// Process detection patterns
	ProcessPatterns []ProcessPattern `json:"process_patterns"`
	
	// Logging
	LogLevel string `json:"log_level"`
	LogFile  string `json:"log_file"`
}

// ProcessPattern defines patterns for process detection
type ProcessPattern struct {
	Name        string   `json:"name"`
	Patterns    []string `json:"patterns"`
	Ports       []string `json:"ports"`
	HealthCheck string   `json:"health_check"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		ServiceRefreshInterval: 5 * time.Second,
		UIRefreshInterval:      2 * time.Second,
		MetricsHistorySize:     60,
		ShowEmptyServices:      false,
		ColorScheme:           "default",
		EnableDocker:          true,
		EnableKubernetes:      true,
		EnableSystemd:         true,
		EnableProcess:         true,
		ProcessPatterns: []ProcessPattern{
			{
				Name:     "Node.js",
				Patterns: []string{"node", "npm", "yarn"},
				Ports:    []string{"3000", "8000", "8080"},
			},
			{
				Name:     "Python",
				Patterns: []string{"python", "python3", "gunicorn", "uvicorn"},
				Ports:    []string{"5000", "8000", "8080"},
			},
			{
				Name:     "Java",
				Patterns: []string{"java", "spring-boot"},
				Ports:    []string{"8080", "8443", "9090"},
			},
		},
		LogLevel: "INFO",
		LogFile:  "logs/lazyservice.log",
	}
}

// LoadConfig loads configuration from file or creates default
func LoadConfig() (*Config, error) {
	configPath := getConfigPath()
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// If config file doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := config.Save(); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
		return config, nil
	}
	
	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return config, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath := getConfigPath()
	
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// getConfigPath returns the configuration file path
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "config/lazyservice.json"
	}
	return filepath.Join(homeDir, ".config", "lazyservice", "config.json")
}
