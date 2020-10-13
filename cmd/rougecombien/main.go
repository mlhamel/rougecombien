package main

import (
	"context"
	"os"

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

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
