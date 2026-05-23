// Package manifest handles loading and validating service manifest files.
package manifest

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Manifest describes the desired state of a single deployed service.
type Manifest struct {
	Name     string            `yaml:"name"`
	Image    string            `yaml:"image"`
	Replicas int               `yaml:"replicas"`
	Labels   map[string]string `yaml:"labels"`
}

// Load reads a YAML manifest file from path and returns the parsed Manifest.
// It returns an error if the file cannot be read, cannot be decoded, or fails
// basic validation.
func Load(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("manifest: read %q: %w", path, err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("manifest: decode %q: %w", path, err)
	}

	if err := validate(m); err != nil {
		return Manifest{}, fmt.Errorf("manifest: validate %q: %w", path, err)
	}

	return m, nil
}

func validate(m Manifest) error {
	if m.Name == "" {
		return fmt.Errorf("name is required")
	}
	if m.Replicas < 0 {
		return fmt.Errorf("replicas must be non-negative, got %d", m.Replicas)
	}
	return nil
}
