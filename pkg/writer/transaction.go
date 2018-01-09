package writer

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PostTransactionParam struct {
	Currency string `json:"currency" binding:"required"`
	Remain   int64  `json:"remain"`
}

func (r *Rest) PostTransaction(c *gin.Context) {
	var param PostTransactionParam
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if param.Remain < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "remain must be more than zero"})
		return
	}

	influxClient, err := NewInfluxClient(r.cfg.InfluxDB.Writer, r.cfg.InfluxDB.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer influxClient.Close()

	tags := map[string]string{}
	tags["currency"] = strings.ToLower(param.Currency)

	fields := map[string]interface{}{}
	fields["remain"] = param.Remain

	influxClient.AddPoint("transaction", tags, fields, time.Now())

	err = influxClient.Write()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
