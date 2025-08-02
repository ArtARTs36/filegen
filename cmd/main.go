package main

import (
	"context"
	"fmt"
	"github.com/artarts36/filegen/internal/cmd"
	"github.com/artarts36/filegen/internal/config"
	"github.com/artarts36/filegen/internal/filesystem"
	"github.com/artarts36/filegen/internal/generator"
	"github.com/artarts36/filegen/internal/template"
	cli "github.com/artarts36/singlecli"
	"github.com/tyler-sommer/stick"
	"os"
)

var (
	Version   = "0.1.0"
	BuildDate = "2024-03-09 03:09:15"
)

func main() {
	application := cli.App{
		BuildInfo: &cli.BuildInfo{
			Name:      "filegen",
			Version:   Version,
			BuildDate: BuildDate,
		},
		Action: run,
		Args: []*cli.ArgDefinition{
			{
				Name:        "config-path",
				Description: "Path to config file",
			},
		},
	}

	application.RunWithGlobalArgs(context.Background())
}

func run(ctx *cli.Context) error {
	fs := filesystem.NewLocal()

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	command := cmd.NewGenerate(
		config.CreateLoader(fs),
		generator.NewRenderingGenerator(template.NewStickRenderer(stick.NewFilesystemLoader(wd)), fs),
	)

	return command.Execute(ctx.Context, cmd.GenerateParams{
		ConfigPath: ctx.GetArg("config-path"),
	})
}
