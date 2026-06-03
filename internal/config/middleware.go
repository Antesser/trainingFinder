package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type AuthConfig struct {
	Endpoints struct {
		Bearer []string `yaml:"bearer"`
	} `yaml:"auth-endpoints"`
	BearerSet map[string]struct{}
}

func (a *AuthConfig) HasAvailableMethod(method string) bool {
	for k := range a.BearerSet {
		if k == method {
			return true
		}
	}
	return false
}

// NewAuthConfig loads and parses the auth configuration from file
func NewAuthConfig(path string) (*AuthConfig, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	var cfg AuthConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	cfg.BearerSet = make(map[string]struct{}, len(cfg.Endpoints.Bearer))
	for _, m := range cfg.Endpoints.Bearer {
		cfg.BearerSet[m] = struct{}{}
	}
	return &cfg, nil
}
