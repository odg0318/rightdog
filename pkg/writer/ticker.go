package writer

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PostTickerParam struct {
	Time     int64   `json:"time" binding:"required"`
	Exchange string  `json:"exchange" binding:"required"`
	From     string  `json:"from" binding:"required"`
	To       string  `json:"to"`
	Price    float64 `json:"price" binding:"required"`
}

func (r *Rest) PostTicker(c *gin.Context) {
	var param PostTickerParam
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if param.Time <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "time must be larger than 0"})
		return
	}

	if len(param.To) == 0 {
		param.To = "krw"
	}

	influxClient, err := NewInfluxClient(r.cfg.InfluxDB.Writer, r.cfg.InfluxDB.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer influxClient.Close()

	tags := map[string]string{}
	tags["exchange"] = strings.ToLower(param.Exchange)
	tags["fromcurrency"] = strings.ToLower(param.From)
	tags["tocurrency"] = strings.ToLower(param.To)

	fields := map[string]interface{}{}
	fields["price"] = param.Price

	influxClient.AddPoint("ticker", tags, fields, time.Unix(param.Time, 0))

	err = influxClient.Write()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
