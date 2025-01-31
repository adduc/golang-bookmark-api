package internal

import (
	"time"

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
