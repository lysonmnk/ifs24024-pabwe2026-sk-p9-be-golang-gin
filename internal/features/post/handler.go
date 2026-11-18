package post

import (
	"net/http"
	"strconv"

	"app/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetCategories godoc
// @Summary Daftar kategori artikel
// @Tags Public
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/public/categories [get]
func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Daftar kategori berhasil dimuat.", categories)
}

// GetPublicPosts godoc
// @Summary Daftar artikel terbit, mendukung filter kata kunci dan kategori
// @Tags Public
// @Produce json
// @Param q query string false "Kata Kunci Pencarian"
// @Param category query string false "Filter Kategori"
// @Success 200 {object} response.Response
// @Router /api/public/posts [get]
func (h *Handler) GetPublicPosts(c *gin.Context) {
	q := c.Query("q")
	category := c.Query("category")

	posts, err := h.service.GetPublicPosts(q, category)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Daftar artikel publik berhasil dimuat.", posts)
}

// GetPublicPostBySlug godoc
// @Summary Detail artikel terbit berdasarkan slug
// @Tags Public
// @Produce json
// @Param slug path string true "Slug Artikel"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/public/posts/{slug} [get]
func (h *Handler) GetPublicPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	post, err := h.service.GetPublicPostBySlug(slug)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Artikel tidak ditemukan.", nil)
		return
	}
	response.Success(c, http.StatusOK, "Detail artikel berhasil dimuat.", post)
}

// GetAuthorPosts godoc
// @Summary Daftar seluruh artikel milik penulis (terbit dan draf)
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/posts [get]
func (h *Handler) GetAuthorPosts(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Tidak terautentikasi", nil)
		return
	}

	userID := val.(uint)
	posts, err := h.service.GetAuthorPosts(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Daftar artikel penulis berhasil dimuat.", posts)
}

// GetAuthorPostStats godoc
// @Summary Statistik artikel penulis: total, terbit, dan draf
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/posts/stats [get]
func (h *Handler) GetAuthorPostStats(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Tidak terautentikasi", nil)
		return
	}

	userID := val.(uint)
	stats, err := h.service.GetAuthorPostStats(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Statistik artikel berhasil dimuat.", stats)
}

// CreatePost godoc
// @Summary Membuat artikel baru beserta slug unik dan estimasi waktu baca
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePostRequest true "Data Artikel Baru"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/posts [post]
func (h *Handler) CreatePost(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Tidak terautentikasi", nil)
		return
	}

	userID := val.(uint)
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data artikel tidak valid: "+err.Error(), nil)
		return
	}

	post, err := h.service.CreatePost(userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusCreated, "Artikel berhasil dibuat.", post)
}

// GetAuthorPostByID godoc
// @Summary Detail artikel milik penulis untuk keperluan penyuntingan
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Artikel"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/posts/{id} [get]
func (h *Handler) GetAuthorPostByID(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Tidak terautentikasi", nil)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	userID := val.(uint)
	post, err := h.service.GetAuthorPostByID(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Detail artikel berhasil dimuat.", post)
}

// UpdateAuthorPost godoc
// @Summary Memperbarui artikel milik penulis
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Artikel"
// @Param request body UpdatePostRequest true "Data Pembaruan Artikel"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/posts/{id} [put]
func (h *Handler) UpdateAuthorPost(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Tidak terautentikasi", nil)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	userID := val.(uint)
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data pembaruan tidak valid: "+err.Error(), nil)
		return
	}

	post, err := h.service.UpdateAuthorPost(uint(id), userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Artikel berhasil diperbarui.", post)
}

// TogglePublish godoc
// @Summary Mengubah status publikasi artikel (draf/terbit)
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Artikel"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/posts/{id}/toggle-publish [patch]
func (h *Handler) TogglePublish(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Tidak terautentikasi", nil)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	userID := val.(uint)
	post, err := h.service.TogglePublish(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Status publikasi berhasil diubah.", post)
}

// DeleteAuthorPost godoc
// @Summary Menghapus artikel milik penulis
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Artikel"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/posts/{id} [delete]
func (h *Handler) DeleteAuthorPost(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Tidak terautentikasi", nil)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	userID := val.(uint)
	if err := h.service.DeleteAuthorPost(uint(id), userID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Artikel berhasil dihapus.", nil)
}

// AdminGetAllUsers godoc
// @Summary Daftar seluruh pengguna beserta jumlah artikel
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/admin/users [get]
func (h *Handler) AdminGetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Daftar pengguna berhasil dimuat.", users)
}

// AdminGetAllPosts godoc
// @Summary Arsip seluruh artikel dari semua penulis
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/admin/posts [get]
func (h *Handler) AdminGetAllPosts(c *gin.Context) {
	posts, err := h.service.GetAllPostsArchive()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Arsip seluruh artikel berhasil dimuat.", posts)
}

// AdminGetUserPosts godoc
// @Summary Daftar artikel milik pengguna tertentu
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Pengguna"
// @Success 200 {object} response.Response
// @Router /api/admin/users/{id}/posts [get]
func (h *Handler) AdminGetUserPosts(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID pengguna tidak valid", nil)
		return
	}

	posts, err := h.service.GetUserPosts(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Daftar artikel pengguna berhasil dimuat.", posts)
}
