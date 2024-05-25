package config

import (
	"context"

	"gopkg.in/yaml.v3"
)

type YAMLParser struct {
}

func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

func (y *YAMLParser) Parse(_ context.Context, content []byte) (Config, error) {
	var config Config

	err := yaml.Unmarshal(content, &config)

	return config, err
}
