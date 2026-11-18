package middleware

import (
	"net/http"

	"app/internal/response"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusForbidden, "Akses ditolak. Peran pengguna tidak ditemukan.", nil)
			c.Abort()
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			response.Error(c, http.StatusForbidden, "Akses ditolak. Peran pengguna tidak valid.", nil)
			c.Abort()
			return
		}

		allowed := false
		for _, r := range allowedRoles {
			if userRole == r {
				allowed = true
				break
			}
		}

		if !allowed {
			response.Error(c, http.StatusForbidden, "Akses ditolak. Memerlukan hak akses administrator.", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
