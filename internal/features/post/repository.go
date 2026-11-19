package post

import (
	"app/internal/database"

	"gorm.io/gorm"
)

type Repository interface {
	GetCategories() ([]string, error)
	GetPublicPosts(keyword, category string) ([]database.Post, error)
	GetPublicPostBySlug(slug string) (*database.Post, error)
	GetAuthorPosts(userID uint) ([]database.Post, error)
	GetAuthorPostStats(userID uint) (*PostStatsResponse, error)
	CreatePost(post *database.Post) error
	GetPostByID(id uint) (*database.Post, error)
	UpdatePost(post *database.Post) error
	DeletePost(id uint) error
	SlugExists(slug string, excludeID uint) (bool, error)
	GetAllUsersWithPostCount() ([]AdminUserResponse, error)
	GetAllPostsArchive() ([]database.Post, error)
	GetPostsByUserID(userID uint) ([]database.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &postRepository{db: db}
}

func (r *postRepository) GetCategories() ([]string, error) {
	var categories []string
	err := r.db.Model(&database.Post{}).
		Where("is_published = ?", true).
		Distinct("category").
		Order("category asc").
		Pluck("category", &categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *postRepository) GetPublicPosts(keyword, category string) ([]database.Post, error) {
	var posts []database.Post
	query := r.db.Model(&database.Post{}).
		Preload("User").
		Where("is_published = ?", true)

	if keyword != "" {
		likeTerm := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ? OR excerpt LIKE ?", likeTerm, likeTerm, likeTerm)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Order("published_at desc, created_at desc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetPublicPostBySlug(slug string) (*database.Post, error) {
	var post database.Post
	err := r.db.Preload("User").
		Where("slug = ? AND is_published = ?", slug, true).
		First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetAuthorPosts(userID uint) ([]database.Post, error) {
	var posts []database.Post
	err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetAuthorPostStats(userID uint) (*PostStatsResponse, error) {
	var total, published int64

	if err := r.db.Model(&database.Post{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, err
	}
	_ = r.db.Model(&database.Post{}).Where("user_id = ? AND is_published = ?", userID, true).Count(&published)

	drafts := total - published
	return &PostStatsResponse{
		TotalPosts:     total,
		PublishedPosts: published,
		DraftPosts:     drafts,
		TotalArticles:  total,
		TotalSelesai:   published,
		BelumSelesai:   drafts,
	}, nil
}

func (r *postRepository) CreatePost(post *database.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetPostByID(id uint) (*database.Post, error) {
	var post database.Post
	if err := r.db.Preload("User").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) UpdatePost(post *database.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) DeletePost(id uint) error {
	return r.db.Delete(&database.Post{}, id).Error
}

func (r *postRepository) SlugExists(slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&database.Post{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *postRepository) GetAllUsersWithPostCount() ([]AdminUserResponse, error) {
	type userWithCount struct {
		database.User
		PostCount int64 `gorm:"column:post_count"`
	}

	var results []userWithCount
	err := r.db.Table("users").
		Select("users.*, COUNT(posts.id) as post_count").
		Joins("LEFT JOIN posts ON posts.user_id = users.id AND posts.deleted_at IS NULL").
		Where("users.deleted_at IS NULL").
		Group("users.id").
		Order("users.created_at desc").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	users := make([]AdminUserResponse, 0)
	for _, res := range results {
		users = append(users, AdminUserResponse{
			ID:         res.ID,
			Name:       res.Name,
			Username:   res.Username,
			Email:      res.Email,
			Role:       res.Role,
			TotalPosts: res.PostCount,
			CreatedAt:  res.CreatedAt,
		})
	}
	return users, nil
}

func (r *postRepository) GetAllPostsArchive() ([]database.Post, error) {
	var posts []database.Post
	err := r.db.Preload("User").Order("created_at desc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetPostsByUserID(userID uint) ([]database.Post, error) {
	var posts []database.Post
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}
