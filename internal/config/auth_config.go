package config

import (
	"fmt"
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
	fmt.Println("fileCfg", fileCfg)
	for method, endpoint := range fileCfg.Endpoints.Bearer {
		roles := make(map[string]struct{}, len(endpoint.Roles))
		for _, id := range endpoint.Roles {
			roles[id] = struct{}{}
		}
		cfg.BearerSet[method] = roles
	}
	return cfg, nil
}
