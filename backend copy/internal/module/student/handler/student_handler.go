package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	s service.StudentService
}

func NewStudentHandler(s service.StudentService) *StudentHandler {
	return &StudentHandler{s: s}
}

func (h *StudentHandler) GetAll(c *gin.Context) {
	var req dto.StudentQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	list, total, err := h.s.GetPaginated(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	filters := map[string]any{
		"search":     req.Search,
		"filter":     req.Filter,
		"status":     req.Status,
		"entry_year": req.EntryYear,
		"class_id":   req.ClassID,
		"major_id":   req.MajorID,
		"sort":       req.Sort,
	}
	meta := helper.GetPaginationMeta(int(total), req.Page, req.Limit, filters)
	helper.SuccessResponse(c, http.StatusOK, "STUDENT_LIST", "berhasil", list, meta)
}

func (h *StudentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	student, err := h.s.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_DETAIL", "berhasil mengambil data siswa", gin.H{"data": student}, nil)
}

func (h *StudentHandler) Create(c *gin.Context) {
	var req dto.StudentCreateReq
	if err := c.ShouldBind(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	// Handle file upload
	file, _ := c.FormFile("image_path")
	req.ImageFile = file

	result, err := h.s.Create(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "STUDENT_CREATED", "data siswa berhasil dibuat", gin.H{"data": result}, nil)
}

func (h *StudentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	var req dto.StudentUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	file, _ := c.FormFile("image_path")
	req.ImageFile = file

	result, err := h.s.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_UPDATED", "data siswa berhasil diperbarui", gin.H{"data": result}, nil)
}

func (h *StudentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_DELETED", "data siswa berhasil dihapus", nil, nil)
}

func (h *StudentHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_RESTORED", "data siswa berhasil dipulihkan", nil, nil)
}

func (h *StudentHandler) ToggleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_STATUS_TOGGLED", "status siswa berhasil diubah", nil, nil)
}

func (h *StudentHandler) BulkDelete(c *gin.Context) {
	var req dto.StudentBulkDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "data tidak valid", validator.Errors(err))
		return
	}

	if err := h.s.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_BULK_DELETED", "data siswa terpilih berhasil dinonaktifkan", nil, nil)
}

func (h *StudentHandler) BulkRestore(c *gin.Context) {
	var req dto.StudentBulkRestoreReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "data tidak valid", validator.Errors(err))
		return
	}

	if err := h.s.BulkRestore(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_BULK_RESTORED", fmt.Sprintf("%d data siswa berhasil dipulihkan", len(req.IDs)), nil, nil)
}

func (h *StudentHandler) BulkGraduate(c *gin.Context) {
	var req dto.StudentBulkGraduateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "data tidak valid", validator.Errors(err))
		return
	}

	if err := h.s.BulkGraduate(c.Request.Context(), req.ClassID, req.StudentIDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_BULK_GRADUATED", "Kelulusan masal berhasil diproses", nil, nil)
}

func (h *StudentHandler) BulkPromote(c *gin.Context) {
	var req dto.StudentBulkPromoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "data tidak valid", validator.Errors(err))
		return
	}

	if err := h.s.BulkPromote(c.Request.Context(), req.SourceClassID, req.TargetClassID, req.StudentIDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_BULK_PROMOTED", "Perpindahan kelas masal berhasil diproses", nil, nil)
}

func (h *StudentHandler) GetClassHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	list, err := h.s.GetClassHistory(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_CLASS_HISTORY", "berhasil", gin.H{"data": list}, nil)
}

func (h *StudentHandler) GetFilters(c *gin.Context) {
	data, err := h.s.GetAcademicFilters(c.Request.Context())
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "STUDENT_ACADEMIC_FILTERS", "berhasil", gin.H{"data": data}, nil)
}

func (h *StudentHandler) GetDependencyInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_DEPENDENCY_INFO", "berhasil", gin.H{"data": info}, nil)
}

func (h *StudentHandler) CheckUnique(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")
	excludeID, _ := strconv.Atoi(c.DefaultQuery("exclude_id", "0"))

	isUnique, err := h.s.CheckUnique(c.Request.Context(), field, value, uint(excludeID))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "STUDENT_UNIQUE_CHECK", "berhasil", gin.H{"is_unique": isUnique}, nil)
}

func (h *StudentHandler) Export(c *gin.Context) {
	var req dto.StudentQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	data, err := h.s.ExportExcel(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=Daftar_Siswa.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
