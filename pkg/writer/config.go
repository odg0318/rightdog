package writer

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Rest     RestConfig     `yaml:"rest"`
	InfluxDB InfluxDBConfig `yaml:"influxdb"`

	validated bool
}

type RestConfig struct {
	Port int `yaml:"port"`
}

func (c *RestConfig) Validate() error {
	if c.Port == 0 {
		return errors.New("rest port is empty in config")
	}

	return nil
}

type InfluxDBConfig struct {
	Writer string `yaml:"writer"`
	Reader string `yaml:"reader"`
	DB     string `yaml:"db"`
}

func (c *InfluxDBConfig) Validate() error {
	if len(c.Writer) == 0 {
		return errors.New("influxdb writer is empty in config")
	}

	if len(c.Reader) == 0 {
		return errors.New("influxdb reader is empty in config")
	}

	if len(c.DB) == 0 {
		return errors.New("influxdb DB  is empty in config")
	}

	return nil
}

func (c *Config) Validate() error {
	if err := c.Rest.Validate(); err != nil {
		return err
	}

	if err := c.InfluxDB.Validate(); err != nil {
		return err
	}

	c.validated = true

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
