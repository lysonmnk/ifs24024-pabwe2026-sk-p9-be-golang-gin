package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		Success(c, http.StatusOK, "Berhasil", map[string]string{"foo": "bar"})

		assert.Equal(t, http.StatusOK, w.Code)
		var resp Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		assert.Equal(t, "Berhasil", resp.Message)
	})

	t.Run("Error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		Error(c, http.StatusBadRequest, "Gagal", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.Equal(t, "Gagal", resp.Message)
	})
}
