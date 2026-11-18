package ai

import (
	"net/http"

	"app/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GeneratePost godoc
// @Summary Menghasilkan draf artikel dengan AI dan penelusuran Search Engine
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GenerateRequest true "Parameter Pembuatan Artikel AI"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/posts/ai-generate [post]
func (h *Handler) GeneratePost(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Permintaan AI generate tidak valid: "+err.Error(), nil)
		return
	}

	res, err := h.service.GeneratePost(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal menghasilkan konten AI: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Draf artikel AI berhasil digenerate.", res)
}
