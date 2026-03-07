package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ImagePath string `json:"image_path"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".tui")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return &Config{Width: 80, Height: 20}, nil
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	if err := os.MkdirAll(configDir(), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}
