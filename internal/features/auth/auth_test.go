package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"app/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)
	err = db.AutoMigrate(&database.User{}, &database.Post{})
	assert.NoError(t, err)
	return db
}

func TestAuthRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &database.User{
		Name:     "Test User",
		Username: "testuser",
		Email:    "test@contoh.com",
		Password: string(hashed),
		Role:     "author",
	}

	err := repo.CreateUser(user)
	assert.NoError(t, err)

	t.Run("FindByEmail", func(t *testing.T) {
		u, err := repo.FindByEmail("test@contoh.com")
		assert.NoError(t, err)
		assert.Equal(t, "testuser", u.Username)

		_, err = repo.FindByEmail("notfound@contoh.com")
		assert.Error(t, err)
	})

	t.Run("FindByID", func(t *testing.T) {
		u, err := repo.FindByID(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "test@contoh.com", u.Email)

		_, err = repo.FindByID(99999)
		assert.Error(t, err)
	})

	t.Run("FindByUsername", func(t *testing.T) {
		u, err := repo.FindByUsername("testuser")
		assert.NoError(t, err)
		assert.Equal(t, user.ID, u.ID)

		_, err = repo.FindByUsername("unknown")
		assert.Error(t, err)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		user.Name = "Updated Name"
		err := repo.UpdateUser(user)
		assert.NoError(t, err)

		u, _ := repo.FindByID(user.ID)
		assert.Equal(t, "Updated Name", u.Name)
	})
}

func TestAuthService(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	secret := "test-secret-key-32-chars-long-minimum"
	svc := NewService(repo, secret)

	t.Run("Register success", func(t *testing.T) {
		req := &RegisterRequest{
			Name:     "Author One",
			Username: "author1",
			Email:    "author1@contoh.com",
			Password: "password123",
		}
		res, err := svc.Register(req)
		assert.NoError(t, err)
		assert.Equal(t, "author1", res.Username)
	})

	t.Run("Register default username from email", func(t *testing.T) {
		req := &RegisterRequest{
			Name:     "Author Two",
			Username: "",
			Email:    "author2@contoh.com",
			Password: "password123",
		}
		res, err := svc.Register(req)
		assert.NoError(t, err)
		assert.Equal(t, "author2", res.Username)
	})

	t.Run("Register duplicate email error", func(t *testing.T) {
		req := &RegisterRequest{
			Name:     "Author Duplicate",
			Username: "author_dup",
			Email:    "author1@contoh.com",
			Password: "password123",
		}
		_, err := svc.Register(req)
		assert.Error(t, err)
		assert.Equal(t, "email sudah terdaftar", err.Error())
	})

	t.Run("Register duplicate username error", func(t *testing.T) {
		req := &RegisterRequest{
			Name:     "Author Duplicate Username",
			Username: "author1",
			Email:    "other@contoh.com",
			Password: "password123",
		}
		_, err := svc.Register(req)
		assert.Error(t, err)
		assert.Equal(t, "username sudah digunakan", err.Error())
	})

	t.Run("Register bcrypt password too long error", func(t *testing.T) {
		req := &RegisterRequest{
			Name:     "Author Long",
			Username: "author_long",
			Email:    "long@contoh.com",
			Password: strings.Repeat("a", 100),
		}
		_, err := svc.Register(req)
		assert.Error(t, err)
	})

	t.Run("Login success and failure", func(t *testing.T) {
		loginReq := &LoginRequest{
			Email:    "author1@contoh.com",
			Password: "password123",
		}
		res, err := svc.Login(loginReq)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Token)

		// Wrong password
		loginReq.Password = "wrongpass"
		_, err = svc.Login(loginReq)
		assert.Error(t, err)
		assert.Equal(t, "email atau password salah", err.Error())

		// Login with username
		loginUserReq := &LoginRequest{
			Username: "author1",
			Password: "password123",
		}
		resUser, err := svc.Login(loginUserReq)
		assert.NoError(t, err)
		assert.NotEmpty(t, resUser.Token)

		// Login with username in email field (no @)
		loginEmailAsUserReq := &LoginRequest{
			Email:    "author1",
			Password: "password123",
		}
		resEmailAsUser, err := svc.Login(loginEmailAsUserReq)
		assert.NoError(t, err)
		assert.NotEmpty(t, resEmailAsUser.Token)

		// Login with non-existent username
		loginNonExistentUserReq := &LoginRequest{
			Username: "nonexistent_user",
			Password: "password123",
		}
		_, err = svc.Login(loginNonExistentUserReq)
		assert.Error(t, err)
		assert.Equal(t, "email atau password salah", err.Error())
	})

	t.Run("GetProfile and UpdateProfile", func(t *testing.T) {
		user, err := repo.FindByEmail("author1@contoh.com")
		assert.NoError(t, err)

		prof, err := svc.GetProfile(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Author One", prof.Name)

		// Non-existent user profile
		_, err = svc.GetProfile(99999)
		assert.Error(t, err)

		// Update profile success
		upReq := &UpdateProfileRequest{
			Name:     "Author One Renamed",
			Username: "author1_renamed",
		}
		updated, err := svc.UpdateProfile(user.ID, upReq)
		assert.NoError(t, err)
		assert.Equal(t, "Author One Renamed", updated.Name)

		// Update profile with same username
		sameUserReq := &UpdateProfileRequest{
			Name:     "Author One Again",
			Username: "author1_renamed",
		}
		_, err = svc.UpdateProfile(user.ID, sameUserReq)
		assert.NoError(t, err)

		// Update profile username taken by another user
		upConflictReq := &UpdateProfileRequest{
			Name:     "Author One",
			Username: "author2",
		}
		_, err = svc.UpdateProfile(user.ID, upConflictReq)
		assert.Error(t, err)
		assert.Equal(t, "username sudah digunakan", err.Error())

		// Update non-existent user
		_, err = svc.UpdateProfile(99999, upReq)
		assert.Error(t, err)
	})

	t.Run("CreateUser and UpdateUser db errors on closed db", func(t *testing.T) {
		closedDB := setupTestDB(t)
		sqlDB, _ := closedDB.DB()
		_ = sqlDB.Close()
		closedRepo := NewRepository(closedDB)
		closedSvc := NewService(closedRepo, secret)

		req := &RegisterRequest{
			Name:     "Err User",
			Username: "erruser",
			Email:    "err@contoh.com",
			Password: "password123",
		}
		_, err := closedSvc.Register(req)
		assert.Error(t, err)

		// Mock repo to trigger UpdateUser error in UpdateProfile
		upReq := &UpdateProfileRequest{Name: "Err", Username: "err"}
		mockRepo := &mockAuthRepoUpdateErr{repo}
		mockSvc := NewService(mockRepo, secret)
		_, err = mockSvc.UpdateProfile(1, upReq)
		assert.Error(t, err)
	})
}

