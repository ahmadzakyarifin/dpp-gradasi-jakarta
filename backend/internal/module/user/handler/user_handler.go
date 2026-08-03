package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	s   service.UserService
	cfg *config.Config
}

func NewUserHandler(s service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{s: s, cfg: cfg}
}

func (h *UserHandler) secureCookie() bool {
	return h.cfg != nil && h.cfg.App.Env == "production"
}

func (h *UserHandler) GetAll(c *gin.Context) {
	var req dto.UserQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	req.Normalize()

	users, total, err := h.s.GetPaginated(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	filters := map[string]any{
		"search": req.Search,
		"role":   req.Role,
		"status": req.Status,
		"sort":   req.Sort,
	}
	meta := helper.GetPaginationMeta(int(total), req.Page, req.Limit, filters)
	helper.SuccessResponse(c, http.StatusOK, "ADMIN_LIST", "Daftar admin berhasil diambil", gin.H{
		"items": mapper.UsersEntityToResponse(users),
		"meta":  meta,
	}, nil)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", validator.Errors(err))
		return
	}

	user, err := h.s.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "USER_DETAIL", "berhasil mengambil data user", gin.H{"user": mapper.UserEntityToResponse(*user)}, nil)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	created, err := h.s.Create(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "USER_CREATED", "user berhasil dibuat, email aktivasi dikirim", gin.H{"user": mapper.UserEntityToResponse(*created)}, nil)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	var req dto.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	updated, err := h.s.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "USER_UPDATED", "user berhasil diperbarui", gin.H{"user": mapper.UserEntityToResponse(*updated)}, nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	currentUserID, _ := c.Get("user_id")
	if uid, ok := currentUserID.(uint); ok && uint(id) == uid {
		helper.ErrorResponse(c, http.StatusBadRequest, "UNAUTHORIZED_ACTION", "Anda tidak diperbolehkan menghapus akun Anda sendiri", nil)
		return
	}

	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "USER_DELETED", "user berhasil dihapus", nil, nil)
}

func (h *UserHandler) ToggleStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := h.s.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	currentUserID, _ := c.Get("user_id")
	if uid, ok := currentUserID.(uint); ok && uint(id) == uid {
		helper.ErrorResponse(c, http.StatusBadRequest, "UNAUTHORIZED_ACTION", "Anda tidak diperbolehkan menonaktifkan akun Anda sendiri", nil)
		return
	}

	if user.RoleID == 1 && user.Status == "active" {
		admins, _ := h.s.GetAll(c.Request.Context(), 1)
		activeAdmins := 0
		for _, a := range admins {
			if a.Status == "active" {
				activeAdmins++
			}
		}
		if activeAdmins <= 1 {
			helper.ErrorResponse(c, http.StatusBadRequest, "UNAUTHORIZED_ACTION", "Gagal: Ini adalah satu-satunya akun Super Admin yang aktif di sistem", nil)
			return
		}
	}

	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "USER_STATUS_UPDATED", "status user berhasil diubah", nil, nil)
}

func (h *UserHandler) Activate(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	ctx := c.Request.Context()
	user, err := h.s.ActivateAccount(ctx, req.Token, req.Password)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	accessToken, _ := helper.GenerateAccessToken(user.ID, user.Email, user.RoleID, user.Name, user.RoleName, h.cfg.JWT.Secret, h.cfg.JWT.AccessTTLMinutes)
	refreshToken, expiry, _ := helper.GenerateRefreshToken(h.cfg.JWT.RefreshTTLHours)
	if err := h.s.SaveRefreshToken(ctx, user.ID, refreshToken, expiry); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "GAGAL_MENYIMPAN_SESSION", "gagal menyimpan session", nil)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", refreshToken, int(time.Until(expiry).Seconds()), "/", "", h.secureCookie(), true)

	helper.SuccessResponse(c, http.StatusOK, "AKUN_BERHASIL_DIAKTIFKAN", "akun berhasil diaktifkan", gin.H{
		"access_token": accessToken,
		"user": gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"role_id": user.RoleID,
		},
	}, nil)
}

func (h *UserHandler) ResendNotification(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.s.ResendNotification(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "LINK_AKTIVASI_DIKIRIM_ULANG", "link aktivasi berhasil dikirim ulang", nil, nil)
}

