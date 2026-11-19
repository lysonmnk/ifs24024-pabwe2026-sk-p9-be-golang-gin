package post

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupPostTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)
	err = db.AutoMigrate(&database.User{}, &database.Post{})
	assert.NoError(t, err)
	return db
}

func TestUtils(t *testing.T) {
	t.Run("GenerateSlug", func(t *testing.T) {
		assert.Equal(t, "halo-dunia", GenerateSlug("Halo Dunia!"))
		assert.Equal(t, "artikel", GenerateSlug("!!!@@@###"))
	})

	t.Run("CalculateReadingTime", func(t *testing.T) {
		assert.Equal(t, 1, CalculateReadingTime(""))
		assert.Equal(t, 1, CalculateReadingTime("Satu dua tiga empat"))
		longContent := ""
		for i := 0; i < 450; i++ {
			longContent += "kata "
		}
		assert.Equal(t, 3, CalculateReadingTime(longContent))
	})
}

func TestPostRepository(t *testing.T) {
	db := setupPostTestDB(t)
	repo := NewRepository(db)

	user := &database.User{
		Name:     "Author Test",
		Username: "author_post",
		Email:    "author_post@contoh.com",
		Password: "password123",
		Role:     "author",
	}
	assert.NoError(t, db.Create(user).Error)

	now := time.Now()
	post1 := &database.Post{
		UserID:      user.ID,
		Title:       "Belajar Go Gin",
		Slug:        "belajar-go-gin",
		Category:    "Golang",
		Excerpt:     "Ringkasan Go Gin",
		Content:     "Konten lengkap tutorial Go Gin",
		ReadingTime: 2,
		IsPublished: true,
		PublishedAt: &now,
	}
	assert.NoError(t, repo.CreatePost(post1))

	post2 := &database.Post{
		UserID:      user.ID,
		Title:       "Draf Artikel Baru",
		Slug:        "draf-artikel-baru",
		Category:    "Tutorial",
		Excerpt:     "Draf excerpt",
		Content:     "Draf content",
		ReadingTime: 1,
		IsPublished: false,
	}
	assert.NoError(t, repo.CreatePost(post2))

	t.Run("GetCategories", func(t *testing.T) {
		cats, err := repo.GetCategories()
		assert.NoError(t, err)
		assert.Contains(t, cats, "Golang")
	})

	t.Run("GetPublicPosts with and without query", func(t *testing.T) {
		posts, err := repo.GetPublicPosts("", "")
		assert.NoError(t, err)
		assert.Len(t, posts, 1)

		postsQ, err := repo.GetPublicPosts("Gin", "Golang")
		assert.NoError(t, err)
		assert.Len(t, postsQ, 1)

		postsEmpty, err := repo.GetPublicPosts("TidakAda", "")
		assert.NoError(t, err)
		assert.Len(t, postsEmpty, 0)
	})

	t.Run("GetPublicPostBySlug", func(t *testing.T) {
		p, err := repo.GetPublicPostBySlug("belajar-go-gin")
		assert.NoError(t, err)
		assert.Equal(t, post1.ID, p.ID)

		_, err = repo.GetPublicPostBySlug("not-exist")
		assert.Error(t, err)
	})

	t.Run("GetAuthorPosts and Stats", func(t *testing.T) {
		posts, err := repo.GetAuthorPosts(user.ID)
		assert.NoError(t, err)
		assert.Len(t, posts, 2)

		stats, err := repo.GetAuthorPostStats(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), stats.TotalPosts)
		assert.Equal(t, int64(1), stats.PublishedPosts)
		assert.Equal(t, int64(1), stats.DraftPosts)
	})

	t.Run("GetPostByID, UpdatePost, DeletePost", func(t *testing.T) {
		p, err := repo.GetPostByID(post1.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Belajar Go Gin", p.Title)

		p.Title = "Belajar Go Gin Updated"
		assert.NoError(t, repo.UpdatePost(p))

		exists, err := repo.SlugExists("belajar-go-gin", 0)
		assert.NoError(t, err)
		assert.True(t, exists)

		existsSelf, err := repo.SlugExists("belajar-go-gin", post1.ID)
		assert.NoError(t, err)
		assert.False(t, existsSelf)

		assert.NoError(t, repo.DeletePost(post2.ID))
	})

	t.Run("Admin queries", func(t *testing.T) {
		users, err := repo.GetAllUsersWithPostCount()
		assert.NoError(t, err)
		assert.NotEmpty(t, users)

		archive, err := repo.GetAllPostsArchive()
		assert.NoError(t, err)
		assert.NotEmpty(t, archive)

		userPosts, err := repo.GetPostsByUserID(user.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, userPosts)
	})

	t.Run("Closed DB repository errors", func(t *testing.T) {
		closedDB := setupPostTestDB(t)
		sqlDB, _ := closedDB.DB()
		_ = sqlDB.Close()
		cRepo := NewRepository(closedDB)

		_, err := cRepo.GetCategories()
		assert.Error(t, err)

		_, err = cRepo.GetPublicPosts("", "")
		assert.Error(t, err)

		_, err = cRepo.GetAuthorPosts(1)
		assert.Error(t, err)

		_, err = cRepo.GetAuthorPostStats(1)
		assert.Error(t, err)

		_, err = cRepo.SlugExists("test", 0)
		assert.Error(t, err)

		_, err = cRepo.GetAllUsersWithPostCount()
		assert.Error(t, err)

		_, err = cRepo.GetAllPostsArchive()
		assert.Error(t, err)

		_, err = cRepo.GetPostsByUserID(1)
		assert.Error(t, err)
	})
}

