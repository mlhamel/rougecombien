package main

import (
	"context"
	"fmt"
	"os"
	"path"

	"cirello.io/runner/procfile"
	"github.com/urfave/cli"

	"github.com/mlhamel/rougecombien/pkg/config"
	"github.com/mlhamel/rougecombien/pkg/consumer"
	"github.com/mlhamel/rougecombien/pkg/scraper"
)

func main() {
	cfg := config.NewConfig()
	ctx := context.Background()

	app := cli.App{
		Name: "rougecombien",
		Action: func(*cli.Context) error {
			if err := scraper.NewScraper(cfg).Run(ctx); err != nil {
				panic(err)
			}

			return nil
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "dev",
			Action: func(c *cli.Context) error {
				profile := path.Join(cfg.ProjectPath(), "Procfile")
				file, err := os.Open(profile)
				if err != nil {
					return fmt.Errorf("Cannot find Procfile at %s: %w", profile, err)
				}
				runner, err := procfile.Parse(file)
				if err != nil {
					return err
				}
				cfg.Logger().Info().Str("workdir", runner.WorkDir).Str("profile", profile).Msg(fmt.Sprintf("Running"))
				runner.BasePort = cfg.HTTPPort()
				return runner.Start(ctx)
			},
		},
		{
			Name: "subscribe",
			Action: func(c *cli.Context) error {
				if err := consumer.New(cfg).Run(ctx); err != nil {
					panic(err)
				}
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
