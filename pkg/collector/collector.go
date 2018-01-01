package collector

import (
	"time"
)

type Collector interface {
	Collect() error
}

type Runner struct {
	collectors []Collector
}

func (r *Runner) Run() {
	for _, collector := range r.collectors {
		go collector.Collect()
	}

	for true {
		time.Sleep(time.Second)
	}
}

func NewRunner(cfg *Config) (*Runner, error) {
	runner := Runner{}

	if cfg.Coinone.Enabled == true {
		c, err := NewCoinoneCollector(cfg)
		if err != nil {
			return nil, err
		}

		runner.collectors = append(runner.collectors, c)
	}

	return &runner, nil
}
