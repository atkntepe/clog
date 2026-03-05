package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type RepoInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Config struct {
	Author string     `json:"author"`
	Model  string     `json:"model,omitempty"`
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

const defaultModel = "claude-sonnet-4-20250514"

func EnvPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "clog", ".env"), nil
}

func GetAPIKey() string {
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return key
	}
	val, _ := loadEnvValue("ANTHROPIC_API_KEY")
	return val
}

func GetModel(cfg *Config) string {
	if val := os.Getenv("ANTHROPIC_MODEL"); val != "" {
		return val
	}
	if cfg.Model != "" {
		return cfg.Model
	}
	return defaultModel
}

func SaveEnvValue(key, value string) error {
	p, err := EnvPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	env := make(map[string]string)
	if data, err := os.ReadFile(p); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if k, v, ok := strings.Cut(line, "="); ok {
				env[k] = v
			}
		}
	}

	env[key] = value

	var sb strings.Builder
	for k, v := range env {
		sb.WriteString(k + "=" + v + "\n")
	}
	return os.WriteFile(p, []byte(sb.String()), 0600)
}

func loadEnvValue(key string) (string, error) {
	p, err := EnvPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if k, v, ok := strings.Cut(line, "="); ok && k == key {
			return v, nil
		}
	}
	return "", nil
}
