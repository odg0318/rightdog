package collector

import (
	"time"
)

type Collector interface {
	Collect() error
	Run()
}

type Runner struct {
	collectors []Collector
}

func (r *Runner) Run() {
	for _, collector := range r.collectors {
		go collector.Run()
	}

	for true {
		time.Sleep(time.Second)
	}
}

func NewRunner(cfg *Config) (*Runner, error) {
	runner := Runner{}

	if cfg.Upbit.Enabled == true {
		c, err := NewUpbitCollector(cfg)
		if err != nil {
			return nil, err
		}

		runner.collectors = append(runner.collectors, c)
	}

	return &runner, nil
}
