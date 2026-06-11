package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type authYamlConfig struct {
	Endpoints struct {
		Bearer []string `yaml:"bearer"`
	} `yaml:"auth-endpoints"`
}
type AuthConfig struct {
	BearerSet map[string]struct{}
}

// NewAuthConfig loads and parses the auth configuration from file
func NewAuthConfig(path string) (*AuthConfig, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	var fileCfg authYamlConfig
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return nil, err
	}
	var cfg AuthConfig
	cfg.BearerSet = make(map[string]struct{}, len(fileCfg.Endpoints.Bearer))
	for _, m := range fileCfg.Endpoints.Bearer {
		cfg.BearerSet[m] = struct{}{}
	}
	return &cfg, nil
}
