package post

import (
	"errors"
	"fmt"
	"time"

	"app/internal/database"
)

type Service interface {
	GetCategories() ([]string, error)
	GetPublicPosts(keyword, category string) ([]PostResponse, error)
	GetPublicPostBySlug(slug string) (*PostDetailResponse, error)
	GetAuthorPosts(userID uint) ([]PostResponse, error)
	GetAuthorPostStats(userID uint) (*PostStatsResponse, error)
	CreatePost(userID uint, req *CreatePostRequest) (*PostDetailResponse, error)
	GetAuthorPostByID(id uint, userID uint) (*PostDetailResponse, error)
	UpdateAuthorPost(id uint, userID uint, req *UpdatePostRequest) (*PostDetailResponse, error)
	TogglePublish(id uint, userID uint) (*PostResponse, error)
	DeleteAuthorPost(id uint, userID uint) error
	GetAllUsers() ([]AdminUserResponse, error)
	GetAllPostsArchive() ([]PostResponse, error)
	GetUserPosts(userID uint) ([]PostResponse, error)
}

type postService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &postService{repo: repo}
}

func (s *postService) ensureUniqueSlug(title string, excludeID uint) (string, error) {
	baseSlug := GenerateSlug(title)
	slug := baseSlug
	counter := 1

	for {
		exists, err := s.repo.SlugExists(slug, excludeID)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return slug, nil
}

func (s *postService) GetCategories() ([]string, error) {
	cats, err := s.repo.GetCategories()
	if err != nil {
		cats = []string{}
	}
	defaultCats := []string{"Bisnis", "Gaya Hidup", "Opini", "Pendidikan", "Teknologi", "Umum"}
	seen := make(map[string]bool)
	var result []string
	for _, c := range cats {
		if c != "" && !seen[c] {
			seen[c] = true
			result = append(result, c)
		}
	}
	for _, c := range defaultCats {
		if !seen[c] {
			seen[c] = true
			result = append(result, c)
		}
	}
	return result, nil
}

func (s *postService) GetPublicPosts(keyword, category string) ([]PostResponse, error) {
	posts, err := s.repo.GetPublicPosts(keyword, category)
	if err != nil {
		return nil, err
	}

	responses := make([]PostResponse, 0)
	for _, p := range posts {
		var author *AuthorSummary
		if p.User != nil {
			author = &AuthorSummary{
				ID:       p.User.ID,
				Name:     p.User.Name,
				Username: p.User.Username,
			}
		}

		responses = append(responses, PostResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			Title:       p.Title,
			Slug:        p.Slug,
			Category:    p.Category,
			Excerpt:     p.Excerpt,
			ReadingTime: p.ReadingTime,
			IsPublished: p.IsPublished,
			PublishedAt: p.PublishedAt,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			Author:      author,
		})
	}

	return responses, nil
}

func (s *postService) GetPublicPostBySlug(slug string) (*PostDetailResponse, error) {
	p, err := s.repo.GetPublicPostBySlug(slug)
	if err != nil || p == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	var author *AuthorSummary
	if p.User != nil {
		author = &AuthorSummary{
			ID:       p.User.ID,
			Name:     p.User.Name,
			Username: p.User.Username,
		}
	}

	return &PostDetailResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		Title:       p.Title,
		Slug:        p.Slug,
		Category:    p.Category,
		Excerpt:     p.Excerpt,
		Content:     p.Content,
		ReadingTime: p.ReadingTime,
		IsPublished: p.IsPublished,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Author:      author,
	}, nil
}

