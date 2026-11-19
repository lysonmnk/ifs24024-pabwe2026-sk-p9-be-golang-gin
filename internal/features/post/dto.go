package post

import "time"

type CreatePostRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Category    string `json:"category" binding:"max=100"`
	Excerpt     string `json:"excerpt"`
	Content     string `json:"content" binding:"required"`
	IsPublished bool   `json:"is_published"`
}

type UpdatePostRequest struct {
	Title       string `json:"title"`
	Category    string `json:"category"`
	Excerpt     string `json:"excerpt"`
	Content     string `json:"content"`
	IsPublished bool   `json:"is_published"`
}

type AuthorSummary struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type PostResponse struct {
	ID          uint           `json:"id"`
	UserID      uint           `json:"user_id,omitempty"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	Category    string         `json:"category"`
	Excerpt     string         `json:"excerpt"`
	ReadingTime int            `json:"reading_time"`
	AuthorName  string         `json:"author_name"`
	IsPublished bool           `json:"is_published"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Author      *AuthorSummary `json:"author,omitempty"`
}

type PostDetailResponse struct {
	ID          uint           `json:"id"`
	UserID      uint           `json:"user_id,omitempty"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	Category    string         `json:"category"`
	Excerpt     string         `json:"excerpt"`
	Content     string         `json:"content"`
	ReadingTime int            `json:"reading_time"`
	AuthorName  string         `json:"author_name"`
	IsPublished bool           `json:"is_published"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Author      *AuthorSummary `json:"author,omitempty"`
}

type PostStatsResponse struct {
	TotalPosts     int64 `json:"total_posts"`
	PublishedPosts int64 `json:"published_posts"`
	DraftPosts     int64 `json:"draft_posts"`
	TotalArticles  int64 `json:"total_articles"`
	TotalSelesai   int64 `json:"total_selesai"`
	BelumSelesai   int64 `json:"belum_selesai"`
}

type AdminUserResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	TotalPosts int64     `json:"total_posts"`
	CreatedAt  time.Time `json:"created_at"`
}
