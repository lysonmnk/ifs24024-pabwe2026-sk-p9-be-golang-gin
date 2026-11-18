package auth

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

// Register godoc
// @Summary Pendaftaran akun penulis baru
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Data Registrasi"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data registrasi tidak valid: "+err.Error(), nil)
		return
	}

	res, err := h.service.Register(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Registrasi berhasil.", res)
}

// Login godoc
// @Summary Autentikasi dan penerbitan token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Kredensial Login"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Email dan password wajib diisi", nil)
		return
	}

	if req.Email == "" && req.Username == "" {
		response.Error(c, http.StatusBadRequest, "Email dan password wajib diisi", nil)
		return
	}

	res, err := h.service.Login(&req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil.", res)
}

// GetProfile godoc
// @Summary Mengambil profil penulis yang sedang login
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/auth/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Tidak terautentikasi", nil)
		return
	}

	userID := val.(uint)
	res, err := h.service.GetProfile(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Profil berhasil dimuat.", res)
}

// UpdateProfile godoc
// @Summary Memperbarui profil penulis
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Data Pembaruan Profil"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/auth/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Tidak terautentikasi", nil)
		return
	}

	userID := val.(uint)
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data pembaruan tidak valid: "+err.Error(), nil)
		return
	}

	res, err := h.service.UpdateProfile(userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Profil berhasil diperbarui.", res)
}

// Logout godoc
// @Summary Mengakhiri sesi pengguna
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	response.Success(c, http.StatusOK, "Logout berhasil.", nil)
}
