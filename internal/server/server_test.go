package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"app/internal/config"
	"app/internal/features/ai"
	"app/internal/features/auth"
	"app/internal/features/post"
	"app/internal/middleware"

	"github.com/stretchr/testify/assert"
)

type mockAuthService struct {
	auth.Service
}

func (m *mockAuthService) Login(req *auth.LoginRequest) (*auth.AuthResponse, error) {
	return &auth.AuthResponse{Token: "test-token", TokenType: "Bearer"}, nil
}
func (m *mockAuthService) GetProfile(userID uint) (*auth.UserProfileResponse, error) {
	return &auth.UserProfileResponse{ID: userID, Name: "User", Role: "author"}, nil
}

type mockPostService struct {
	post.Service
}

func (m *mockPostService) GetCategories() ([]string, error) {
	return []string{"Teknologi"}, nil
}
func (m *mockPostService) GetAuthorPosts(userID uint) ([]post.PostResponse, error) {
	return []post.PostResponse{}, nil
}
func (m *mockPostService) GetAllUsers() ([]post.AdminUserResponse, error) {
	return []post.AdminUserResponse{}, nil
}

func TestServerRouter(t *testing.T) {
	cfg := &config.Config{
		AppEnv:  "production",
		AppPort: "8081",
		AppKey:  "test-secret-key-32-chars-long-minimum",
	}

	authH := auth.NewHandler(&mockAuthService{})
	postH := post.NewHandler(&mockPostService{})
	aiH := ai.NewHandler(nil)

	r := SetupRouter(cfg, authH, postH, aiH)
	assert.NotNil(t, r)

	t.Run("404 Endpoint not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/random-route", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Endpoint not found")
	})

	t.Run("Swagger and api.json endpoints", func(t *testing.T) {
		wAPI := httptest.NewRecorder()
		reqAPI, _ := http.NewRequest(http.MethodGet, "/api.json", nil)
		r.ServeHTTP(wAPI, reqAPI)
		assert.Equal(t, http.StatusMovedPermanently, wAPI.Code)

		wSwag := httptest.NewRecorder()
		reqSwag, _ := http.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		reqSwag.RequestURI = "/swagger/index.html"
		r.ServeHTTP(wSwag, reqSwag)
		assert.Equal(t, http.StatusOK, wSwag.Code)
	})

	t.Run("Public categories endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/public/categories", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Teknologi")
	})

	t.Run("Author posts endpoint with valid token", func(t *testing.T) {
		token, err := middleware.GenerateToken(1, "author", cfg.AppKey)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/posts", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Admin users endpoint with admin token", func(t *testing.T) {
		adminToken, err := middleware.GenerateToken(1, "admin", cfg.AppKey)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/admin/users", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("SetupRouter with development env", func(t *testing.T) {
		devCfg := &config.Config{
			AppEnv: "development",
			AppKey: "test-secret-key-32-chars-long-minimum",
		}
		devR := SetupRouter(devCfg, authH, postH, aiH)
		assert.NotNil(t, devR)
	})
}
