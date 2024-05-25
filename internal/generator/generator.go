package generator

import (
	"context"

	"github.com/artarts36/filegen/internal/config"
)

type GeneratingFile struct {
	File       config.File
	GlobalVars map[string]interface{}
}

type Generator interface {
	Generate(ctx context.Context, file GeneratingFile) error
}
