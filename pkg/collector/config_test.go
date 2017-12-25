package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoinoneConfig(t *testing.T) {
	rawConfig := `
coinone:
  enabled: true
  interval: 10s
  currencies:
   - btc
   - xrp
   - eth
`
	c, err := NewConfig([]byte(rawConfig))

	assert.Nil(t, err)

	assert.NotNil(t, c.Coinone)
	assert.True(t, c.Coinone.Enabled)
	assert.Equal(t, 3, len(c.Coinone.Currencies))
}
