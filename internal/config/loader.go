package config

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

type Loader struct {
	storage configStorage
	parsers map[string]Parser
}

type configStorage interface {
	Get(path string) ([]byte, error)
}

func (l *Loader) Load(ctx context.Context, path string) (Config, error) {
	configContent, err := l.storage.Get(path)
	if err != nil {
		return Config{}, fmt.Errorf("error loading config: %w", err)
	}

	parser, err := l.getParser(path)
	if err != nil {
		return Config{}, err
	}

	config, err := parser.Parse(ctx, configContent)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config: %w", err)
	}

	return config, err
}

func (l *Loader) getParser(path string) (Parser, error) {
	const minExtLength = 3

	ext := filepath.Ext(path)
	if len(ext) < minExtLength {
		return nil, fmt.Errorf("invalid config file extension")
	}

	ext = strings.ToLower(ext[1:])

	parser, ok := l.parsers[ext]
	if !ok {
		return nil, fmt.Errorf("parser not found for extension: %s", ext)
	}

	return parser, nil
}
