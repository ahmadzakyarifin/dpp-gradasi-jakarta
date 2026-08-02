package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/service"
	"github.com/gin-gonic/gin"
)

type SlidersHandler struct {
	svc    service.SlidersService
	logSvc activitylogservice.ActivityLogService
}

func NewSlidersHandler(svc service.SlidersService, logSvc activitylogservice.ActivityLogService) *SlidersHandler {
	return &SlidersHandler{svc: svc, logSvc: logSvc}
}

// GET /api/v1/sliders — publik
func (h *SlidersHandler) GetAll(c *gin.Context) {
	// Publik SELALU hanya melihat slider aktif
	resp, err := h.svc.GetAll(true)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDERS_FETCHED", "Data slider berhasil diambil.", resp, nil)
}

// GET /api/v1/sliders/admin — admin
func (h *SlidersHandler) GetAllAdmin(c *gin.Context) {
	activeOnly := c.DefaultQuery("active", "true") == "true"
	resp, err := h.svc.GetAll(activeOnly)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDERS_FETCHED", "Data slider berhasil diambil.", resp, nil)
}

// GET /api/v1/sliders/:id — publik
func (h *SlidersHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID slider tidak valid.", nil)
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
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_FETCHED", "Data slider berhasil diambil.", resp, nil)
}

// POST /api/v1/sliders — admin
func (h *SlidersHandler) Create(c *gin.Context) {
	var req dto.SliderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "title", Tag: "required", Message: "Title dan image_path wajib diisi."},
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
	helper.SuccessResponse(c, http.StatusCreated, "SLIDER_CREATED", "Slider berhasil dibuat.", resp, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "slider.create",
		EntityType:  "slider",
		EntityID:    &resp.ID,
		EntityLabel: resp.Title,
		Description: "Membuat slider baru",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// PUT /api/v1/sliders/:id — admin
func (h *SlidersHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID slider tidak valid.", nil)
		return
	}
	var req dto.SliderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "title", Tag: "required", Message: "Title dan image_path wajib diisi."},
		})
		return
	}
	resp, err := h.svc.Update(uint(id), &req)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		httpCode := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			httpCode = http.StatusNotFound
		}
		helper.ErrorResponse(c, httpCode, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_UPDATED", "Slider berhasil diupdate.", resp, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "slider.update",
		EntityType:  "slider",
		EntityID:    &resp.ID,
		EntityLabel: resp.Title,
		Description: "Memperbarui slider",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// DELETE /api/v1/sliders/:id — admin
func (h *SlidersHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID slider tidak valid.", nil)
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
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_DELETED", "Slider berhasil dihapus.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(id)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "slider.delete",
		EntityType:  "slider",
		EntityID:    &eID,
		Description: "Menghapus slider",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// PUT /api/v1/sliders/reorder — admin
func (h *SlidersHandler) Reorder(c *gin.Context) {
	var req dto.ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "ids", Tag: "required", Message: "Daftar ID wajib dikirim."},
		})
		return
	}

	if err := h.svc.Reorder(req.IDs); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDERS_REORDERED", "Urutan slider berhasil diperbarui.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "slider.update",
		EntityType:  "slider",
		Description: "Mengubah urutan (reorder) slider",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/sliders/:id/restore — admin
func (h *SlidersHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID slider tidak valid.", nil)
		return
	}
	if err := h.svc.Restore(uint(id)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_RESTORED", "Slider berhasil dipulihkan.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(id)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "slider.restore",
		EntityType:  "slider",
		EntityID:    &eID,
		Description: "Memulihkan slider",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/sliders/bulk-delete — admin (bulk soft delete)
func (h *SlidersHandler) BulkDelete(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_BULK_DELETED", "Slider berhasil dihapus massal.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "slider.bulk_delete",
		EntityType:  "slider",
		Description: "Menghapus slider secara massal",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/sliders/bulk-restore — admin (bulk restore)
func (h *SlidersHandler) BulkRestore(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_BULK_RESTORED", "Slider berhasil dipulihkan massal.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "slider.bulk_restore",
		EntityType:  "slider",
		Description: "Memulihkan slider secara massal",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}
