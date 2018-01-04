package collector

import (
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

type InfluxClient struct {
	c  client.Client
	bp client.BatchPoints
	db string
}

func (c *InfluxClient) Close() {
	c.c.Close()
}

func (c *InfluxClient) AddPoint(name string, tags map[string]string, fields map[string]interface{}, t ...time.Time) error {
	pt, err := client.NewPoint(name, tags, fields, t...)
	if err != nil {
		return err
	}
	c.bp.AddPoint(pt)

	return nil
}

func (c *InfluxClient) Write() error {
	if err := c.c.Write(c.bp); err != nil {
		return err
	}

	return nil
}

func (c *InfluxClient) Query(query string) ([]client.Result, error) {
	q := client.NewQuery(query, c.db, "ns")
	resp, err := c.c.Query(q)

	if err != nil {
		return nil, err
	}
	if resp.Error() != nil {
		return nil, resp.Error()
	}

	return resp.Results, nil
}

func NewInfluxClient(addr, db string) (*InfluxClient, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: addr,
	})
	if err != nil {
		return nil, err
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: db,
	})
	if err != nil {
		return nil, err
	}

	return &InfluxClient{
		c:  c,
		bp: bp,
		db: db,
	}, nil
}
