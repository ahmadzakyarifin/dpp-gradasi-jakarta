package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/service"
	"github.com/gin-gonic/gin"
)

type PengurusHandler struct {
	svc    service.PengurusService
	logSvc activitylogservice.ActivityLogService
}

func NewPengurusHandler(svc service.PengurusService, logSvc activitylogservice.ActivityLogService) *PengurusHandler {
	return &PengurusHandler{svc: svc, logSvc: logSvc}
}

// GET /api/v1/pengurus
func (h *PengurusHandler) ListPublic(c *gin.Context) {
	var q dto.PengurusQuery
	c.ShouldBindQuery(&q)

	resp, err := h.svc.GetAllPublic(q)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_LIST", "Daftar pengurus berhasil diambil", resp, nil)
}

// GET /api/v1/pengurus/regions
func (h *PengurusHandler) Regions(c *gin.Context) {
	resp, err := h.svc.GetRegions()
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_REGIONS", "Daftar wilayah berhasil diambil", resp, nil)
}

// GET /api/v1/admin/pengurus
func (h *PengurusHandler) ListAdmin(c *gin.Context) {
	var q dto.PengurusQuery
	c.ShouldBindQuery(&q)

	resp, err := h.svc.GetAllAdmin(q)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_LIST", "Daftar pengurus berhasil diambil", resp, nil)
}

// POST /api/v1/admin/pengurus
func (h *PengurusHandler) Create(c *gin.Context) {
	var req dto.PengurusRequest
	if err := c.ShouldBind(&req); err != nil { // ShouldBind for form-data
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "form", Tag: "invalid", Message: "Data form tidak valid atau tidak lengkap."},
		})
		return
	}

	// Validate required file manually since gin binding file is tricky sometimes
	file, err := c.FormFile("image")
	if err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "image", Tag: "required", Message: "Image wajib diunggah."},
		})
		return
	}
	req.Image = file

	resp, err := h.svc.Create(&req)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "PENGURUS_CREATED", "Pengurus berhasil ditambahkan", resp, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "pengurus.create",
		EntityType:  "pengurus",
		EntityID:    &resp.ID,
		EntityLabel: resp.Name,
		Description: "Menambahkan pengurus baru",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// PUT /api/v1/admin/pengurus/:id
func (h *PengurusHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID pengurus tidak valid.", nil)
		return
	}

	var req dto.PengurusRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "form", Tag: "invalid", Message: "Data form tidak valid atau tidak lengkap."},
		})
		return
	}

	// Optional file upload
	file, _ := c.FormFile("image")
	if file != nil {
		req.Image = file
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
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_UPDATED", "Data pengurus berhasil diperbarui", resp, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "pengurus.update",
		EntityType:  "pengurus",
		EntityID:    &resp.ID,
		EntityLabel: resp.Name,
		Description: "Memperbarui data pengurus",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// DELETE /api/v1/admin/pengurus/:id
func (h *PengurusHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID pengurus tidak valid.", nil)
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
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_DELETED", "Pengurus berhasil dihapus", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(id)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "pengurus.delete",
		EntityType:  "pengurus",
		EntityID:    &eID,
		Description: "Menghapus pengurus",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/admin/pengurus/:id/restore
func (h *PengurusHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID pengurus tidak valid.", nil)
		return
	}
	if err := h.svc.Restore(uint(id)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_RESTORED", "Pengurus berhasil dipulihkan", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(id)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "pengurus.restore",
		EntityType:  "pengurus",
		EntityID:    &eID,
		Description: "Memulihkan pengurus",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/admin/pengurus/bulk-delete — admin (bulk soft delete)
func (h *PengurusHandler) BulkDelete(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_BULK_DELETED", "Pengurus berhasil dihapus massal", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "pengurus.bulk_delete",
		EntityType:  "pengurus",
		Description: "Menghapus pengurus secara massal",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/admin/pengurus/bulk-restore — admin (bulk restore)
func (h *PengurusHandler) BulkRestore(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_BULK_RESTORED", "Pengurus berhasil dipulihkan massal", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "pengurus.bulk_restore",
		EntityType:  "pengurus",
		Description: "Memulihkan pengurus secara massal",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}
