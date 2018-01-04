package collector

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type WriterClient struct {
	cfg WriterConfig
}

func (c *WriterClient) PostTicker(exchange string, currency string, price float64, t time.Time) error {
	url := fmt.Sprintf("%s/ticker", c.cfg.Addr)
	values := map[string]interface{}{
		"exchange": exchange,
		"from":     currency,
		"price":    price,
		"time":     t.Unix(),
	}

	marshaledJson, err := json.Marshal(values)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(marshaledJson))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("POST /ticker failed; %s", string(resBody)))
	}

	return nil
}

func NewWriterClient(cfg WriterConfig) *WriterClient {
	return &WriterClient{
		cfg: cfg,
	}
}
