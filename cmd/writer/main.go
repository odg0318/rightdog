package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"rightdog/pkg/writer"
)

func main() {
	app := cli.NewApp()
	app.Name = "Rightdog Writer"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			EnvVar: "Writer_CONFIG",
			Usage:  "Configuration file path",
			Value:  "../../writer.yml",
		},
	}
	app.Action = mainAction

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Writer cannot run; %+v\n", err)
		os.Exit(1)
	}
}

func mainAction(ctx *cli.Context) error {
	cfg, err := writer.NewConfigFromFile(ctx.String("config"))
	if err != nil {
		return err
	}

	err = cfg.Validate()
	if err != nil {
		return err
	}

	rest := writer.NewRest(cfg)
	return rest.Run()
}
