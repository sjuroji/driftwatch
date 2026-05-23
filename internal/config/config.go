// Package config handles loading and validating driftwatch runtime configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds the runtime configuration for driftwatch.
type Config struct {
	ManifestDir string   `yaml:"manifest_dir"`
	OutputFormat string   `yaml:"output_format"`
	Services    []string `yaml:"services"`
	FailOnDrift bool     `yaml:"fail_on_drift"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		ManifestDir:  "./manifests",
		OutputFormat: "text",
		Services:     []string{},
		FailOnDrift:  false,
	}
}

// Load reads a YAML config file from path and merges it with defaults.
// If path is empty, the default config is returned.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file %q: %w", path, err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// validate checks that Config fields contain acceptable values.
func (c Config) validate() error {
	valid := map[string]bool{"text": true, "json": true}
	fmt := strings.ToLower(c.OutputFormat)
	if !valid[fmt] {
		return errors.New("output_format must be \"text\" or \"json\"")
	}
	if c.ManifestDir == "" {
		return errors.New("manifest_dir must not be empty")
	}
	return nil
}
