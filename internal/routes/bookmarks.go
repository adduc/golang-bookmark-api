package routes

import (
	"log"
	"net/http"
	"net/url"

	"github.com/adduc/exercise-golang-bookmark-db/internal"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func addBookmarkRoutes(r *gin.RouterGroup, db *gorm.DB) {

	r.GET("/me/bookmarks", func(c *gin.Context) {
		// Handler logic for listing user bookmarks
	})

	r.POST("/me/bookmarks", func(c *gin.Context) {

		// CORS: Allow POST requests to this endpoint
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
		// CORS: we intentionally want to allow requests to this endpoint
		// from any origin to allow our bookmarklet for saving bookmarks
		// to work on any website
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Parse and validate the request
		var json struct {
			URL  string `json:"url" binding:"required"`
			Note string `json:"note"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find or create the bookmark
		bookmark, err := findOrCreateBookmark(db, json.URL)
		if err != nil && err.(*internal.ValidationError) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			log.Printf("Failed to find or create bookmark: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bookmark"})
			return
		}

		///
		// Find and update the user bookmark if it exists, otherwise create it
		// Then return the created bookmark
		///
		// Assuming user ID is 1 for now, this should be replaced with actual user ID from authentication
		userID := uint(1)

		var userBookmark internal.UserBookmark
		err = db.Limit(1).Find(&userBookmark, "user_id = ? AND bookmark_id = ?", userID, bookmark.ID).Error

		if err != nil {
			log.Printf("Failed to query user bookmark: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user bookmark"})
			return
		}

		if userBookmark.ID == 0 {
			// user bookmark not found
			userBookmark = internal.UserBookmark{
				UserID:     userID,
				BookmarkID: bookmark.ID,
				Note:       json.Note,
			}
			if err := db.Create(&userBookmark).Error; err != nil {
				log.Printf("Failed to create user bookmark: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user bookmark"})
				return
			}
		} else {
			// existing user bookmark found
			if json.Note != "" && userBookmark.Note != json.Note {
				userBookmark.Note = json.Note
				if err := db.Save(&userBookmark).Error; err != nil {
					log.Printf("Failed to update user bookmark: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user bookmark"})
					return
				}
			}
		}

		response := internal.UserBookmarkResponse{
			BookmarkID:     bookmark.ID,
			UserBookmarkID: userBookmark.ID,
			UserID:         userBookmark.UserID,
			URL:            bookmark.URL,
			Title:          bookmark.Title,
			Description:    bookmark.Description,
			Note:           userBookmark.Note,
			CreatedAt:      userBookmark.CreatedAt,
			UpdatedAt:      userBookmark.UpdatedAt,
		}

		c.JSON(http.StatusOK, response)
	})

	r.GET("/bookmarks", func(c *gin.Context) {
		// Handler logic for listing all bookmarks
	})
}

func findOrCreateBookmark(db *gorm.DB, bookmarkURL string) (*internal.Bookmark, error) {
	// Check if the URL is valid
	parsedURL, err := url.ParseRequestURI(bookmarkURL)

	if err != nil || parsedURL.Hostname() == "" {
		return nil, internal.NewValidationError("Invalid URL")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, internal.NewValidationError("Unsupported URL scheme")
	}

	// Check if the bookmark already exists
	var bookmark internal.Bookmark
	err = db.Limit(1).Find(&bookmark, "url = ?", bookmarkURL).Error

	if err != nil {
		return nil, err
	}

	if bookmark.ID != 0 {
		return &bookmark, nil
	}

	// Create the bookmark
	bookmark = internal.Bookmark{
		URL: bookmarkURL,
	}

	if err := db.Create(&bookmark).Error; err != nil {
		return nil, err
	}

	return &bookmark, nil
}
