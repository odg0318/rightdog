package collector

import (
	"errors"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Writer WriterConfig `yaml:"writer"`

	Upbit UpbitConfig `yaml:"upbit"`

	validated bool
}

func (c *Config) Validate() error {
	if err := c.Writer.Validate(); err != nil {
		return err
	}

	if err := c.Upbit.Validate(); err != nil {
		return err
	}

	c.validated = true

	return nil
}

type WriterConfig struct {
	Addr string `yaml:"addr"`
}

func (c *WriterConfig) Validate() error {
	if len(c.Addr) == 0 {
		return errors.New("writer addr is empty in config")
	}

	return nil
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
