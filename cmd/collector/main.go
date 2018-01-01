package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"rightdog/pkg/collector"
)

func main() {
	app := cli.NewApp()
	app.Name = "Rightdog Collector"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			EnvVar: "COLLECTOR_CONFIG",
			Usage:  "Configuration file path",
			Value:  "../../collector.yml",
		},
	}
	app.Action = mainAction

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Collector cannot run; %+v\n", err)
		os.Exit(1)
	}
}

func mainAction(ctx *cli.Context) error {
	cfg, err := collector.NewConfigFromFile(ctx.String("config"))
	if err != nil {
		return err
	}

	err = cfg.Validate()
	if err != nil {
		return err
	}

	runner, err := collector.NewRunner(cfg)
	if err != nil {
		return err
	}

	runner.Run()

	return nil
}
