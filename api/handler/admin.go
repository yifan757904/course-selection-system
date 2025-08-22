package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liuyifan1996/course-selection-system/api/service"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

type CreateAdminRequest struct {
	JobNo    string `json:"job_no" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateAdminRequest struct {
	ID       int64   `json:"id" binding:"required"`
	JobNo    *string `json:"job_no"`
	Password *string `json:"password"`
}

func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input := service.CreateAdminInput{
		JobNo:    req.JobNo,
		Password: req.Password,
	}
	admin, err := h.adminService.CreateAdmin(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": admin.ID, "job_no": admin.JobNo})
}

type DeleteAdminRequest struct {
	JobNo    string `json:"job_no" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	var req DeleteAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input := service.DeleteAdminInput{
		JobNo:    req.JobNo,
		Password: req.Password,
	}
	if err := h.adminService.DeleteAdmin(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input := service.UpdateAdminInput{
		ID:       req.ID,
		JobNo:    req.JobNo,
		Password: req.Password,
	}
	admin, err := h.adminService.UpdateAdmin(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": admin.ID, "job_no": admin.JobNo})
}

func (h *AdminHandler) GetAdminByID(c *gin.Context) {
	id := c.Param("id")
	var adminID int64
	_, err := fmt.Sscan(id, &adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}
	admin, err := h.adminService.GetAdminByID(adminID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": admin.ID, "job_no": admin.JobNo})
}

func (h *AdminHandler) GetAdminByJobNo(c *gin.Context) {
	jobNo := c.Param("job_no")
	admin, err := h.adminService.GetAdminByJobNo(jobNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": admin.ID, "job_no": admin.JobNo})
}

func (h *AdminHandler) ListAdmins(c *gin.Context) {
	admins, err := h.adminService.ListAdmins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"admins": admins})
}

func (h *AdminHandler) AdminLogin(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	admin, err := h.adminService.GetAdminByJobNo(req.JobNo)
	if err != nil || admin.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "登录失败，用户名或密码错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "admin_id": admin.ID})
}
