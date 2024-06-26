package generator

import (
	"context"
	"fmt"

	"github.com/artarts36/filegen/internal/filesystem"
	"github.com/artarts36/filegen/internal/template"
)

type RenderingGenerator struct {
	renderer   template.Renderer
	filesystem filesystem.FileSystem
}

func NewRenderingGenerator(renderer template.Renderer, filesystem filesystem.FileSystem) *RenderingGenerator {
	return &RenderingGenerator{
		renderer:   renderer,
		filesystem: filesystem,
	}
}

func (g *RenderingGenerator) Generate(ctx context.Context, file GeneratingFile) error {
	vars := map[string]interface{}{
		"vars": map[string]interface{}{
			"local":  file.File.Vars,
			"global": file.GlobalVars,
		},
	}

	tmpl, err := g.filesystem.Get(file.File.TemplatePath)
	if err != nil {
		return fmt.Errorf("could not get template: %w", err)
	}

	content, err := g.renderer.Render(ctx, tmpl, vars)
	if err != nil {
		return fmt.Errorf("could not render template: %w", err)
	}

	outputPath, err := g.renderer.Render(ctx, []byte(file.File.OutputPath), vars)
	if err != nil {
		return fmt.Errorf("could not render output path: %w", err)
	}

	err = g.filesystem.Save(string(outputPath), content)
	if err != nil {
		return fmt.Errorf("could not save file: %w", err)
	}

	return nil
}
