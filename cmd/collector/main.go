package main

import (
	"fmt"

	"rightdog/pkg/collector"
)

func main() {
	cfg, err := collector.NewConfigFromFile("../../collector.yml")
	if err != nil {
		fmt.Printf("Collector cannot run; %+v\n", err)
		return
	}

	err = cfg.Validate()
	if err != nil {
		fmt.Printf("Collector cannot run; %+v\n", err)
		return
	}
	runner, err := collector.NewRunner(cfg)
	if err != nil {
		fmt.Printf("Collector cannot run; %+v\n", err)
		return
	}

	runner.Run()
}
