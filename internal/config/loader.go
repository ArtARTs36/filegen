package config

import (
	"context"
	"errors"
	"fmt"
	"os"
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

	err = l.injectEnv(&config)
	if err != nil {
		return Config{}, err
	}

	return config, err
}

func (l *Loader) injectEnv(cfg *Config) error {
	vars, err := l.injectEnvToMap(cfg.Vars)
	if err != nil {
		return err
	}
	cfg.Vars = vars

	for i, file := range cfg.Files {
		fVars, fErr := l.injectEnvToMap(file.Vars)
		if fErr != nil {
			return fErr
		}

		cfg.Files[i].Vars = fVars
	}

	return nil
}

func (l *Loader) injectEnvToMap(vars map[string]interface{}) (map[string]interface{}, error) {
	transform := func(stringWithVar string) (string, error) {
		if stringWithVar == "" {
			return stringWithVar, nil
		}

		if stringWithVar[0] != '$' {
			return stringWithVar, nil
		}

		varName := stringWithVar[1:]
		varValue := os.Getenv(varName)
		if varValue == "" {
			return "", errors.New(varName + " is not set")
		}

		return varValue, nil
	}

	preparedVars := map[string]interface{}{}
	for key, val := range vars {
		switch v := val.(type) {
		case map[string]interface{}:
			pv, err := l.injectEnvToMap(v)
			if err != nil {
				return nil, err
			}

			preparedVars[key] = pv
		case string:
			pv, err := transform(v)
			if err != nil {
				return nil, err
			}
			preparedVars[key] = pv
		default:
			preparedVars[key] = v
		}
	}

	return preparedVars, nil
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
