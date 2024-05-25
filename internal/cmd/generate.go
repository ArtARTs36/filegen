package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/artarts36/filegen/internal/config"
	"github.com/artarts36/filegen/internal/generator"
)

type Generate struct {
	configLoader *config.Loader
	generator    generator.Generator
}

func NewGenerate(configLoader *config.Loader, generator generator.Generator) *Generate {
	return &Generate{
		configLoader: configLoader,
		generator:    generator,
	}
}

type GenerateParams struct {
	ConfigPath string
}

func (cmd *Generate) Execute(ctx context.Context, params GenerateParams) error {
	cfg, err := cmd.configLoader.Load(ctx, params.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	wg := &sync.WaitGroup{}

	for _, file := range cfg.Files {
		wg.Add(1)

		file := file

		go func() {
			defer wg.Done()

			slog.DebugContext(ctx, fmt.Sprintf("[cmd] generating %q", file.OutputPath))

			genErr := cmd.generator.Generate(ctx, generator.GeneratingFile{
				File:       file,
				GlobalVars: cfg.Vars,
			})
			if genErr != nil {
				slog.
					With(slog.Any("err", genErr)).
					ErrorContext(ctx, fmt.Sprintf("failed to generate file %q", file.OutputPath))
			}
		}()
	}

	wg.Wait()

	return nil
}
