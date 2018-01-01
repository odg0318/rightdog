package collector

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"rightdog/pkg/collector/influx"
)

type upbitTickerRaw struct {
	Code         string  `json:"code"`
	OpeningPrice float64 `json:"openingPrice"`
	HighPrice    float64 `json:"highPrice"`
	LowPrice     float64 `json:"lowPrice"`
	TradePrice   float64 `json:"tradePrice"`
}

func (r *upbitTickerRaw) GetRate() float64 {
	return r.TradePrice
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
	influxClient, err := influx.NewInfluxClient(c.cfg.InfluxDB.Writer, c.cfg.InfluxDB.DB)
	if err != nil {
		return err
	}

	defer influxClient.Close()

	for k, v := range c.cfg.Upbit.Currencies {
		tickerRaw, err := c.collectTicker(v)
		if err != nil {
			return err
		}

		c.addPoint(k, tickerRaw, influxClient)
	}

	err = influxClient.Write()
	if err != nil {
		return err
	}

	return nil
}

func (c *UpbitCollector) addPoint(currency string, v *upbitTickerRaw, influxClient *influx.InfluxClient) {
	tags := map[string]string{}
	tags["exchange"] = c.name
	tags["currency"] = currency

	fields := map[string]interface{}{}
	fields["rate"] = v.GetRate()

	influxClient.AddPoint("ticker", tags, fields, time.Now())
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
