package main

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/confmap"
	"gopkg.in/yaml.v3"
)

const schemeName = "test"

// Simple test provider
// It takes the config as a string and allows it to be retreived as raw yaml
// Its schema begins with `test:<anything-can-go-here>`
type provider struct {
	yamlString string
}

func NewTestProvider(yamlString string) confmap.Provider {
	return &provider{
		yamlString: yamlString,
	}
}

func (s *provider) Retrieve(_ context.Context, uri string, _ confmap.WatcherFunc) (*confmap.Retrieved, error) {
	if !strings.HasPrefix(uri, schemeName+":") {
		return nil, fmt.Errorf("%q uri is not supported by %q provider", uri, schemeName)
	}

	var rawConf any
	if err := yaml.Unmarshal([]byte(s.yamlString), &rawConf); err != nil {
		return nil, err
	}
	return confmap.NewRetrieved(rawConf)
}

func (*provider) Scheme() string {
	return schemeName
}

func (s *provider) Shutdown(context.Context) error {
	return nil
}
