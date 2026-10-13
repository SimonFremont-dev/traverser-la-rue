package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
	"gopkg.in/yaml.v3"
)

// ConfigParser handles loading and parsing configuration files
type ConfigParser struct {
	configPath string
}

// NewConfigParser creates a new config parser
func NewConfigParser(configPath string) *ConfigParser {
	return &ConfigParser{
		configPath: configPath,
	}
}

// ParseConfig loads and parses the configuration file
func (cp *ConfigParser) ParseConfig() (*models.Config, error) {
	// Check if config file exists
	if _, err := os.Stat(cp.configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return models.NewDefaultConfig(), nil
	}

	// Read the config file
	data, err := os.ReadFile(cp.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse based on file extension
	ext := filepath.Ext(cp.configPath)
	switch ext {
	case ".json":
		var config models.Config
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
		return &config, nil

	case ".yaml", ".yml":
		var config models.Config
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
		return &config, nil

	default:
		// Try JSON first, then YAML
		var config models.Config
		if err := json.Unmarshal(data, &config); err != nil {
			if err := yaml.Unmarshal(data, &config); err != nil {
				return nil, fmt.Errorf("failed to parse config file (tried JSON and YAML): %w", err)
			}
		}
		return &config, nil
	}
}

// LoadConfig loads configuration from a file path with fallback to default
func LoadConfig(configPath string) (*models.Config, error) {
	parser := NewConfigParser(configPath)
	return parser.ParseConfig()
}

// GetResumePath returns the full path to the resume file
func (cp *ConfigParser) GetResumePath(config *models.Config) (string, error) {
	// Check if the path is a file
	info, err := os.Stat(config.Resume.Path)
	if err == nil && !info.IsDir() {
		// It's a file, return it directly
		return config.Resume.Path, nil
	}

	// If it's a directory, find the first resume file
	if err != nil {
		return "", fmt.Errorf("failed to access resume path: %w", err)
	}

	// Find the first resume file in the directory
	files, err := os.ReadDir(config.Resume.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read resume directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			return filepath.Join(config.Resume.Path, file.Name()), nil
		}
	}

	return "", fmt.Errorf("no resume files found in %s", config.Resume.Path)
}

// GetJobOfferPath returns the full path to the job offer file
func (cp *ConfigParser) GetJobOfferPath(config *models.Config) (string, error) {
	// Check if the path is a file
	info, err := os.Stat(config.JobOffer.Path)
	if err == nil && !info.IsDir() {
		// It's a file, return it directly
		return config.JobOffer.Path, nil
	}

	// If it's a directory, find the first job offer file
	if err != nil {
		return "", fmt.Errorf("failed to access job offer path: %w", err)
	}

	// Find the first job offer file in the directory
	files, err := os.ReadDir(config.JobOffer.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read job offer directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			return filepath.Join(config.JobOffer.Path, file.Name()), nil
		}
	}

	return "", fmt.Errorf("no job offer files found in %s", config.JobOffer.Path)
}

// GetOutputPath returns the full output path based on company and role
func GetOutputPath(config *models.Config, companyName, roleTitle string) string {
	if config.Output.UseCompanyAndRole {
		return filepath.Join(config.Output.Directory, companyName, roleTitle)
	}
	return config.Output.Directory
}
