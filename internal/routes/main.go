package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	group := router.Group("/")
	addMarketingRoutes(group)
	addBookmarkRoutes(group, db)
	addListRoutes(group)
	addTagRoutes(group)

	return router
}
