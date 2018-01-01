package collector

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"rightdog/pkg/collector/influx"
)

type coinoneTickerRaw struct {
	BCH  coinoneTickerValueRaw `json:"bch"`
	QTUM coinoneTickerValueRaw `json:"qtum"`
	IOTA coinoneTickerValueRaw `json:"iota"`
	LTC  coinoneTickerValueRaw `json:"ltc"`
	ETC  coinoneTickerValueRaw `json:"etc"`
	BTG  coinoneTickerValueRaw `json:"btg"`
	BTC  coinoneTickerValueRaw `json:"btc"`
	ETH  coinoneTickerValueRaw `json:"eth"`
	XRP  coinoneTickerValueRaw `json:"xrp"`
}

type coinoneTickerValueRaw struct {
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

func (r coinoneTickerValueRaw) GetRate() float64 {
	f, _ := strconv.ParseFloat(r.Last, 64)
	return f
}

type CoinoneCollector struct {
	cfg    *Config
	logger *log.Logger
	name   string
}

func (c *CoinoneCollector) Run() {
	for true {
		err := c.Collect()
		if err != nil {
			c.logger.Printf("%+v\n", err)
		}

		c.logger.Printf("updated.\n")
		time.Sleep(c.cfg.Coinone.Interval)
	}
}

func (c *CoinoneCollector) Collect() error {
	influxClient, err := influx.NewInfluxClient(c.cfg.InfluxDB.Writer, c.cfg.InfluxDB.DB)
	if err != nil {
		return err
	}

	defer influxClient.Close()

	tickerRaw, err := c.collectTicker()
	if err != nil {
		return err
	}

	c.addPoint("bch", tickerRaw.BCH, influxClient)
	c.addPoint("qtum", tickerRaw.QTUM, influxClient)
	c.addPoint("iota", tickerRaw.IOTA, influxClient)
	c.addPoint("ltc", tickerRaw.LTC, influxClient)
	c.addPoint("etc", tickerRaw.ETC, influxClient)
	c.addPoint("btg", tickerRaw.BTG, influxClient)
	c.addPoint("btc", tickerRaw.BTC, influxClient)
	c.addPoint("eth", tickerRaw.ETH, influxClient)
	c.addPoint("xrp", tickerRaw.XRP, influxClient)

	err = influxClient.Write()
	if err != nil {
		return err
	}

	return nil
}

func (c *CoinoneCollector) addPoint(currency string, v coinoneTickerValueRaw, influxClient *influx.InfluxClient) {
	tags := map[string]string{}
	tags["exchange"] = c.name
	tags["currency"] = currency

	fields := map[string]interface{}{}
	fields["rate"] = v.GetRate()

	influxClient.AddPoint("ticker", tags, fields, time.Now())
}

func (c *CoinoneCollector) collectTicker() (*coinoneTickerRaw, error) {
	tickerUrl := "https://api.coinone.co.kr/ticker?currency=all"
	res, err := http.Get(tickerUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var raw coinoneTickerRaw
	err = decoder.Decode(&raw)
	if err != nil {
		return nil, err
	}

	return &raw, nil
}

func NewCoinoneCollector(cfg *Config) (Collector, error) {
	name := "coinone"

	return &CoinoneCollector{
		cfg:    cfg,
		logger: log.New(os.Stdout, fmt.Sprintf("[%s] ", name), log.LstdFlags),
		name:   name,
	}, nil
}
