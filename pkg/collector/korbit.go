package collector

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type korbitAccessTokenRaw struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type korbitTickerRaw struct {
	Timestamp int    `json:"timestamp"`
	Last      string `json:"last"`
}

func (r korbitTickerRaw) GetRate() float64 {
	f, _ := strconv.ParseFloat(r.Last, 64)
	return f
}

type KorbitCollector struct {
	cfg    *Config
	logger *log.Logger
	name   string
}

func (c *KorbitCollector) Run() {
	for true {
		err := c.Collect()
		if err != nil {
			c.logger.Printf("%+v\n", err)
		}
		time.Sleep(c.cfg.Korbit.GetInterval())
	}
}

func (c *KorbitCollector) Collect() error {
	for _, currency := range c.cfg.Korbit.Currencies {
		tickerRaw, err := c.collectTicker(currency)
		if err != nil {
			return err
		}

		fmt.Printf("%s %+v\n", currency, tickerRaw)
	}

	return nil
}

func (c *KorbitCollector) getAccessToken() (*korbitAccessTokenRaw, error) {
	apiUrl := "https://api.korbit.co.kr/v1/oauth2/access_token"

	v := url.Values{}
	v.Set("client_id", c.cfg.Korbit.Auth.ClientId)
	v.Set("client_secret", c.cfg.Korbit.Auth.ClientSecret)
	v.Set("username", c.cfg.Korbit.Auth.Username)
	v.Set("password", c.cfg.Korbit.Auth.Password)
	v.Set("grant_type", "password")

	res, err := http.PostForm(apiUrl, v)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var raw korbitAccessTokenRaw
	err = decoder.Decode(&raw)
	if err != nil {
		return nil, err
	}

	return &raw, nil
}

func (c *KorbitCollector) collectTicker(currency string) (*korbitTickerRaw, error) {
	accessToken, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("https://api.korbit.co.kr/v1/ticker?currency_pair=%s", currency)

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", accessToken.TokenType, accessToken.AccessToken))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var raw korbitTickerRaw
	err = decoder.Decode(&raw)
	if err != nil {
		return nil, err
	}

	return &raw, nil
}

func NewKorbitCollector(cfg *Config) (Collector, error) {
	name := "korbit"

	return &KorbitCollector{
		cfg:    cfg,
		logger: log.New(os.Stdout, fmt.Sprintf("[%s] ", name), log.LstdFlags),
		name:   name,
	}, nil
}
