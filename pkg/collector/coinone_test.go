package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoinoneCollectTicker(t *testing.T) {
	c := CoinoneCollector{}

	raw, err := c.collectTicker()

	assert.Nil(t, err)
	assert.NotNil(t, raw)
}
