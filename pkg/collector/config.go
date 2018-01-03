package collector

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	InfluxDB InfluxDBConfig `yaml:"influxdb"`
	Coinone  CoinoneConfig  `yaml:"coinone"`
	Korbit   KorbitConfig   `yaml:"korbit"`
	Upbit    UpbitConfig    `yaml:"upbit"`

	validated bool
}

func (c *Config) Validate() error {
	if err := c.Coinone.Validate(); err != nil {
		return err
	}

	if err := c.Korbit.Validate(); err != nil {
		return err
	}

	if err := c.Upbit.Validate(); err != nil {
		return err
	}

	c.validated = true

	return nil
}

type InfluxDBConfig struct {
	Writer string `yaml:"writer"`
	Reader string `yaml:"reader"`
	DB     string `yaml:"db"`
}

type CoinoneConfig struct {
	Enabled     bool   `yaml:"enabled"`
	RawInterval string `yaml:"interval"`

	Interval time.Duration `yaml:"_,omitempty"`
}

func (c *CoinoneConfig) Validate() error {
	if c.Enabled == false {
		return nil
	}

	var err error
	c.Interval, err = time.ParseDuration(c.RawInterval)
	if err != nil {
		return err
	}

	return nil
}

type KorbitConfig struct {
	Enabled     bool              `yaml:"enabled"`
	RawInterval string            `yaml:"interval"`
	Auth        KorbitAuthConfig  `yaml:"auth"`
	Currencies  map[string]string `yaml:"currencies"`

	Interval time.Duration `yaml:"_,omitempty"`
}

func (c *KorbitConfig) Validate() error {
	if c.Enabled == false {
		return nil
	}

	var err error
	c.Interval, err = time.ParseDuration(c.RawInterval)
	if err != nil {
		return err
	}

	return nil
}

type KorbitAuthConfig struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

type UpbitConfig struct {
	Enabled     bool              `yaml:"enabled"`
	RawInterval string            `yaml:"interval"`
	Currencies  map[string]string `yaml:"currencies"`

	Interval time.Duration `yaml:"_,omitempty"`
}

func (c *UpbitConfig) Validate() error {
	if c.Enabled == false {
		return nil
	}

	var err error
	c.Interval, err = time.ParseDuration(c.RawInterval)
	if err != nil {
		return err
	}

	return nil
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
