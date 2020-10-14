package main

import (
	"context"
	"fmt"
	"os"
	"path"

	"cirello.io/runner/procfile"
	"github.com/pior/runnable"
	"github.com/urfave/cli"

	"github.com/mlhamel/rougecombien/pkg/config"
	"github.com/mlhamel/rougecombien/pkg/web"
)

func main() {
	cfg := config.NewConfig()
	ctx := context.Background()

	app := cli.App{
		Name: "rougecombien",
		Action: func(*cli.Context) error {
			manager := runnable.Manager(nil)
			manager.Add(web.NewController(cfg))
			return manager.Build().Run(ctx)
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
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