type mockAuthRepoUpdateErr struct {
	Repository
}

func (m *mockAuthRepoUpdateErr) FindByID(id uint) (*database.User, error) {
	return &database.User{ID: id, Name: "Old", Username: "old", Email: "old@c.com"}, nil
}

func (m *mockAuthRepoUpdateErr) FindByUsername(username string) (*database.User, error) {
	return nil, nil
}

func (m *mockAuthRepoUpdateErr) UpdateUser(user *database.User) error {
	return errors.New("update db error")
}

func TestAuthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	repo := NewRepository(db)
	secret := "test-secret-key-32-chars-long-minimum"
	svc := NewService(repo, secret)
	handler := NewHandler(svc)

	t.Run("Register handler", func(t *testing.T) {
		// Valid
		body, _ := json.Marshal(RegisterRequest{
			Name:     "Handler User",
			Username: "handleruser",
			Email:    "handler@contoh.com",
			Password: "password123",
		})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.Register(c)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Invalid JSON
		wBad := httptest.NewRecorder()
		cBad, _ := gin.CreateTestContext(wBad)
		cBad.Request, _ = http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer([]byte("invalid json")))
		cBad.Request.Header.Set("Content-Type", "application/json")
		handler.Register(cBad)
		assert.Equal(t, http.StatusBadRequest, wBad.Code)

		// Service error (duplicate)
		wDup := httptest.NewRecorder()
		cDup, _ := gin.CreateTestContext(wDup)
		cDup.Request, _ = http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
		cDup.Request.Header.Set("Content-Type", "application/json")
		handler.Register(cDup)
		assert.Equal(t, http.StatusBadRequest, wDup.Code)
	})

	t.Run("Login handler", func(t *testing.T) {
		// Valid
		body, _ := json.Marshal(LoginRequest{
			Email:    "handler@contoh.com",
			Password: "password123",
		})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.Login(c)
		assert.Equal(t, http.StatusOK, w.Code)

		// Invalid JSON
		wBad := httptest.NewRecorder()
		cBad, _ := gin.CreateTestContext(wBad)
		cBad.Request, _ = http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer([]byte("{}")))
		cBad.Request.Header.Set("Content-Type", "application/json")
		handler.Login(cBad)
		assert.Equal(t, http.StatusBadRequest, wBad.Code)

		// Wrong password -> 401
		bodyWrong, _ := json.Marshal(LoginRequest{
			Email:    "handler@contoh.com",
			Password: "wrongpassword",
		})
		wWrong := httptest.NewRecorder()
		cWrong, _ := gin.CreateTestContext(wWrong)
		cWrong.Request, _ = http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(bodyWrong))
		cWrong.Request.Header.Set("Content-Type", "application/json")
		handler.Login(cWrong)
		assert.Equal(t, http.StatusUnauthorized, wWrong.Code)

		// Missing both email and username -> 400
		bodyEmpty, _ := json.Marshal(LoginRequest{
			Password: "password123",
		})
		wEmpty := httptest.NewRecorder()
		cEmpty, _ := gin.CreateTestContext(wEmpty)
		cEmpty.Request, _ = http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(bodyEmpty))
		cEmpty.Request.Header.Set("Content-Type", "application/json")
		handler.Login(cEmpty)
		assert.Equal(t, http.StatusBadRequest, wEmpty.Code)

		// Valid login with username -> 200
		bodyUser, _ := json.Marshal(LoginRequest{
			Username: "handleruser",
			Password: "password123",
		})
		wUser := httptest.NewRecorder()
		cUser, _ := gin.CreateTestContext(wUser)
		cUser.Request, _ = http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(bodyUser))
		cUser.Request.Header.Set("Content-Type", "application/json")
		handler.Login(cUser)
		assert.Equal(t, http.StatusOK, wUser.Code)
	})

	t.Run("GetProfile handler", func(t *testing.T) {
		user, _ := repo.FindByEmail("handler@contoh.com")

		// Valid
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", user.ID)
		handler.GetProfile(c)
		assert.Equal(t, http.StatusOK, w.Code)

		// Missing context
		wNoAuth := httptest.NewRecorder()
		cNoAuth, _ := gin.CreateTestContext(wNoAuth)
		handler.GetProfile(cNoAuth)
		assert.Equal(t, http.StatusUnauthorized, wNoAuth.Code)

		// Not found
		wNotFound := httptest.NewRecorder()
		cNotFound, _ := gin.CreateTestContext(wNotFound)
		cNotFound.Set("userID", uint(99999))
		handler.GetProfile(cNotFound)
		assert.Equal(t, http.StatusNotFound, wNotFound.Code)
	})

	t.Run("UpdateProfile handler", func(t *testing.T) {
		user, _ := repo.FindByEmail("handler@contoh.com")

		// Valid
		body, _ := json.Marshal(UpdateProfileRequest{
			Name:     "Handler User Updated",
			Username: "handler_updated",
		})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", user.ID)
		c.Request, _ = http.NewRequest(http.MethodPut, "/api/auth/profile", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateProfile(c)
		assert.Equal(t, http.StatusOK, w.Code)

		// Missing context
		wNoAuth := httptest.NewRecorder()
		cNoAuth, _ := gin.CreateTestContext(wNoAuth)
		handler.UpdateProfile(cNoAuth)
		assert.Equal(t, http.StatusUnauthorized, wNoAuth.Code)

		// Invalid JSON
		wBad := httptest.NewRecorder()
		cBad, _ := gin.CreateTestContext(wBad)
		cBad.Set("userID", user.ID)
		cBad.Request, _ = http.NewRequest(http.MethodPut, "/api/auth/profile", bytes.NewBuffer([]byte("bad")))
		cBad.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateProfile(cBad)
		assert.Equal(t, http.StatusBadRequest, wBad.Code)

		// Service error
		wErr := httptest.NewRecorder()
		cErr, _ := gin.CreateTestContext(wErr)
		cErr.Set("userID", uint(99999))
		cErr.Request, _ = http.NewRequest(http.MethodPut, "/api/auth/profile", bytes.NewBuffer(body))
		cErr.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateProfile(cErr)
		assert.Equal(t, http.StatusBadRequest, wErr.Code)
	})

	t.Run("Logout handler", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.Logout(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
