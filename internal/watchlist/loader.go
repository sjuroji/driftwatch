package watchlist

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// fileSchema mirrors the YAML structure for a watchlist file.
type fileSchema struct {
	Services []struct {
		Name   string            `yaml:"name"`
		Tags   []string          `yaml:"tags"`
		Labels map[string]string `yaml:"labels"`
	} `yaml:"services"`
}

// LoadFile reads a YAML watchlist file and populates a Watchlist.
// The file must contain a top-level "services" list.
//
// Example YAML:
//
//	services:
//	  - name: auth-service
//	    tags: [prod, critical]
//	    labels:
//	      env: production
func LoadFile(path string) (*Watchlist, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("watchlist: read file: %w", err)
	}

	var schema fileSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("watchlist: parse yaml: %w", err)
	}

	w := New()
	for _, svc := range schema.Services {
		e := Entry{
			Name:   svc.Name,
			Tags:   svc.Tags,
			Labels: svc.Labels,
		}
		if err := w.Add(e); err != nil {
			return nil, fmt.Errorf("watchlist: invalid entry: %w", err)
		}
	}
	return w, nil
}