func (h *UserHandler) BulkResendNotification(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	result, err := h.s.BulkResendNotification(c.Request.Context(), req.IDs)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "BULK_RESEND_NOTIFICATION", fmt.Sprintf("%d link aktivasi diproses", result.Sent), gin.H{"result": result}, nil)
}

func (h *UserHandler) BulkDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", validator.Errors(err))
		return
	}

	currentUserID, _ := c.Get("user_id")
	for _, id := range req.IDs {
		if uid, ok := currentUserID.(uint); ok && id == uid {
			helper.ErrorResponse(c, http.StatusBadRequest, "UNAUTHORIZED_ACTION", "Operasi dibatalkan: Terdapat akun Anda sendiri dalam daftar hapus", nil)
			return
		}

		user, err := h.s.GetByID(c.Request.Context(), id)
		if err == nil && user != nil && user.RoleID == 1 {
			helper.ErrorResponse(c, http.StatusBadRequest, "UNAUTHORIZED_ACTION", fmt.Sprintf("Operasi dibatalkan: Akun %s adalah Admin dan tidak boleh dihapus", user.Name), nil)
			return
		}
	}

	if err := h.s.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "BULK_DELETE_USERS", fmt.Sprintf("%d pengguna berhasil dihapus", len(req.IDs)), nil, nil)
}

func (h *UserHandler) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "USER_RESTORE_SUCCESS", "akun berhasil dipulihkan", nil, nil)
}

func (h *UserHandler) BulkRestore(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", validator.Errors(err))
		return
	}

	if err := h.s.BulkRestore(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "BULK_RESTORE_USERS", fmt.Sprintf("%d pengguna berhasil dipulihkan", len(req.IDs)), nil, nil)
}

func (h *UserHandler) GetDependencyInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(
		c,
		http.StatusOK,
		"USER_DEPENDENCY_INFO_RETRIEVED",
		"User dependency information retrieved successfully.",
		info,
		nil,
	)
}

func (h *UserHandler) CheckUnique(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")
	excludeIDStr := c.DefaultQuery("exclude_id", "0")
	excludeID, _ := strconv.Atoi(excludeIDStr)

	isUnique, err := h.s.CheckUnique(c.Request.Context(), field, value, uint(excludeID))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "USER_UNIQUENESS_CHECKED", "berhasil", gin.H{
		"is_unique": isUnique,
	}, nil)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required,min=2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	value, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED_ACTION", "akses tidak sah", nil)
		return
	}

	userID, ok := value.(uint)
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED_ACTION", "user tidak valid", nil)
		return
	}

	updated, err := h.s.UpdateProfile(c.Request.Context(), userID, req.Name)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "USER_PROFILE_UPDATED", "profil berhasil diperbarui", gin.H{
		"user": mapper.UserEntityToResponse(*updated),
	}, nil)
}

// GetProfile mengembalikan profil user yang sedang login (dari token).
func (h *UserHandler) GetProfile(c *gin.Context) {
	value, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED_ACTION", "akses tidak sah", nil)
		return
	}

	userID, ok := value.(uint)
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED_ACTION", "user tidak valid", nil)
		return
	}

	user, err := h.s.GetProfile(c.Request.Context(), userID)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "PROFILE_RETRIEVED", "Profil berhasil diambil", gin.H{
		"user": mapper.UserEntityToResponse(*user),
	}, nil)
}

// ChangePassword mengganti password akun sendiri (old_password + new_password).
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	value, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED_ACTION", "akses tidak sah", nil)
		return
	}

	userID, ok := value.(uint)
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED_ACTION", "user tidak valid", nil)
		return
	}

	if err := h.s.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "PASSWORD_CHANGED", "Password berhasil diubah", nil, nil)
}

// VerifyEmail memverifikasi token aktivasi email untuk user yang sedang login.
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	value, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED_ACTION", "akses tidak sah", nil)
		return
	}

	userID, ok := value.(uint)
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED_ACTION", "user tidak valid", nil)
		return
	}

	if err := h.s.VerifyEmail(c.Request.Context(), userID, req.Token); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "EMAIL_VERIFIED", "Email berhasil diverifikasi dan diperbarui", nil, nil)
}
