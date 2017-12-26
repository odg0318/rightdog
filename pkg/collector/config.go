package collector

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Coinone CoinoneConfig `yaml:"coinone"`
}

type CoinoneConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Interval string `yaml:"interval"`
}

func NewConfig(data []byte) (*Config, error) {
	var c Config
	err := yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func NewConfigFromFile(path string) (*Config, error) {
	return &Config{}, nil
}
