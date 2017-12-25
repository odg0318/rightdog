package collector

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/prometheus/client_golang/prometheus"
)

type coinoneTickerRaw struct {
	Volume    string `json:"volume"`
	Last      string `json:"last"`
	Timestamp string `json:"timestamp"`
	High      string `json:"high"`
	Result    string `json:"result"`
	ErrorCode string `json:"errorCode"`
	First     string `json:"first"`
	Low       string `json:"low"`
	Currency  string `json:"currency"`
}

type CoinoneCollector struct {
	cfg *CoinoneConfig
}

func (c *CoinoneCollector) Collect() error {
	for _, currency := range c.cfg.Currencies {
		err := c.collectTicker(currency)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CoinoneCollector) collectTicker(currency string) error {
	tickerUrl := fmt.Sprintf("https://api.coinone.co.kr/ticker?currency=%s", currency)
	res, err := http.Get(tickerUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var raw coinoneTickerRaw
	err = decoder.Decode(&raw)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", raw)

	return nil
}

func NewCoinoneCollector(cfg *CoinoneConfig) (Collector, error) {
	return &CoinoneCollector{cfg}, nil
}
