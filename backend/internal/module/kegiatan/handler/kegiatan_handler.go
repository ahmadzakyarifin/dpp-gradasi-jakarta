package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/service"
	"github.com/gin-gonic/gin"
)

type KegiatanHandler struct {
	svc    service.KegiatanService
	logSvc activitylogservice.ActivityLogService
}

func NewKegiatanHandler(svc service.KegiatanService, logSvc activitylogservice.ActivityLogService) *KegiatanHandler {
	return &KegiatanHandler{svc: svc, logSvc: logSvc}
}

// GET /api/v1/kegiatan — publik (published only)
func (h *KegiatanHandler) List(c *gin.Context) {
	var q dto.KegiatanQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Query tidak valid.", nil)
		return
	}

	resp, err := h.svc.GetPublished(q)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_LIST", "Daftar kegiatan berhasil diambil.", resp, nil)
}

// GET /api/v1/kegiatan/categories — publik (daftar kategori unik)
func (h *KegiatanHandler) GetCategories(c *gin.Context) {
	categories, err := h.svc.GetCategories()
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Gagal mengambil daftar kategori.", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_CATEGORIES", "Daftar kategori kegiatan berhasil diambil.", categories, nil)
}

// GET /api/v1/kegiatan/admin — admin (all status)
func (h *KegiatanHandler) ListAdmin(c *gin.Context) {
	var q dto.KegiatanQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Query tidak valid.", nil)
		return
	}

	resp, err := h.svc.GetAll(q)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_LIST", "Daftar kegiatan berhasil diambil.", resp, nil)
}

// GET /api/v1/kegiatan/:slug — publik
func (h *KegiatanHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	resp, err := h.svc.GetBySlug(slug)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		code := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			code = http.StatusNotFound
		}
		helper.ErrorResponse(c, code, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_DETAIL", "Detail kegiatan berhasil diambil.", resp, nil)
}

// GET /api/v1/kegiatan/id/:id — admin (by ID)
func (h *KegiatanHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID tidak valid.", nil)
		return
	}
	resp, err := h.svc.GetByID(uint(id))
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		code := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			code = http.StatusNotFound
		}
		helper.ErrorResponse(c, code, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_DETAIL", "Detail kegiatan berhasil diambil.", resp, nil)
}

// POST /api/v1/kegiatan — admin
func (h *KegiatanHandler) Create(c *gin.Context) {
	var req dto.KegiatanCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "title", Tag: "required", Message: "Title dan content wajib diisi."},
		})
		return
	}
	userID, _ := middleware.GetUserID(c)
	resp, err := h.svc.Create(&req, userID)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "KEGIATAN_CREATED", "Kegiatan berhasil ditambahkan.", resp, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "kegiatan.create",
		EntityType:  "kegiatan",
		EntityID:    &resp.ID,
		EntityLabel: resp.Title,
		Description: "Membuat kegiatan baru",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// PUT /api/v1/kegiatan/:id — admin
func (h *KegiatanHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID tidak valid.", nil)
		return
	}
	var req dto.KegiatanUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "title", Tag: "required", Message: "Title dan content wajib diisi."},
		})
		return
	}
	resp, err := h.svc.Update(uint(id), &req)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		code := http.StatusInternalServerError
		switch svcErr.Code {
		case "NOT_FOUND":
			code = http.StatusNotFound
		case "VALIDATION_ERROR", "DUPLICATE_TITLE":
			code = http.StatusUnprocessableEntity
		}
		helper.ErrorResponse(c, code, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_UPDATED", "Kegiatan berhasil diperbarui.", resp, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "kegiatan.update",
		EntityType:  "kegiatan",
		EntityID:    &resp.ID,
		EntityLabel: resp.Title,
		Description: "Memperbarui kegiatan",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// DELETE /api/v1/kegiatan/:id — admin (soft delete)
func (h *KegiatanHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID tidak valid.", nil)
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		code := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			code = http.StatusNotFound
		}
		helper.ErrorResponse(c, code, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_DELETED", "Kegiatan berhasil dihapus.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(id)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "kegiatan.delete",
		EntityType:  "kegiatan",
		EntityID:    &eID,
		Description: "Menghapus kegiatan",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/kegiatan/:id/restore — admin
func (h *KegiatanHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID tidak valid.", nil)
		return
	}
	if err := h.svc.Restore(uint(id)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_RESTORED", "Kegiatan berhasil dipulihkan.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(id)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "kegiatan.restore",
		EntityType:  "kegiatan",
		EntityID:    &eID,
		Description: "Memulihkan kegiatan",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// DELETE /api/v1/kegiatan/gallery/:gallery_id — admin
func (h *KegiatanHandler) DeleteGalleryImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("gallery_id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID galeri tidak valid.", nil)
		return
	}
	if err := h.svc.DeleteGalleryImage(uint(id)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Gagal menghapus foto galeri.", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "GALLERY_IMAGE_DELETED", "Foto galeri berhasil dihapus.", nil, nil)
}

// POST /api/v1/kegiatan/bulk-delete — admin (bulk soft delete)
func (h *KegiatanHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "ids", Tag: "required", Message: "IDs wajib diisi minimal 1."},
		})
		return
	}
	if err := h.svc.BulkDelete(req.IDs); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_BULK_DELETED", "Kegiatan berhasil dihapus massal.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "kegiatan.bulk_delete",
		EntityType:  "kegiatan",
		Description: "Menghapus kegiatan secara massal",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/kegiatan/bulk-restore — admin (bulk restore)
func (h *KegiatanHandler) BulkRestore(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "ids", Tag: "required", Message: "IDs wajib diisi minimal 1."},
		})
		return
	}
	if err := h.svc.BulkRestore(req.IDs); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KEGIATAN_BULK_RESTORED", "Kegiatan berhasil dipulihkan massal.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "kegiatan.bulk_restore",
		EntityType:  "kegiatan",
		Description: "Memulihkan kegiatan secara massal",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}
