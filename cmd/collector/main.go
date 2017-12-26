package main

import (
	"fmt"

	"rightdog/pkg/collector"
)

func main() {
	cfg := collector.CoinoneConfig{
		Enabled:  true,
		Interval: "10s",
	}
	collector, err := collector.NewCoinoneCollector(&cfg)
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}

	err = collector.Collect()
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
}
