package routes

import "github.com/gin-gonic/gin"

func addListRoutes(r *gin.RouterGroup) {
	r.GET("/me/lists", func(c *gin.Context) {
		// Handler logic for listing user lists
	})

	r.GET("/me/lists/:list_id", func(c *gin.Context) {
		// Handler logic for listing bookmarks in a specific list
	})

	r.GET("/lists", func(c *gin.Context) {
		// Handler logic for listing all lists
	})
}
