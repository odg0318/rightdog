package writer

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	cfg *Config
}

func (r *Rest) Run() error {
	router := gin.Default()
	router.POST("/ticker", r.PostTicker)

	return router.Run(fmt.Sprintf(":%d", r.cfg.Rest.Port))
}

func NewRest(cfg *Config) *Rest {
	return &Rest{
		cfg: cfg,
	}
}