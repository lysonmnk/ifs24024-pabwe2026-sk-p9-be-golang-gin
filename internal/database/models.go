package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Username  string         `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Role      string         `gorm:"size:50;not null;default:'author'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Posts     []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

type Post struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Slug        string         `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Category    string         `gorm:"size:100;not null;index" json:"category"`
	Excerpt     string         `gorm:"type:text" json:"excerpt"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	ReadingTime int            `gorm:"not null;default:1" json:"reading_time"`
	IsPublished bool           `gorm:"not null;default:false;index" json:"is_published"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        *User          `gorm:"foreignKey:UserID" json:"author,omitempty"`
}
