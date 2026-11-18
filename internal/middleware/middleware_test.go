package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret-key-32-chars-long-minimum"

	t.Run("Generate and validate valid token", func(t *testing.T) {
		token, err := GenerateToken(1, "author", secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := ValidateToken(token, secret)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), claims.UserID)
		assert.Equal(t, "author", claims.Role)
	})

	t.Run("Validate token with invalid secret or malformed", func(t *testing.T) {
		_, err := ValidateToken("invalid.token.string", secret)
		assert.Error(t, err)

		token, err := GenerateToken(1, "author", secret)
		assert.NoError(t, err)
		_, err = ValidateToken(token, "wrong-secret")
		assert.Error(t, err)
	})

	t.Run("Validate token with unexpected signing method", func(t *testing.T) {
		claims := JWTClaims{
			UserID: 1,
			Role:   "author",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		assert.NoError(t, err)

		_, err = ValidateToken(tokenString, secret)
		assert.Error(t, err)
	})

	t.Run("AuthMiddleware missing header", func(t *testing.T) {
		r := gin.New()
		r.Use(AuthMiddleware(secret))
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("AuthMiddleware invalid header format", func(t *testing.T) {
		r := gin.New()
		r.Use(AuthMiddleware(secret))
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Basic 12345")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("AuthMiddleware invalid token string", func(t *testing.T) {
		r := gin.New()
		r.Use(AuthMiddleware(secret))
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("AuthMiddleware success", func(t *testing.T) {
		token, err := GenerateToken(1, "author", secret)
		assert.NoError(t, err)

		r := gin.New()
		r.Use(AuthMiddleware(secret))
		r.GET("/test", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			role, _ := c.Get("role")
			c.JSON(http.StatusOK, gin.H{"userID": userID, "role": role})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("RequireRole role not found in context", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole("admin"))
		r.GET("/admin", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("RequireRole role in context is not string", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("role", 12345)
			c.Next()
		})
		r.Use(RequireRole("admin"))
		r.GET("/admin", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("RequireRole role not allowed", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("role", "author")
			c.Next()
		})
		r.Use(RequireRole("admin"))
		r.GET("/admin", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("RequireRole success", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("role", "admin")
			c.Next()
		})
		r.Use(RequireRole("admin"))
		r.GET("/admin", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(CORSMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	t.Run("Options preflight request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodOptions, "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("Standard GET request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewRateLimiter(2, 50*time.Millisecond)
	r := gin.New()
	r.Use(limiter.Middleware())
	r.GET("/rate", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Request 1: OK
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodGet, "/rate", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Request 2: OK
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/rate", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Request 3: Exceeded -> 429
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest(http.MethodGet, "/rate", nil)
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusTooManyRequests, w3.Code)

	// Wait for reset window
	time.Sleep(60 * time.Millisecond)

	// Request 4: Reset -> OK
	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest(http.MethodGet, "/rate", nil)
	r.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
}
