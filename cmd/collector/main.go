package main

import (
	"fmt"

	"rightdog/pkg/collector"
)

func main() {
	cfg := collector.Config{
		InfluxDB: collector.InfuxDBConfig{
			Reader: "http://127.0.0.1:8086",
			Writer: "http://127.0.0.1:8086",
			DB:     "ticker",
		},
		Coinone: collector.CoinoneConfig{
			Enabled:  true,
			Interval: "10s",
		},
	}

	runner, err := collector.NewRunner(&cfg)
	if err != nil {
		fmt.Printf("Collector cannot run; %+v\n", err)
		return
	}

	runner.Run()
}
