package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/service"
)

type EnrollmentHandler struct {
	service *service.EnrollmentService
}

func NewEnrollmentHandler(service *service.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{service: service}
}

func (h *EnrollmentHandler) Enroll(c *gin.Context) {
	studentIDCard := c.GetString("user_id")
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效课程ID"})
		return
	}

	if err := h.service.Enroll(studentIDCard, courseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "选课成功"})
}

func (h *EnrollmentHandler) GetStudentCourses(c *gin.Context) {
	studentIDCard := c.GetString("user_id")

	// 绑定分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "ASC")

	fieldsParam := c.DefaultQuery("fields", "")
	var fields []string
	if fieldsParam != "" {
		fields = strings.Split(fieldsParam, ",")
	}

	// 验证排序字段
	if !model.AllowedSortFields[sortBy] {
		sortBy = "id"
	}

	// 验证排序方向
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "ASC"
	}

	response, err := h.service.GetStudentCourses(studentIDCard, page, pageSize, sortBy, sortOrder, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *EnrollmentHandler) DeleteEnroll(c *gin.Context) {
	studentIDCard := c.GetString("user_id")
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效课程ID"})
		return
	}

	if err := h.service.DeleteEnrollment(studentIDCard, courseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程退选成功"})
}
