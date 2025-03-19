package httpinterface

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (chs *CentaureissiHttpInterface) InitializeStatsRoutes() {
	protected := chs.g.Group("/api/stats").Use(chs.AuthMiddleware())
	{
		protected.GET("/", chs.stats)
	}
}

func (chs *CentaureissiHttpInterface) stats(c *gin.Context) {
	result := chs.services.Stats()

	c.JSON(http.StatusOK, result)
}
