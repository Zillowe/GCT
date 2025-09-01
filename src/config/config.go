package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type CacheConfig struct {
	Enabled bool `yaml:"enabled" envconfig:"GCT_CACHE_ENABLED"`
}

type GuidesConfig struct {
	Paths []string `yaml:"guides"`
}

type Config struct {
	Name               string       `yaml:"name" envconfig:"GCT_NAME"`
	Provider           string       `yaml:"provider" envconfig:"GCT_PROVIDER"`
	Model              string       `yaml:"model" envconfig:"GCT_MODEL"`
	APIKey             string       `yaml:"api" envconfig:"GCT_API_KEY"`
	Endpoint           string       `yaml:"endpoint,omitempty" envconfig:"GCT_ENDPOINT"`
	Commits            GuidesConfig `yaml:"commits"`
	Changelogs         GuidesConfig `yaml:"changelogs"`
	GCPProjectID       string       `yaml:"gcp_project_id,omitempty" envconfig:"GCT_GCP_PROJECT_ID"`
	GCPRegion          string       `yaml:"gcp_region,omitempty" envconfig:"GCT_GCP_REGION"`
	AWSAccessKeyID     string       `yaml:"aws_access_key_id,omitempty" envconfig:"GCT_AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey string       `yaml:"aws_secret_access_key,omitempty" envconfig:"GCT_AWS_SECRET_ACCESS_KEY"`
	AWSRegion          string       `yaml:"aws_region,omitempty" envconfig:"GCT_AWS_REGION"`
	AzureResourceName  string       `yaml:"azure_resource_name,omitempty" envconfig:"GCT_AZURE_RESOURCE_NAME"`
	Cache              CacheConfig  `yaml:"cache,omitempty"`
}

func loadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml config %s: %w", path, err)
	}
	return &config, nil
}

func findLocalConfig() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}

	for {
		configPath := filepath.Join(dir, "gct.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, true
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", false
		}
		dir = parentDir
	}
}

func findGlobalConfig() (string, bool) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", false
	}

	configPath := filepath.Join(configDir, "gct", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		return configPath, true
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	legacyPath := filepath.Join(homeDir, ".gct", "config.yaml")
	if _, err := os.Stat(legacyPath); err == nil {
		return legacyPath, true
	}

	return "", false
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	var cfg *Config
	var err error

	if localPath, found := findLocalConfig(); found {
		cfg, err = loadConfigFromFile(localPath)
		if err != nil {
			return nil, err
		}
	}

	if cfg == nil {
		if globalPath, found := findGlobalConfig(); found {
			cfg, err = loadConfigFromFile(globalPath)
			if err != nil {
				return nil, err
			}
		}
	}

	if cfg == nil {
		cfg = &Config{}
	}

	err = envconfig.Process("gct", cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	if cfg.Provider == "" {
		return nil, fmt.Errorf("no AI provider configured. Please run 'gct init' or configure environment variables")
	}

	return cfg, nil
}
