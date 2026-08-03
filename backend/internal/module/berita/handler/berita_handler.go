package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type BeritaHandler struct {
	svc service.BeritaService
}

func NewBeritaHandler(svc service.BeritaService) *BeritaHandler {
	return &BeritaHandler{svc: svc}
}

// GET /api/v1/berita — publik (published only)
func (h *BeritaHandler) List(c *gin.Context) {
	var q dto.BeritaQuery
	_ = c.ShouldBindQuery(&q)

	resp, err := h.svc.GetPublished(c.Request.Context(), q)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_LIST", "Daftar berita berhasil diambil.", resp, nil)
}

// GET /api/v1/berita/categories — publik (daftar kategori unik)
func (h *BeritaHandler) GetCategories(c *gin.Context) {
	categories, err := h.svc.GetCategories(c.Request.Context())
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_CATEGORIES", "Daftar kategori berita berhasil diambil.", categories, nil)
}

// GET /api/v1/berita/admin — admin (all status)
func (h *BeritaHandler) ListAdmin(c *gin.Context) {
	var q dto.BeritaQuery
	_ = c.ShouldBindQuery(&q)

	resp, err := h.svc.GetAll(c.Request.Context(), q)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_LIST", "Daftar berita berhasil diambil.", resp, nil)
}

// GET /api/v1/berita/:slug — publik (detail by slug)
func (h *BeritaHandler) GetBySlug(c *gin.Context) {
	resp, err := h.svc.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_DETAIL", "Detail berita berhasil diambil.", resp, nil)
}

// GET /api/v1/berita/id/:id — admin (by ID)
func (h *BeritaHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	resp, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_DETAIL", "Detail berita berhasil diambil.", resp, nil)
}

// POST /api/v1/berita — admin
func (h *BeritaHandler) Create(c *gin.Context) {
	var req dto.BeritaCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	userID, _ := middleware.GetUserID(c)
	resp, err := h.svc.Create(c.Request.Context(), &req, userID)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "BERITA_CREATED", "Berita berhasil diterbitkan.", resp, nil)
}

// POST /api/v1/berita/upload-image — admin (multipart field "image")
// Upload gambar cover berita → path relatif /uploads/berita/<file>
func (h *BeritaHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "File gambar wajib diunggah (field 'image')", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)
	_ = userID // upload tidak mencatat audit per-user (hanya simpan file)
	resp, err := h.svc.UploadImage(c.Request.Context(), file)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_IMAGE_UPLOADED", "Gambar berhasil diunggah.", resp, nil)
}

// PUT /api/v1/berita/:id — admin
func (h *BeritaHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	var req dto.BeritaUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	resp, err := h.svc.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_UPDATED", "Berita berhasil diperbarui.", resp, nil)
}

// DELETE /api/v1/berita/:id — admin (soft delete)
func (h *BeritaHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_DELETED", "Berita berhasil dihapus.", nil, nil)
}

// POST /api/v1/berita/:id/restore — admin
func (h *BeritaHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.svc.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_RESTORED", "Berita berhasil dipulihkan.", nil, nil)
}

// POST /api/v1/berita/bulk-delete — admin (bulk soft delete)
func (h *BeritaHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	if err := h.svc.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_BULK_DELETED", "Berita berhasil dihapus massal.", nil, nil)
}

// POST /api/v1/berita/bulk-restore — admin (bulk restore)
func (h *BeritaHandler) BulkRestore(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	if err := h.svc.BulkRestore(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_BULK_RESTORED", "Berita berhasil dipulihkan massal.", nil, nil)
}
