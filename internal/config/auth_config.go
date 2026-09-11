package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type authYamlConfig struct {
	Endpoints struct {
		Bearer map[string]struct {
			Roles []string `yaml:"roles"`
		} `yaml:"bearer"`
	} `yaml:"auth-endpoints"`
}
type AuthConfig struct {
	BearerSet map[string]map[string]struct{}
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
	cfg := &AuthConfig{
		BearerSet: make(map[string]map[string]struct{}, len(fileCfg.Endpoints.Bearer)),
	}
	for method, endpoint := range fileCfg.Endpoints.Bearer {
		roles := make(map[string]struct{}, len(endpoint.Roles))
		for _, r := range endpoint.Roles {
			roles[r] = struct{}{}
		}
		cfg.BearerSet[method] = roles
	}

	return cfg, nil
}
