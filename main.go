package main

/*

Synopsis:
- A bookmark database (in the vein of Del.icio.us or Pocket)

Goal:
- Create a simple web server that listens on port 8080

Entities:
- "Shared" entities used by mutiple users:
  - Bookmark: The URL, and possibly the title and description if we can
    get a scraper in place
  - Tag: A keyword associated with a bookmark
  - Auth Method: A method of authentication (e.g. email/password)
- User-specific entities:
  - User: A user of the system
  - UserAuth: One of the user's authentication methods (e.g. email/password)
  - UserBookmark: Bookmark metadata associated with a user (e.g. note, tags)
  - List: A collection of bookmarks created by a user
  - BookmarkTag: A tag associated with a bookmark by a user

Database (SQLite3):

- shared tables
   - auth_methods: id, name
   - bookmarks: id, url, title, description
   - tags: id, name
- user tables
   - users: id, username
   - user_auths: id, user_id, method, value
   - user_bookmarks: id, user_id, bookmark_id, note
   - lists: id, user_id, name
   - list_bookmarks: id, list_id, bookmark_id
   - user_bookmark_tags: id, bookmark_id, tag_id, user_id

Routes:
- GET /: Landing Page (Welcome to the Bookmark API)
- GET /me/bookmarks: List of Bookmarks for the user
- POST /me/bookmarks: Add a new bookmark
- GET /me/lists: List of Lists for the user
- GET /me/lists/:list_id: List of Bookmarks for a specific list
- GET /me/tags: List of Tags for the user
- GET /lists: List of all lists
- GET /bookmarks: List of all bookmarks
- GET /tags: List of all tags

Opportunities for Expansion:
- Add a scraper to get the title and description of a URL
- Add a search endpoint
- Social features like following lists or users
- Gamification features like badges or points to encourage usage

Constraints:
- gin should be used for the web server
- gorm should be used for the ORM
- sqlite3 should be used for the database

*/

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Bookmark struct {
	gorm.Model
	ID  uint   `gorm:"primaryKey"`
	URL string `gorm:"unique"`
	// title and description should be optional, as we may not always have them
	Title       *string `gorm:"default:null"`
	Description *string `gorm:"default:null"`
}

type Tag struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type AuthMethod struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
}

type UserAuth struct {
	gorm.Model
	ID     uint `gorm:"primaryKey"`
	UserID uint
	Method string
	Value  string

	User User
}

type UserBookmark struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	BookmarkID uint
	Note       string

	User     *User
	Bookmark *Bookmark
}

type List struct {
	gorm.Model
	ID     uint `gorm:"primaryKey"`
	UserID uint
	Name   string

	User User
}

type ListBookmark struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	ListID     uint
	BookmarkID uint

	List     List
	Bookmark Bookmark
}

type UserBookmarkTag struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	BookmarkID uint
	TagID      uint

	User     User
	Bookmark Bookmark
	Tag      Tag
}

type UserBookmarkResponse struct {
	BookmarkID     uint      `json:"bookmark_id"`
	UserBookmarkID uint      `json:"user_bookmark_id"`
	UserID         uint      `json:"user_id"`
	URL            string    `json:"url"`
	Title          *string   `json:"title"`
	Description    *string   `json:"description"`
	Note           string    `json:"note"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Errors

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func findOrCreateBookmark(db *gorm.DB, bookmarkURL string) (*Bookmark, error) {
	// Check if the URL is valid
	parsedURL, err := url.ParseRequestURI(bookmarkURL)

	if err != nil || parsedURL.Hostname() == "" {
		return nil, &ValidationError{"Invalid URL"}
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, &ValidationError{"Unsupported URL scheme"}
	}

	// Check if the bookmark already exists
	var bookmark Bookmark
	err = db.Limit(1).Find(&bookmark, "url = ?", bookmarkURL).Error

	if err != nil {
		return nil, err
	}

	if bookmark.ID != 0 {
		return &bookmark, nil
	}

	// Create the bookmark
	bookmark = Bookmark{
		URL: bookmarkURL,
	}

	if err := db.Create(&bookmark).Error; err != nil {
		return nil, err
	}

	return &bookmark, nil
}

func main() {
	r := gin.Default()

	db, err := gorm.Open(sqlite.Open("bookmarks.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Bookmark{}, &Tag{}, &AuthMethod{}, &User{}, &UserAuth{}, &UserBookmark{}, &List{}, &ListBookmark{}, &UserBookmarkTag{})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the Bookmark API")
	})

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
		if err != nil && err.(*ValidationError) != nil {
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

		var userBookmark UserBookmark
		err = db.Limit(1).Find(&userBookmark, "user_id = ? AND bookmark_id = ?", userID, bookmark.ID).Error

		if err != nil {
			log.Printf("Failed to query user bookmark: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user bookmark"})
			return
		}

		if userBookmark.ID == 0 {
			// user bookmark not found
			userBookmark = UserBookmark{
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

		response := UserBookmarkResponse{
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

	r.GET("/me/lists", func(c *gin.Context) {
		// Handler logic for listing user lists
	})

	r.GET("/me/lists/:list_id", func(c *gin.Context) {
		// Handler logic for listing bookmarks in a specific list
	})

	r.GET("/me/tags", func(c *gin.Context) {
		// Handler logic for listing user tags
	})

	r.GET("/lists", func(c *gin.Context) {
		// Handler logic for listing all lists
	})

	r.GET("/bookmarks", func(c *gin.Context) {
		// Handler logic for listing all bookmarks
	})

	r.GET("/tags", func(c *gin.Context) {
		// Handler logic for listing all tags
	})

	r.Run(":8080")
}
