package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type RepoInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Config struct {
	Author string     `json:"author"`
	Repos  []RepoInfo `json:"repos"`
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "clog", "config.json"), nil
}

func LoadConfig() (*Config, error) {
	p, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := &Config{
				Author: "",
				Repos:  []RepoInfo{},
			}
			if saveErr := SaveConfig(cfg); saveErr != nil {
				return nil, saveErr
			}
			return cfg, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Repos == nil {
		cfg.Repos = []RepoInfo{}
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	p, err := ConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}

func GetAPIKey() string {
	return os.Getenv("ANTHROPIC_API_KEY")
}
