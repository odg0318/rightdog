package writer

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PostLatencyParam struct {
	Exchange string  `json:"exchange" binding:"required"`
	Latency  float64 `json:"latency"`
}

func (r *Rest) PostLatency(c *gin.Context) {
	var param PostLatencyParam
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// status ok returns when latency is zero.
	if param.Latency <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
		return
	}

	influxClient, err := NewInfluxClient(r.cfg.InfluxDB.Writer, r.cfg.InfluxDB.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer influxClient.Close()

	tags := map[string]string{}
	tags["exchange"] = strings.ToLower(param.Exchange)

	fields := map[string]interface{}{}
	fields["latency"] = param.Latency

	influxClient.AddPoint("latency", tags, fields, time.Now())

	err = influxClient.Write()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