func TestPostService(t *testing.T) {
	db := setupPostTestDB(t)
	repo := NewRepository(db)
	svc := NewService(repo)

	user := &database.User{
		Name:     "Author Service",
		Username: "author_svc",
		Email:    "author_svc@contoh.com",
		Password: "password123",
		Role:     "author",
	}
	assert.NoError(t, db.Create(user).Error)

	otherUser := &database.User{
		Name:     "Other User",
		Username: "other_user",
		Email:    "other@contoh.com",
		Password: "password123",
		Role:     "author",
	}
	assert.NoError(t, db.Create(otherUser).Error)

	t.Run("CreatePost draft and published", func(t *testing.T) {
		reqDraft := &CreatePostRequest{
			Title:       "Post Pertama",
			Category:    "Teknologi",
			Excerpt:     "Excerpt",
			Content:     "Konten pertama",
			IsPublished: false,
		}
		p1, err := svc.CreatePost(user.ID, reqDraft)
		assert.NoError(t, err)
		assert.Equal(t, "post-pertama", p1.Slug)
		assert.False(t, p1.IsPublished)
		assert.Nil(t, p1.PublishedAt)

		// Duplicate title -> slug receives suffix
		p2, err := svc.CreatePost(user.ID, reqDraft)
		assert.NoError(t, err)
		assert.Equal(t, "post-pertama-1", p2.Slug)

		// Published post
		reqPub := &CreatePostRequest{
			Title:       "Post Terbit",
			Category:    "Teknologi",
			Excerpt:     "Excerpt terbit",
			Content:     "Konten terbit",
			IsPublished: true,
		}
		pPub, err := svc.CreatePost(user.ID, reqPub)
		assert.NoError(t, err)
		assert.True(t, pPub.IsPublished)
		assert.NotNil(t, pPub.PublishedAt)
	})

	t.Run("GetCategories and GetPublicPosts and Detail", func(t *testing.T) {
		cats, err := svc.GetCategories()
		assert.NoError(t, err)
		assert.Contains(t, cats, "Teknologi")

		posts, err := svc.GetPublicPosts("", "")
		assert.NoError(t, err)
		assert.NotEmpty(t, posts)

		detail, err := svc.GetPublicPostBySlug("post-terbit")
		assert.NoError(t, err)
		assert.Equal(t, "Post Terbit", detail.Title)

		_, err = svc.GetPublicPostBySlug("slug-palsu")
		assert.Error(t, err)
	})

	t.Run("GetAuthorPosts and Stats", func(t *testing.T) {
		posts, err := svc.GetAuthorPosts(user.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, posts)

		stats, err := svc.GetAuthorPostStats(user.ID)
		assert.NoError(t, err)
		assert.True(t, stats.TotalPosts >= 3)
	})

	t.Run("GetAuthorPostByID and UpdateAuthorPost", func(t *testing.T) {
		authorPosts, _ := svc.GetAuthorPosts(user.ID)
		firstID := authorPosts[0].ID

		// Success get
		p, err := svc.GetAuthorPostByID(firstID, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, firstID, p.ID)

		// Unauthorized get
		_, err = svc.GetAuthorPostByID(firstID, otherUser.ID)
		assert.Error(t, err)

		// Not found get
		_, err = svc.GetAuthorPostByID(99999, user.ID)
		assert.Error(t, err)

		// Update post title and publish status
		upReq := &UpdatePostRequest{
			Title:       "Judul Diubah",
			Category:    "Tutorial",
			Excerpt:     "Excerpt baru",
			Content:     "Konten baru",
			IsPublished: true,
		}
		upRes, err := svc.UpdateAuthorPost(firstID, user.ID, upReq)
		assert.NoError(t, err)
		assert.Equal(t, "Judul Diubah", upRes.Title)
		assert.Equal(t, "judul-diubah", upRes.Slug)

		// Unpublish post update
		upReq.IsPublished = false
		upRes2, err := svc.UpdateAuthorPost(firstID, user.ID, upReq)
		assert.NoError(t, err)
		assert.False(t, upRes2.IsPublished)

		// Update unauthorized
		_, err = svc.UpdateAuthorPost(firstID, otherUser.ID, upReq)
		assert.Error(t, err)

		// Update not found
		_, err = svc.UpdateAuthorPost(99999, user.ID, upReq)
		assert.Error(t, err)
	})

	t.Run("TogglePublish and DeleteAuthorPost", func(t *testing.T) {
		authorPosts, _ := svc.GetAuthorPosts(user.ID)
		firstID := authorPosts[0].ID

		// Toggle to publish (true)
		res, err := svc.TogglePublish(firstID, user.ID)
		assert.NoError(t, err)
		assert.True(t, res.IsPublished)

		// Toggle to draft (false)
		res2, err := svc.TogglePublish(firstID, user.ID)
		assert.NoError(t, err)
		assert.False(t, res2.IsPublished)

		// Toggle unauthorized
		_, err = svc.TogglePublish(firstID, otherUser.ID)
		assert.Error(t, err)

		// Toggle not found
		_, err = svc.TogglePublish(99999, user.ID)
		assert.Error(t, err)

		// Delete unauthorized
		err = svc.DeleteAuthorPost(firstID, otherUser.ID)
		assert.Error(t, err)

		// Delete not found
		err = svc.DeleteAuthorPost(99999, user.ID)
		assert.Error(t, err)

		// Delete success
		err = svc.DeleteAuthorPost(firstID, user.ID)
		assert.NoError(t, err)
	})

	t.Run("Admin service queries", func(t *testing.T) {
		users, err := svc.GetAllUsers()
		assert.NoError(t, err)
		assert.NotEmpty(t, users)

		archive, err := svc.GetAllPostsArchive()
		assert.NoError(t, err)
		assert.NotEmpty(t, archive)

		userPosts, err := svc.GetUserPosts(user.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, userPosts)
	})

	t.Run("Service errors with mock repo", func(t *testing.T) {
		mockR := &mockPostRepoErr{repo}
		mockS := NewService(mockR)

		_, err := mockS.CreatePost(1, &CreatePostRequest{Title: "A", Category: "B", Content: "C"})
		assert.Error(t, err)

		_, err = mockS.GetPublicPosts("", "")
		assert.Error(t, err)

		_, err = mockS.GetAuthorPosts(1)
		assert.Error(t, err)

		_, err = mockS.GetAllPostsArchive()
		assert.Error(t, err)

		_, err = mockS.GetUserPosts(1)
		assert.Error(t, err)

		// CreatePost repo error
		_, err = mockS.CreatePost(1, &CreatePostRequest{Title: "Judul Post", Category: "B", Content: "C"})
		assert.Error(t, err)

		// CreatePost slug error
		_, err = mockS.CreatePost(1, &CreatePostRequest{Title: "force-slug-err", Category: "B", Content: "C"})
		assert.Error(t, err)

		// UpdateAuthorPost update repo error
		_, err = mockS.UpdateAuthorPost(1, 1, &UpdatePostRequest{Title: "Judul Post", Category: "B", Content: "C"})
		assert.Error(t, err)

		// UpdateAuthorPost slug error
		_, err = mockS.UpdateAuthorPost(1, 1, &UpdatePostRequest{Title: "force-slug-err", Category: "B", Content: "C"})
		assert.Error(t, err)

		// TogglePublish update repo error
		_, err = mockS.TogglePublish(1, 1)
		assert.Error(t, err)
	})
}

