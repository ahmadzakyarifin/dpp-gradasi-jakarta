package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc    service.UserService
	logSvc activitylogservice.ActivityLogService
	cfg    *config.Config
}

func NewUserHandler(svc service.UserService, logSvc activitylogservice.ActivityLogService, cfg *config.Config) *UserHandler {
	return &UserHandler{svc: svc, logSvc: logSvc, cfg: cfg}
}

func (h *UserHandler) getAuthUserID(c *gin.Context) (uint, error) {
	val, exists := c.Get("user_id")
	if !exists {
		return 0, helper.NewServiceError("UNAUTHORIZED", "User tidak terautentikasi", nil)
	}
	// Tergantung middleware, bisa float64 atau uint
	if id, ok := val.(float64); ok {
		return uint(id), nil
	}
	if id, ok := val.(uint); ok {
		return id, nil
	}
	return 0, helper.NewServiceError("UNAUTHORIZED", "ID User tidak valid", nil)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := h.getAuthUserID(c)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	resp, err := h.svc.GetProfile(userID)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusNotFound, svcErr.Code, svcErr.Message, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "PROFILE_RETRIEVED", "Profil berhasil diambil", resp, nil)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := h.getAuthUserID(c)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	var req dto.ProfileUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{{Field: "input", Tag: "invalid", Message: err.Error()}})
		return
	}

	file, _ := c.FormFile("photo")
	req.Photo = file

	resp, msg, err := h.svc.UpdateProfile(userID, &req)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusBadRequest, svcErr.Code, svcErr.Message, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "PROFILE_UPDATED", msg, resp, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "user.update_profile",
		EntityType:  "user",
		EntityID:    &userID,
		EntityLabel: resp.Name,
		Description: "Mengupdate profil sendiri",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	userID, err := h.getAuthUserID(c)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{{Field: "token", Tag: "required", Message: "Token wajib diisi"}})
		return
	}

	if err := h.svc.VerifyEmail(userID, &req); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusBadRequest, svcErr.Code, svcErr.Message, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "EMAIL_VERIFIED", "Email berhasil diverifikasi dan diperbarui", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "user.verify_email",
		EntityType:  "user",
		EntityID:    &userID,
		Description: "Memverifikasi email baru",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, err := h.getAuthUserID(c)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{{Field: "password", Tag: "invalid", Message: "Format password tidak valid"}})
		return
	}

	if err := h.svc.ChangePassword(userID, &req); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusBadRequest, svcErr.Code, svcErr.Message, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "PASSWORD_CHANGED", "Password berhasil diubah", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "user.change_password",
		EntityType:  "user",
		EntityID:    &userID,
		Description: "Mengubah password sendiri",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

func (h *UserHandler) GetAdmins(c *gin.Context) {
	resp, err := h.svc.GetAdmins()
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "ADMIN_LIST", "Daftar admin berhasil diambil", resp, nil)
}

func (h *UserHandler) CreateAdmin(c *gin.Context) {
	var req dto.AdminCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{{Field: "input", Tag: "invalid", Message: err.Error()}})
		return
	}

	resp, err := h.svc.CreateAdmin(&req, h.cfg.App.URL)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusBadRequest, svcErr.Code, svcErr.Message, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "ADMIN_CREATED", "Undangan aktivasi admin berhasil dikirim ke email.", resp, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "user.create",
		EntityType:  "user",
		EntityID:    &resp.ID,
		EntityLabel: resp.Name,
		Description: "Mengundang admin baru (Email Invitation)",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/admin/users/:id/resend-activation
func (h *UserHandler) ResendActivation(c *gin.Context) {
	adminID, err := h.getAuthUserID(c)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	targetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID tidak valid", nil)
		return
	}

	if err := h.svc.ResendActivation(adminID, uint(targetID), h.cfg.App.URL); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		code := http.StatusBadRequest
		if svcErr.Code == "NOT_FOUND" {
			code = http.StatusNotFound
		}
		helper.ErrorResponse(c, code, svcErr.Code, svcErr.Message, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ACTIVATION_RESENT", "Undangan aktivasi berhasil dikirim ulang ke email admin.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(targetID)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "user.resend_activation",
		EntityType:  "user",
		EntityID:    &eID,
		Description: "Mengirim ulang undangan aktivasi admin",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// PUT /api/v1/admin/users/:id/status
func (h *UserHandler) SetAdminStatus(c *gin.Context) {
	adminID, err := h.getAuthUserID(c)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	targetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID tidak valid", nil)
		return
	}

	var req dto.AdminStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{{Field: "status", Tag: "invalid", Message: "Status harus 'active' atau 'inactive'"}})
		return
	}

	if err := h.svc.SetAdminStatus(adminID, uint(targetID), &req); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		code := http.StatusBadRequest
		if svcErr.Code == "NOT_FOUND" {
			code = http.StatusNotFound
		}
		helper.ErrorResponse(c, code, svcErr.Code, svcErr.Message, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ADMIN_STATUS_UPDATED", "Status admin berhasil diperbarui", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(targetID)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "user.update_status",
		EntityType:  "user",
		EntityID:    &eID,
		Description: "Mengubah status admin menjadi " + req.Status,
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

func (h *UserHandler) DeleteAdmin(c *gin.Context) {
	adminID, err := h.getAuthUserID(c)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	targetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID tidak valid", nil)
		return
	}

	if err := h.svc.DeleteAdmin(adminID, uint(targetID)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		code := http.StatusBadRequest
		if svcErr.Code == "NOT_FOUND" {
			code = http.StatusNotFound
		}
		helper.ErrorResponse(c, code, svcErr.Code, svcErr.Message, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ADMIN_DELETED", "Admin berhasil dihapus", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(targetID)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "user.delete",
		EntityType:  "user",
		EntityID:    &eID,
		Description: "Menghapus admin",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/admin/users/:id/restore
func (h *UserHandler) RestoreAdmin(c *gin.Context) {
	adminID, err := h.getAuthUserID(c)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	targetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID tidak valid", nil)
		return
	}

	if err := h.svc.RestoreAdmin(adminID, uint(targetID)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "ADMIN_RESTORED", "Admin berhasil dipulihkan", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	eID := uint(targetID)
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "user.restore",
		EntityType:  "user",
		EntityID:    &eID,
		Description: "Memulihkan admin",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/admin/users/bulk-delete
func (h *UserHandler) BulkDeleteAdmin(c *gin.Context) {
	adminID, err := h.getAuthUserID(c)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "ids", Tag: "required", Message: "IDs wajib diisi minimal 1."},
		})
		return
	}

	if err := h.svc.BulkDeleteAdmin(adminID, req.IDs); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "ADMIN_BULK_DELETED", "Admin berhasil dihapus massal", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "user.bulk_delete",
		EntityType:  "user",
		Description: "Menghapus admin secara massal",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// POST /api/v1/admin/users/bulk-restore
func (h *UserHandler) BulkRestoreAdmin(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "ids", Tag: "required", Message: "IDs wajib diisi minimal 1."},
		})
		return
	}

	if err := h.svc.BulkRestoreAdmin(req.IDs); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "ADMIN_BULK_RESTORED", "Admin berhasil dipulihkan massal", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "user.bulk_restore",
		EntityType:  "user",
		Description: "Memulihkan admin secara massal",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}
