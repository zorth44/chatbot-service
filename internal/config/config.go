package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	OpenRouter OpenRouterConfig `yaml:"openrouter"`
}

// OpenRouterConfig represents the OpenRouter client configuration
type OpenRouterConfig struct {
	BaseURL  string `yaml:"base_url"`
	APIKey   string `yaml:"api_key"`
	SiteURL  string `yaml:"site_url"`
	SiteName string `yaml:"site_name"`
}

// LoadConfig loads the configuration from the specified file
func LoadConfig(path string) (*Config, error) {
	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Load API key from environment variable
	if apiKey := os.Getenv("OPENROUTER_API_KEY"); apiKey != "" {
		config.OpenRouter.APIKey = apiKey
	}

	return &config, nil
}
