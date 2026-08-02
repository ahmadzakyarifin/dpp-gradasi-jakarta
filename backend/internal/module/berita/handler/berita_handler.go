package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/service"
	internalvalidator "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type BeritaHandler struct {
	svc    service.BeritaService
	logSvc activitylogservice.ActivityLogService
}

func NewBeritaHandler(svc service.BeritaService, logSvc activitylogservice.ActivityLogService) *BeritaHandler {
	return &BeritaHandler{svc: svc, logSvc: logSvc}
}

// bindErrors memetakan error binding gin (validator) menjadi per-field ValidationErrorItem
func (h *BeritaHandler) bindErrors(c *gin.Context, err error) {
	vitems := internalvalidator.Errors(err)
	if len(vitems) == 0 {
		vitems = []internalvalidator.ValidationErrorItem{{Field: "input", Tag: "invalid", Message: err.Error()}}
	}
	items := make([]helper.ValidationErrorItem, 0, len(vitems))
	for _, v := range vitems {
		items = append(items, helper.ValidationErrorItem{Field: v.Field, Tag: v.Tag, Param: v.Param, Message: v.Message})
	}
	helper.ValidationErrorResponse(c, items)
}

// GET /api/v1/berita — publik (published only)
func (h *BeritaHandler) List(c *gin.Context) {
	var q dto.BeritaQuery
	c.ShouldBindQuery(&q)

	resp, err := h.svc.GetPublished(q)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_LIST", "Daftar berita berhasil diambil.", resp, nil)
}

// GET /api/v1/berita/categories — publik (daftar kategori unik)
func (h *BeritaHandler) GetCategories(c *gin.Context) {
	categories, err := h.svc.GetCategories()
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Gagal mengambil daftar kategori.", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_CATEGORIES", "Daftar kategori berita berhasil diambil.", categories, nil)
}

// GET /api/v1/berita/admin — admin (all status)
func (h *BeritaHandler) ListAdmin(c *gin.Context) {
	var q dto.BeritaQuery
	c.ShouldBindQuery(&q)

	resp, err := h.svc.GetAll(q)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_LIST", "Daftar berita berhasil diambil.", resp, nil)
}

// GET /api/v1/berita/:slug — publik (detail by slug)
func (h *BeritaHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	resp, err := h.svc.GetBySlug(slug)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		httpCode := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			httpCode = http.StatusNotFound
		}
		helper.ErrorResponse(c, httpCode, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_DETAIL", "Detail berita berhasil diambil.", resp, nil)
}

// GET /api/v1/berita/id/:id — admin (by ID)
func (h *BeritaHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID berita tidak valid.", nil)
		return
	}
	resp, err := h.svc.GetByID(uint(id))
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		httpCode := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			httpCode = http.StatusNotFound
		}
		helper.ErrorResponse(c, httpCode, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_DETAIL", "Detail berita berhasil diambil.", resp, nil)
}

// mapServiceError mengubah ServiceError menjadi HTTP status code yang tepat
func mapServiceError(err error) (int, string, string) {
	svcErr, ok := err.(*helper.ServiceError)
	if !ok || svcErr == nil {
		return http.StatusInternalServerError, "SERVER_ERROR", "Terjadi kesalahan pada server."
	}
	switch svcErr.Code {
	case "NOT_FOUND":
		return http.StatusNotFound, svcErr.Code, svcErr.Message
	case "VALIDATION_ERROR", "DUPLICATE_TITLE":
		return http.StatusUnprocessableEntity, svcErr.Code, svcErr.Message
	default:
		return http.StatusInternalServerError, svcErr.Code, svcErr.Message
	}
}

// POST /api/v1/berita — admin
func (h *BeritaHandler) Create(c *gin.Context) {
	var req dto.BeritaCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.bindErrors(c, err)
		return
	}
	userID, _ := middleware.GetUserID(c)
	resp, err := h.svc.Create(&req, userID)
	if err != nil {
		httpCode, code, msg := mapServiceError(err)
		helper.ErrorResponse(c, httpCode, code, msg, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "BERITA_CREATED", "Berita berhasil diterbitkan.", resp, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "berita.create",
		EntityType:  "berita",
		EntityID:    &resp.ID,
		EntityLabel: resp.Title,
		Description: "Membuat berita baru",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// PUT /api/v1/berita/:id — admin
func (h *BeritaHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID berita tidak valid.", nil)
		return
	}
	var req dto.BeritaUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.bindErrors(c, err)
		return
	}
	resp, err := h.svc.Update(uint(id), &req)
	if err != nil {
		httpCode, code, msg := mapServiceError(err)
		helper.ErrorResponse(c, httpCode, code, msg, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_UPDATED", "Berita berhasil diperbarui.", resp, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "berita.update",
		EntityType:  "berita",
		EntityID:    &resp.ID,
		EntityLabel: resp.Title,
		Description: "Memperbarui berita",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// DELETE /api/v1/berita/:id — admin (soft delete)
func (h *BeritaHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID berita tidak valid.", nil)
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		httpCode := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			httpCode = http.StatusNotFound
		}
		helper.ErrorResponse(c, httpCode, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_DELETED", "Berita berhasil dihapus.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(id)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "berita.delete",
		EntityType:  "berita",
		EntityID:    &eID,
		Description: "Menghapus berita",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/berita/:id/restore — admin
func (h *BeritaHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID berita tidak valid.", nil)
		return
	}
	if err := h.svc.Restore(uint(id)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "BERITA_RESTORED", "Berita berhasil dipulihkan.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(id)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "berita.restore",
		EntityType:  "berita",
		EntityID:    &eID,
		Description: "Memulihkan berita",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/berita/bulk-delete — admin (bulk soft delete)
func (h *BeritaHandler) BulkDelete(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "BERITA_BULK_DELETED", "Berita berhasil dihapus massal.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "berita.bulk_delete",
		EntityType:  "berita",
		Description: "Menghapus berita secara massal",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/berita/bulk-restore — admin (bulk restore)
func (h *BeritaHandler) BulkRestore(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "BERITA_BULK_RESTORED", "Berita berhasil dipulihkan massal.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "berita.bulk_restore",
		EntityType:  "berita",
		Description: "Memulihkan berita secara massal",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}
