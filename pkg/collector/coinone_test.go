package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoinoneCollectTicker(t *testing.T) {
	c := CoinoneCollector{}

	err := c.collectTicker("btc")

	assert.Nil(t, err)
}

func TestCoinoneCollect(t *testing.T) {
	cfg := Config{
		Coinone: CoinoneConfig{
			Currencies: []string{"btc", "xrp"},
		},
	}

	c, err := NewCoinoneCollector(&cfg.Coinone)
	assert.Nil(t, err)

	err = c.Collect()
	assert.Nil(t, err)
}