type mockPostRepoErr struct {
	Repository
}

func (m *mockPostRepoErr) SlugExists(slug string, excludeID uint) (bool, error) {
	if slug == "force-slug-err" {
		return false, errors.New("slug error")
	}
	return false, nil
}
func (m *mockPostRepoErr) CreatePost(p *database.Post) error {
	return errors.New("create error")
}
func (m *mockPostRepoErr) UpdatePost(p *database.Post) error {
	return errors.New("update error")
}
func (m *mockPostRepoErr) GetPostByID(id uint) (*database.Post, error) {
	return &database.Post{ID: id, UserID: 1, Title: "Existing", Category: "C", Content: "K"}, nil
}
func (m *mockPostRepoErr) GetPublicPosts(q, cat string) ([]database.Post, error) {
	return nil, errors.New("db error")
}
func (m *mockPostRepoErr) GetAuthorPosts(userID uint) ([]database.Post, error) {
	return nil, errors.New("db error")
}
func (m *mockPostRepoErr) GetAllPostsArchive() ([]database.Post, error) {
	return nil, errors.New("db error")
}
func (m *mockPostRepoErr) GetPostsByUserID(userID uint) ([]database.Post, error) {
	return nil, errors.New("db error")
}

func TestPostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupPostTestDB(t)
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	user := &database.User{Name: "Post Author", Username: "post_author", Email: "p@c.com", Password: "p", Role: "author"}
	db.Create(user)
	now := time.Now()
	pubPost := &database.Post{
		UserID:      user.ID,
		Title:       "Post Publik Sukses",
		Slug:        "post-publik-sukses",
		Category:    "Tech",
		Content:     "Konten",
		IsPublished: true,
		PublishedAt: &now,
	}
	db.Create(pubPost)

	t.Run("GetCategories and GetPublicPosts", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetCategories(c)
		assert.Equal(t, http.StatusOK, w.Code)

		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		c2.Request, _ = http.NewRequest(http.MethodGet, "/api/public/posts?q=test&category=tech", nil)
		handler.GetPublicPosts(c2)
		assert.Equal(t, http.StatusOK, w2.Code)
	})

	t.Run("GetPublicPostBySlug", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "slug", Value: "not-found"}}
		handler.GetPublicPostBySlug(c)
		assert.Equal(t, http.StatusNotFound, w.Code)

		wOK := httptest.NewRecorder()
		cOK, _ := gin.CreateTestContext(wOK)
		cOK.Params = gin.Params{{Key: "slug", Value: "post-publik-sukses"}}
		handler.GetPublicPostBySlug(cOK)
		assert.Equal(t, http.StatusOK, wOK.Code)
	})

	t.Run("Author posts endpoints", func(t *testing.T) {
		// Unauthorized GetAuthorPosts
		wNoAuth := httptest.NewRecorder()
		cNoAuth, _ := gin.CreateTestContext(wNoAuth)
		handler.GetAuthorPosts(cNoAuth)
		assert.Equal(t, http.StatusUnauthorized, wNoAuth.Code)

		// Success GetAuthorPosts
		wAuth := httptest.NewRecorder()
		cAuth, _ := gin.CreateTestContext(wAuth)
		cAuth.Set("userID", user.ID)
		handler.GetAuthorPosts(cAuth)
		assert.Equal(t, http.StatusOK, wAuth.Code)

		// Unauthorized Stats
		wStatNoAuth := httptest.NewRecorder()
		cStatNoAuth, _ := gin.CreateTestContext(wStatNoAuth)
		handler.GetAuthorPostStats(cStatNoAuth)
		assert.Equal(t, http.StatusUnauthorized, wStatNoAuth.Code)

		// Success Stats
		wStat := httptest.NewRecorder()
		cStat, _ := gin.CreateTestContext(wStat)
		cStat.Set("userID", user.ID)
		handler.GetAuthorPostStats(cStat)
		assert.Equal(t, http.StatusOK, wStat.Code)

		// CreatePost Unauthorized
		wCreateNoAuth := httptest.NewRecorder()
		cCreateNoAuth, _ := gin.CreateTestContext(wCreateNoAuth)
		handler.CreatePost(cCreateNoAuth)
		assert.Equal(t, http.StatusUnauthorized, wCreateNoAuth.Code)

		// CreatePost Bad JSON
		wCreateBad := httptest.NewRecorder()
		cCreateBad, _ := gin.CreateTestContext(wCreateBad)
		cCreateBad.Set("userID", user.ID)
		cCreateBad.Request, _ = http.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer([]byte("bad")))
		handler.CreatePost(cCreateBad)
		assert.Equal(t, http.StatusBadRequest, wCreateBad.Code)

		// CreatePost Success
		createBody, _ := json.Marshal(CreatePostRequest{
			Title:       "Post Handler",
			Category:    "Golang",
			Excerpt:     "Excerpt",
			Content:     "Content",
			IsPublished: true,
		})
		wCreate := httptest.NewRecorder()
		cCreate, _ := gin.CreateTestContext(wCreate)
		cCreate.Set("userID", user.ID)
		cCreate.Request, _ = http.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(createBody))
		cCreate.Request.Header.Set("Content-Type", "application/json")
		handler.CreatePost(cCreate)
		assert.Equal(t, http.StatusCreated, wCreate.Code)

		var postResp struct {
			Data PostDetailResponse `json:"data"`
		}
		_ = json.Unmarshal(wCreate.Body.Bytes(), &postResp)
		postID := postResp.Data.ID

		// GetAuthorPostByID Unauthorized
		wGetNoAuth := httptest.NewRecorder()
		cGetNoAuth, _ := gin.CreateTestContext(wGetNoAuth)
		handler.GetAuthorPostByID(cGetNoAuth)
		assert.Equal(t, http.StatusUnauthorized, wGetNoAuth.Code)

		// GetAuthorPostByID Invalid ID
		wGetBadID := httptest.NewRecorder()
		cGetBadID, _ := gin.CreateTestContext(wGetBadID)
		cGetBadID.Set("userID", user.ID)
		cGetBadID.Params = gin.Params{{Key: "id", Value: "abc"}}
		handler.GetAuthorPostByID(cGetBadID)
		assert.Equal(t, http.StatusBadRequest, wGetBadID.Code)

		// GetAuthorPostByID Not Found
		wGetNF := httptest.NewRecorder()
		cGetNF, _ := gin.CreateTestContext(wGetNF)
		cGetNF.Set("userID", user.ID)
		cGetNF.Params = gin.Params{{Key: "id", Value: "99999"}}
		handler.GetAuthorPostByID(cGetNF)
		assert.Equal(t, http.StatusNotFound, wGetNF.Code)

		// GetAuthorPostByID Success
		wGetOK := httptest.NewRecorder()
		cGetOK, _ := gin.CreateTestContext(wGetOK)
		cGetOK.Set("userID", user.ID)
		cGetOK.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.GetAuthorPostByID(cGetOK)
		assert.Equal(t, http.StatusOK, wGetOK.Code)

		// UpdateAuthorPost Unauthorized
		wUpNoAuth := httptest.NewRecorder()
		cUpNoAuth, _ := gin.CreateTestContext(wUpNoAuth)
		handler.UpdateAuthorPost(cUpNoAuth)
		assert.Equal(t, http.StatusUnauthorized, wUpNoAuth.Code)

		// UpdateAuthorPost Invalid ID
		wUpBadID := httptest.NewRecorder()
		cUpBadID, _ := gin.CreateTestContext(wUpBadID)
		cUpBadID.Set("userID", user.ID)
		cUpBadID.Params = gin.Params{{Key: "id", Value: "abc"}}
		handler.UpdateAuthorPost(cUpBadID)
		assert.Equal(t, http.StatusBadRequest, wUpBadID.Code)

		// UpdateAuthorPost Bad JSON
		wUpBad := httptest.NewRecorder()
		cUpBad, _ := gin.CreateTestContext(wUpBad)
		cUpBad.Set("userID", user.ID)
		cUpBad.Params = gin.Params{{Key: "id", Value: "1"}}
		cUpBad.Request, _ = http.NewRequest(http.MethodPut, "/api/posts/1", bytes.NewBuffer([]byte("bad")))
		handler.UpdateAuthorPost(cUpBad)
		assert.Equal(t, http.StatusBadRequest, wUpBad.Code)

		// UpdateAuthorPost Success
		upBody, _ := json.Marshal(UpdatePostRequest{
			Title:       "Post Handler Updated",
			Category:    "Golang",
			Excerpt:     "Ex",
			Content:     "Con",
			IsPublished: true,
		})
		wUpOK := httptest.NewRecorder()
		cUpOK, _ := gin.CreateTestContext(wUpOK)
		cUpOK.Set("userID", user.ID)
		cUpOK.Params = gin.Params{{Key: "id", Value: "1"}}
		cUpOK.Request, _ = http.NewRequest(http.MethodPut, "/api/posts/1", bytes.NewBuffer(upBody))
		cUpOK.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateAuthorPost(cUpOK)
		assert.Equal(t, http.StatusOK, wUpOK.Code)

		// TogglePublish Unauthorized
		wTogNoAuth := httptest.NewRecorder()
		cTogNoAuth, _ := gin.CreateTestContext(wTogNoAuth)
		handler.TogglePublish(cTogNoAuth)
		assert.Equal(t, http.StatusUnauthorized, wTogNoAuth.Code)

		// TogglePublish Invalid ID
		wTogBadID := httptest.NewRecorder()
		cTogBadID, _ := gin.CreateTestContext(wTogBadID)
		cTogBadID.Set("userID", user.ID)
		cTogBadID.Params = gin.Params{{Key: "id", Value: "abc"}}
		handler.TogglePublish(cTogBadID)
		assert.Equal(t, http.StatusBadRequest, wTogBadID.Code)

		// TogglePublish Success
		wTogOK := httptest.NewRecorder()
		cTogOK, _ := gin.CreateTestContext(wTogOK)
		cTogOK.Set("userID", user.ID)
		cTogOK.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.TogglePublish(cTogOK)
		assert.Equal(t, http.StatusOK, wTogOK.Code)

		// DeleteAuthorPost Unauthorized
		wDelNoAuth := httptest.NewRecorder()
		cDelNoAuth, _ := gin.CreateTestContext(wDelNoAuth)
		handler.DeleteAuthorPost(cDelNoAuth)
		assert.Equal(t, http.StatusUnauthorized, wDelNoAuth.Code)

		// DeleteAuthorPost Invalid ID
		wDelBadID := httptest.NewRecorder()
		cDelBadID, _ := gin.CreateTestContext(wDelBadID)
		cDelBadID.Set("userID", user.ID)
		cDelBadID.Params = gin.Params{{Key: "id", Value: "abc"}}
		handler.DeleteAuthorPost(cDelBadID)
		assert.Equal(t, http.StatusBadRequest, wDelBadID.Code)

		// Author post 403 Forbidden checks (post 1 still exists here)
		wGetOther := httptest.NewRecorder()
		cGetOther, _ := gin.CreateTestContext(wGetOther)
		cGetOther.Set("userID", uint(99999))
		cGetOther.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.GetAuthorPostByID(cGetOther)
		assert.Equal(t, http.StatusForbidden, wGetOther.Code)

		wGetNF2 := httptest.NewRecorder()
		cGetNF2, _ := gin.CreateTestContext(wGetNF2)
		cGetNF2.Set("userID", user.ID)
		cGetNF2.Params = gin.Params{{Key: "id", Value: "999999"}}
		handler.GetAuthorPostByID(cGetNF2)
		assert.Equal(t, http.StatusNotFound, wGetNF2.Code)

		wUpOther := httptest.NewRecorder()
		cUpOther, _ := gin.CreateTestContext(wUpOther)
		cUpOther.Set("userID", uint(99999))
		cUpOther.Params = gin.Params{{Key: "id", Value: "1"}}
		cUpOther.Request, _ = http.NewRequest(http.MethodPut, "/api/posts/1", bytes.NewBuffer([]byte(`{"title":"Edit"}`)))
		cUpOther.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateAuthorPost(cUpOther)
		assert.Equal(t, http.StatusForbidden, wUpOther.Code)

		wUpNF := httptest.NewRecorder()
		cUpNF, _ := gin.CreateTestContext(wUpNF)
		cUpNF.Set("userID", user.ID)
		cUpNF.Params = gin.Params{{Key: "id", Value: "999999"}}
		cUpNF.Request, _ = http.NewRequest(http.MethodPut, "/api/posts/999999", bytes.NewBuffer([]byte(`{"title":"Edit"}`)))
		cUpNF.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateAuthorPost(cUpNF)
		assert.Equal(t, http.StatusNotFound, wUpNF.Code)

		wTogOther := httptest.NewRecorder()
		cTogOther, _ := gin.CreateTestContext(wTogOther)
		cTogOther.Set("userID", uint(99999))
		cTogOther.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.TogglePublish(cTogOther)
		assert.Equal(t, http.StatusForbidden, wTogOther.Code)

		wDelOther := httptest.NewRecorder()
		cDelOther, _ := gin.CreateTestContext(wDelOther)
		cDelOther.Set("userID", uint(99999))
		cDelOther.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.DeleteAuthorPost(cDelOther)
		assert.Equal(t, http.StatusForbidden, wDelOther.Code)

		// DeleteAuthorPost Success
		wDelOK := httptest.NewRecorder()
		cDelOK, _ := gin.CreateTestContext(wDelOK)
		cDelOK.Set("userID", user.ID)
		cDelOK.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.DeleteAuthorPost(cDelOK)
		assert.Equal(t, http.StatusOK, wDelOK.Code)

		wDelNF := httptest.NewRecorder()
		cDelNF, _ := gin.CreateTestContext(wDelNF)
		cDelNF.Set("userID", user.ID)
		cDelNF.Params = gin.Params{{Key: "id", Value: "999999"}}
		handler.DeleteAuthorPost(cDelNF)
		assert.Equal(t, http.StatusNotFound, wDelNF.Code)

		_ = postID
	})

	t.Run("Admin handlers", func(t *testing.T) {
		wUsers := httptest.NewRecorder()
		cUsers, _ := gin.CreateTestContext(wUsers)
		handler.AdminGetAllUsers(cUsers)
		assert.Equal(t, http.StatusOK, wUsers.Code)

		wArchive := httptest.NewRecorder()
		cArchive, _ := gin.CreateTestContext(wArchive)
		handler.AdminGetAllPosts(cArchive)
		assert.Equal(t, http.StatusOK, wArchive.Code)

		// AdminGetUserPosts Invalid ID
		wUserPostsBad := httptest.NewRecorder()
		cUserPostsBad, _ := gin.CreateTestContext(wUserPostsBad)
		cUserPostsBad.Params = gin.Params{{Key: "id", Value: "abc"}}
		handler.AdminGetUserPosts(cUserPostsBad)
		assert.Equal(t, http.StatusBadRequest, wUserPostsBad.Code)

		// AdminGetUserPosts Success
		wUserPostsOK := httptest.NewRecorder()
		cUserPostsOK, _ := gin.CreateTestContext(wUserPostsOK)
		cUserPostsOK.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.AdminGetUserPosts(cUserPostsOK)
		assert.Equal(t, http.StatusOK, wUserPostsOK.Code)
	})

	t.Run("Handler error branches with error service", func(t *testing.T) {
		errSvc := &mockPostServiceErr{}
		errH := NewHandler(errSvc)

		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)
		errH.GetCategories(c1)
		assert.Equal(t, http.StatusInternalServerError, w1.Code)

		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		errH.GetPublicPosts(c2)
		assert.Equal(t, http.StatusInternalServerError, w2.Code)

		w3 := httptest.NewRecorder()
		c3, _ := gin.CreateTestContext(w3)
		c3.Set("userID", uint(1))
		errH.GetAuthorPosts(c3)
		assert.Equal(t, http.StatusInternalServerError, w3.Code)

		w4 := httptest.NewRecorder()
		c4, _ := gin.CreateTestContext(w4)
		c4.Set("userID", uint(1))
		errH.GetAuthorPostStats(c4)
		assert.Equal(t, http.StatusInternalServerError, w4.Code)

		w5 := httptest.NewRecorder()
		c5, _ := gin.CreateTestContext(w5)
		c5.Set("userID", uint(1))
		c5.Request, _ = http.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer([]byte(`{"title":"Judul","category":"Kategori","content":"Konten"}`)))
		c5.Request.Header.Set("Content-Type", "application/json")
		errH.CreatePost(c5)
		assert.Equal(t, http.StatusBadRequest, w5.Code)

		w6 := httptest.NewRecorder()
		c6, _ := gin.CreateTestContext(w6)
		c6.Set("userID", uint(1))
		c6.Params = gin.Params{{Key: "id", Value: "1"}}
		c6.Request, _ = http.NewRequest(http.MethodPut, "/api/posts/1", bytes.NewBuffer([]byte(`{"title":"Judul","category":"Kategori","content":"Konten"}`)))
		c6.Request.Header.Set("Content-Type", "application/json")
		errH.UpdateAuthorPost(c6)
		assert.Equal(t, http.StatusBadRequest, w6.Code)

		w7 := httptest.NewRecorder()
		c7, _ := gin.CreateTestContext(w7)
		c7.Set("userID", uint(1))
		c7.Params = gin.Params{{Key: "id", Value: "1"}}
		errH.TogglePublish(c7)
		assert.Equal(t, http.StatusBadRequest, w7.Code)

		w8 := httptest.NewRecorder()
		c8, _ := gin.CreateTestContext(w8)
		c8.Set("userID", uint(1))
		c8.Params = gin.Params{{Key: "id", Value: "1"}}
		errH.DeleteAuthorPost(c8)
		assert.Equal(t, http.StatusBadRequest, w8.Code)

		w9 := httptest.NewRecorder()
		c9, _ := gin.CreateTestContext(w9)
		errH.AdminGetAllUsers(c9)
		assert.Equal(t, http.StatusInternalServerError, w9.Code)

		w10 := httptest.NewRecorder()
		c10, _ := gin.CreateTestContext(w10)
		errH.AdminGetAllPosts(c10)
		assert.Equal(t, http.StatusInternalServerError, w10.Code)

		w11 := httptest.NewRecorder()
		c11, _ := gin.CreateTestContext(w11)
		c11.Params = gin.Params{{Key: "id", Value: "1"}}
		errH.AdminGetUserPosts(c11)
		assert.Equal(t, http.StatusInternalServerError, w11.Code)
	})
}

