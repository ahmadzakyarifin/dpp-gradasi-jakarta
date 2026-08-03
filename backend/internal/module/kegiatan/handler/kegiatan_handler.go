package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type KegiatanHandler struct {
	svc service.KegiatanService
}

func NewKegiatanHandler(svc service.KegiatanService) *KegiatanHandler {
	return &KegiatanHandler{svc: svc}
}

// GET /api/v1/kegiatan — publik (published only)
func (h *KegiatanHandler) List(c *gin.Context) {
	var q dto.KegiatanQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	resp, err := h.svc.GetPublished(c.Request.Context(), q)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_LIST", "Daftar kegiatan berhasil diambil.", resp, nil)
}

// GET /api/v1/kegiatan/categories — publik (daftar kategori unik)
func (h *KegiatanHandler) GetCategories(c *gin.Context) {
	categories, err := h.svc.GetCategories(c.Request.Context())
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_CATEGORIES", "Daftar kategori kegiatan berhasil diambil.", categories, nil)
}

// GET /api/v1/kegiatan/admin — admin (all status)
func (h *KegiatanHandler) ListAdmin(c *gin.Context) {
	var q dto.KegiatanQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	resp, err := h.svc.GetAll(c.Request.Context(), q)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_LIST", "Daftar kegiatan berhasil diambil.", resp, nil)
}

// GET /api/v1/kegiatan/:slug — publik
func (h *KegiatanHandler) GetBySlug(c *gin.Context) {
	resp, err := h.svc.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_DETAIL", "Detail kegiatan berhasil diambil.", resp, nil)
}

// GET /api/v1/kegiatan/id/:id — admin (by ID)
func (h *KegiatanHandler) GetByID(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_DETAIL", "Detail kegiatan berhasil diambil.", resp, nil)
}

// POST /api/v1/kegiatan — admin
func (h *KegiatanHandler) Create(c *gin.Context) {
	var req dto.KegiatanCreateRequest
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
	helper.SuccessResponse(c, http.StatusCreated, "KEGIATAN_CREATED", "Kegiatan berhasil ditambahkan.", resp, nil)
}

// PUT /api/v1/kegiatan/:id — admin
func (h *KegiatanHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	var req dto.KegiatanUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	resp, err := h.svc.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_UPDATED", "Kegiatan berhasil diperbarui.", resp, nil)
}

// POST /api/v1/kegiatan/upload-image — admin (multipart field "image")
// Upload gambar cover/galeri kegiatan → path relatif /uploads/kegiatan/<file>
func (h *KegiatanHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		v := helper.NewValidationError()
		v.Add("image", "File gambar wajib diunggah.")
		helper.HandleServiceError(c, v)
		return
	}

	resp, err := h.svc.UploadImage(c.Request.Context(), file)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_IMAGE_UPLOADED", "Gambar berhasil diunggah", resp, nil)
}

// DELETE /api/v1/kegiatan/:id — admin (soft delete)
func (h *KegiatanHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_DELETED", "Kegiatan berhasil dihapus.", nil, nil)
}

// POST /api/v1/kegiatan/:id/restore — admin
func (h *KegiatanHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.svc.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_RESTORED", "Kegiatan berhasil dipulihkan.", nil, nil)
}

// DELETE /api/v1/kegiatan/gallery/:gallery_id — admin
func (h *KegiatanHandler) DeleteGalleryImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("gallery_id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID galeri tidak valid", nil)
		return
	}
	if err := h.svc.DeleteGalleryImage(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "GALLERY_IMAGE_DELETED", "Foto galeri berhasil dihapus.", nil, nil)
}

// POST /api/v1/kegiatan/bulk-delete — admin (bulk soft delete)
func (h *KegiatanHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	if err := h.svc.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_BULK_DELETED", "Kegiatan berhasil dihapus massal.", nil, nil)
}

// POST /api/v1/kegiatan/bulk-restore — admin (bulk restore)
func (h *KegiatanHandler) BulkRestore(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	if err := h.svc.BulkRestore(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_BULK_RESTORED", "Kegiatan berhasil dipulihkan massal.", nil, nil)
}
