package routes

import "github.com/gin-gonic/gin"

func addTagRoutes(r *gin.RouterGroup) {
	r.GET("/me/tags", func(c *gin.Context) {
		// Handler logic for listing user tags
	})

	r.GET("/tags", func(c *gin.Context) {
		// Handler logic for listing all tags
	})
}
