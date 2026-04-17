package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
	Model   string `json:"model"`
}

func DefaultConfig() *Config {
	return &Config{
		BaseURL: "https://api.openai.com/v1",
		Model:   "gpt-4o-mini",
	}
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir failed: %w", err)
	}
	return filepath.Join(home, ".gocommit.json"), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			if key := os.Getenv("OPENAI_API_KEY"); key != "" {
				cfg.APIKey = key
			}
			if url := os.Getenv("OPENAI_BASE_URL"); url != "" {
				cfg.BaseURL = url
			}
			if model := os.Getenv("OPENAI_MODEL"); model != "" {
				cfg.Model = model
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config failed: %w", err)
	}

	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		cfg.APIKey = key
	}
	if url := os.Getenv("OPENAI_BASE_URL"); url != "" {
		cfg.BaseURL = url
	}
	if model := os.Getenv("OPENAI_MODEL"); model != "" {
		cfg.Model = model
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config failed: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
