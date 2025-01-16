// post.go — post-passport.
// Author: d28035203

package models

import "time"

// Post represents a blog post or article
type Post struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string      `gorm:"size:200" json:"title"`
	Content     string      `gorm:"type:text" json:"content"`
	UserID      uint        `gorm:"not null" json:"user_id"`
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Views       int         `gorm:"default:0" json:"views"`
}

// TableName returns the table name for GORM
func (Post) TableName() string {
	return "posts"
}
