// user.go — fuzzy-adventure.
// Author: d28035203

package models

import "time"

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:100" json:"email"`
	Password  string         `gorm:"size:255" json:"password,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Posts     []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

// TableName returns the table name for GORM
func (User) TableName() string {
	return "users"
}
