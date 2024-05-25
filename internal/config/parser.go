package config

import "context"

type Parser interface {
	Parse(ctx context.Context, content []byte) (Config, error)
}
