package main

import (
	"context"
	"os"

	"github.com/mlhamel/rougecombien/pkg/config"

	"github.com/pior/runnable"
	"github.com/urfave/cli"

	"github.com/mlhamel/rougecombien/pkg/web"
)

func main() {
	cliApp := cli.NewApp()
	cfg := config.NewConfig()
	ctx := context.Background()

	cliApp.Commands = []cli.Command{
		{
			Name: "run",
			Action: func(*cli.Context) error {
				manager := runnable.Manager(nil)
				manager.Add(web.NewController(cfg))
				return manager.Build().Run(ctx)
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}