func (s *postService) GetAuthorPosts(userID uint) ([]PostResponse, error) {
	posts, err := s.repo.GetAuthorPosts(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]PostResponse, 0)
	for _, p := range posts {
		responses = append(responses, PostResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			Title:       p.Title,
			Slug:        p.Slug,
			Category:    p.Category,
			Excerpt:     p.Excerpt,
			ReadingTime: p.ReadingTime,
			IsPublished: p.IsPublished,
			PublishedAt: p.PublishedAt,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *postService) GetAuthorPostStats(userID uint) (*PostStatsResponse, error) {
	return s.repo.GetAuthorPostStats(userID)
}

func (s *postService) CreatePost(userID uint, req *CreatePostRequest) (*PostDetailResponse, error) {
	if req.Category == "" {
		req.Category = "Umum"
	}
	if req.Excerpt == "" {
		if len(req.Content) > 150 {
			req.Excerpt = req.Content[:150] + "..."
		} else {
			req.Excerpt = req.Content
		}
	}

	slug, err := s.ensureUniqueSlug(req.Title, 0)
	if err != nil {
		return nil, err
	}

	readingTime := CalculateReadingTime(req.Content)

	var publishedAt *time.Time
	if req.IsPublished {
		now := time.Now()
		publishedAt = &now
	}

	post := &database.Post{
		UserID:      userID,
		Title:       req.Title,
		Slug:        slug,
		Category:    req.Category,
		Excerpt:     req.Excerpt,
		Content:     req.Content,
		ReadingTime: readingTime,
		IsPublished: req.IsPublished,
		PublishedAt: publishedAt,
	}

	if err := s.repo.CreatePost(post); err != nil {
		return nil, err
	}

	return &PostDetailResponse{
		ID:          post.ID,
		UserID:      post.UserID,
		Title:       post.Title,
		Slug:        post.Slug,
		Category:    post.Category,
		Excerpt:     post.Excerpt,
		Content:     post.Content,
		ReadingTime: post.ReadingTime,
		IsPublished: post.IsPublished,
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}, nil
}

func (s *postService) GetAuthorPostByID(id uint, userID uint) (*PostDetailResponse, error) {
	p, err := s.repo.GetPostByID(id)
	if err != nil || p == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	if p.UserID != userID {
		return nil, errors.New("anda tidak memiliki akses ke artikel ini")
	}

	return &PostDetailResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		Title:       p.Title,
		Slug:        p.Slug,
		Category:    p.Category,
		Excerpt:     p.Excerpt,
		Content:     p.Content,
		ReadingTime: p.ReadingTime,
		IsPublished: p.IsPublished,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func (s *postService) UpdateAuthorPost(id uint, userID uint, req *UpdatePostRequest) (*PostDetailResponse, error) {
	p, err := s.repo.GetPostByID(id)
	if err != nil || p == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	if p.UserID != userID {
		return nil, errors.New("anda tidak memiliki akses ke artikel ini")
	}

	if p.Title != req.Title {
		slug, err := s.ensureUniqueSlug(req.Title, p.ID)
		if err != nil {
			return nil, err
		}
		p.Slug = slug
	}

	category := req.Category
	if category == "" {
		category = "Umum"
	}
	excerpt := req.Excerpt
	if excerpt == "" {
		if len(req.Content) > 150 {
			excerpt = req.Content[:150] + "..."
		} else {
			excerpt = req.Content
		}
	}

	p.Title = req.Title
	p.Category = category
	p.Excerpt = excerpt
	p.Content = req.Content
	p.ReadingTime = CalculateReadingTime(req.Content)

	if req.IsPublished && !p.IsPublished {
		now := time.Now()
		p.PublishedAt = &now
	} else if !req.IsPublished && p.IsPublished {
		p.PublishedAt = nil
	}
	p.IsPublished = req.IsPublished

	if err := s.repo.UpdatePost(p); err != nil {
		return nil, err
	}

	return &PostDetailResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		Title:       p.Title,
		Slug:        p.Slug,
		Category:    p.Category,
		Excerpt:     p.Excerpt,
		Content:     p.Content,
		ReadingTime: p.ReadingTime,
		IsPublished: p.IsPublished,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func (s *postService) TogglePublish(id uint, userID uint) (*PostResponse, error) {
	p, err := s.repo.GetPostByID(id)
	if err != nil || p == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	if p.UserID != userID {
		return nil, errors.New("anda tidak memiliki akses ke artikel ini")
	}

	p.IsPublished = !p.IsPublished
	if p.IsPublished {
		now := time.Now()
		p.PublishedAt = &now
	} else {
		p.PublishedAt = nil
	}

	if err := s.repo.UpdatePost(p); err != nil {
		return nil, err
	}

	return &PostResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		Title:       p.Title,
		Slug:        p.Slug,
		Category:    p.Category,
		Excerpt:     p.Excerpt,
		ReadingTime: p.ReadingTime,
		IsPublished: p.IsPublished,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func (s *postService) DeleteAuthorPost(id uint, userID uint) error {
	p, err := s.repo.GetPostByID(id)
	if err != nil || p == nil {
		return errors.New("artikel tidak ditemukan")
	}

	if p.UserID != userID {
		return errors.New("anda tidak memiliki akses ke artikel ini")
	}

	return s.repo.DeletePost(id)
}

func (s *postService) GetAllUsers() ([]AdminUserResponse, error) {
	return s.repo.GetAllUsersWithPostCount()
}

func (s *postService) GetAllPostsArchive() ([]PostResponse, error) {
	posts, err := s.repo.GetAllPostsArchive()
	if err != nil {
		return nil, err
	}

	responses := make([]PostResponse, 0)
	for _, p := range posts {
		var author *AuthorSummary
		if p.User != nil {
			author = &AuthorSummary{
				ID:       p.User.ID,
				Name:     p.User.Name,
				Username: p.User.Username,
			}
		}

		responses = append(responses, PostResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			Title:       p.Title,
			Slug:        p.Slug,
			Category:    p.Category,
			Excerpt:     p.Excerpt,
			ReadingTime: p.ReadingTime,
			IsPublished: p.IsPublished,
			PublishedAt: p.PublishedAt,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			Author:      author,
		})
	}
	return responses, nil
}

func (s *postService) GetUserPosts(userID uint) ([]PostResponse, error) {
	posts, err := s.repo.GetPostsByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]PostResponse, 0)
	for _, p := range posts {
		responses = append(responses, PostResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			Title:       p.Title,
			Slug:        p.Slug,
			Category:    p.Category,
			Excerpt:     p.Excerpt,
			ReadingTime: p.ReadingTime,
			IsPublished: p.IsPublished,
			PublishedAt: p.PublishedAt,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}
	return responses, nil
}
