package collector

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	InfluxDB InfuxDBConfig `yaml:"influxdb"`
	Coinone  CoinoneConfig `yaml:"coinone"`
	Korbit   KorbitConfig  `yaml:"korbit"`

	validated bool
}

func (c *Config) Validate() error {
	var err error
	if c.Coinone.Enabled {
		c.Coinone.interval, err = time.ParseDuration(c.Coinone.Interval)
		if err != nil {
			return err
		}
	}

	if c.Korbit.Enabled {
		c.Korbit.interval, err = time.ParseDuration(c.Korbit.Interval)
		if err != nil {
			return err
		}
	}

	c.validated = true

	return nil
}

type InfuxDBConfig struct {
	Writer string `yaml:"writer"`
	Reader string `yaml:"reader"`
	DB     string `yaml:"db"`
}

type CoinoneConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Interval string `yaml:"interval"`

	interval time.Duration
}

func (c CoinoneConfig) GetInterval() time.Duration {
	return c.interval
}

type KorbitConfig struct {
	Enabled    bool             `yaml:"enabled"`
	Interval   string           `yaml:"interval"`
	Auth       KorbitAuthConfig `yaml:"auth"`
	Currencies []string         `yaml:"currencies"`

	interval time.Duration
}

func (c KorbitConfig) GetInterval() time.Duration {
	return c.interval
}

type KorbitAuthConfig struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
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
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewConfig(data)
}
