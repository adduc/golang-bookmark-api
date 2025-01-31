package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addMarketingRoutes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the Bookmark API")
	})
}