type mockPostServiceErr struct {
	Service
}

func (m *mockPostServiceErr) GetCategories() ([]string, error) {
	return nil, errors.New("err")
}
func (m *mockPostServiceErr) GetPublicPosts(q, c string) ([]PostResponse, error) {
	return nil, errors.New("err")
}
func (m *mockPostServiceErr) GetPublicPostBySlug(s string) (*PostDetailResponse, error) {
	return nil, errors.New("err")
}
func (m *mockPostServiceErr) GetAuthorPosts(u uint) ([]PostResponse, error) {
	return nil, errors.New("err")
}
func (m *mockPostServiceErr) GetAuthorPostByID(id, u uint) (*PostDetailResponse, error) {
	return nil, errors.New("err")
}
func (m *mockPostServiceErr) GetAuthorPostStats(u uint) (*PostStatsResponse, error) {
	return nil, errors.New("err")
}
func (m *mockPostServiceErr) CreatePost(u uint, req *CreatePostRequest) (*PostDetailResponse, error) {
	return nil, errors.New("err")
}
func (m *mockPostServiceErr) UpdateAuthorPost(id, u uint, req *UpdatePostRequest) (*PostDetailResponse, error) {
	return nil, errors.New("err")
}
func (m *mockPostServiceErr) TogglePublish(id, u uint) (*PostResponse, error) {
	return nil, errors.New("err")
}
func (m *mockPostServiceErr) DeleteAuthorPost(id, u uint) error {
	return errors.New("err")
}
func (m *mockPostServiceErr) GetAllUsers() ([]AdminUserResponse, error) {
	return nil, errors.New("err")
}
func (m *mockPostServiceErr) GetAllPostsArchive() ([]PostResponse, error) {
	return nil, errors.New("err")
}
func (m *mockPostServiceErr) GetUserPosts(u uint) ([]PostResponse, error) {
	return nil, errors.New("err")
}
