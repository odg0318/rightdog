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
`
	c, err := NewConfig([]byte(rawConfig))

	assert.Nil(t, err)

	assert.NotNil(t, c.Coinone)
	assert.True(t, c.Coinone.Enabled)
}

func TestKorbitConfig(t *testing.T) {
	rawConfig := `
korbit:
  enabled: true
  interval: 10s
`
	c, err := NewConfig([]byte(rawConfig))

	assert.Nil(t, err)

	assert.NotNil(t, c.Korbit)
	assert.True(t, c.Korbit.Enabled)
	assert.Equal(t, 3, len(c.Korbit.Currencies))
}
