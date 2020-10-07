package main

import (
	"context"
	"os"

	"github.com/pior/runnable"
	"github.com/urfave/cli"
)

func main() {
	cliApp := cli.NewApp()
	ctx := context.Background()

	cliApp.Commands = []cli.Command{
		{
			Name: "run",
			Action: func(*cli.Context) error {
				return runnable.Manager(nil).Build().Run(ctx)
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
