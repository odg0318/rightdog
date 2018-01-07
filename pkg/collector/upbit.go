package collector

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

//[{"code":"CRIX.UPBIT.KRW-BTC","candleDateTime":"2018-01-03T19:59:00+00:00","candleDateTimeKst":"2018-01-04T04:59:00+09:00","openingPrice":20349000.00000000,"highPrice":20350000.00000000,"lowPrice":20320000.00000000,"tradePrice":20331000.00000000,"candleAccTradeVolume":10.57938696,"candleAccTradePrice":215157469.67828000,"timestamp":1515009599203,"unit":1} ]
type upbitTickerRaw struct {
	Code         string  `json:"code"`
	OpeningPrice float64 `json:"openingPrice"`
	HighPrice    float64 `json:"highPrice"`
	LowPrice     float64 `json:"lowPrice"`
	TradePrice   float64 `json:"tradePrice"`
	Timestamp    int64   `json:"timestamp"`
}

func (r *upbitTickerRaw) GetPrice() float64 {
	return r.TradePrice
}

func (r *upbitTickerRaw) GetTime() time.Time {
	return time.Unix(r.Timestamp/1000, 0)
}

type UpbitCollector struct {
	cfg    *Config
	logger *log.Logger
	name   string
}

func (c *UpbitCollector) Run() {
	for true {
		err := c.Collect()
		if err != nil {
			c.logger.Printf("%+v\n", err)
		}

		c.logger.Printf("updated.\n")

		time.Sleep(c.cfg.Upbit.Interval)
	}
}

func (c *UpbitCollector) Collect() error {
	writer := NewWriterClient(c.cfg.Writer)

	for currency, v := range c.cfg.Upbit.Currencies {
		now := time.Now()
		tickerRaw, err := c.collectTicker(v)
		if err != nil {
			c.logger.Printf("collecting failed; %+v", err)
			continue
		}

		latency := time.Since(now).Seconds()
		if latency > 0 {
			err = writer.PostLatency(c.name, latency)
			if err != nil {
				c.logger.Printf("writing failed; %+v", err)
			}
		}

		err = writer.PostTicker(c.name, currency, tickerRaw.GetPrice(), tickerRaw.GetTime())
		if err != nil {
			c.logger.Printf("writing failed; %+v", err)
			continue
		}
	}

	return nil
}

func (c *UpbitCollector) collectTicker(currency string) (*upbitTickerRaw, error) {
	now := time.Now().UTC().Add(time.Hour * 9)

	apiUrl := fmt.Sprintf("https://crix-api-endpoint.upbit.com/v1/crix/candles/minutes/1?code=CRIX.UPBIT.%s&count=1&to=%s%%20%s", currency, now.Format("2006-01-02"), now.Format("15:04:05"))

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var raw []upbitTickerRaw
	err = decoder.Decode(&raw)
	if err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		return nil, errors.New("data length is zero")
	}

	return &raw[0], nil
}

func NewUpbitCollector(cfg *Config) (Collector, error) {
	name := "upbit"

	return &UpbitCollector{
		cfg:    cfg,
		logger: log.New(os.Stdout, fmt.Sprintf("[%s] ", name), log.LstdFlags),
		name:   name,
	}, nil
}
