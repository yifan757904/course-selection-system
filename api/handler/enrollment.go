package handler

import (
	"net/http"
	"strconv"

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

	// 验证排序字段
	if !model.AllowedSortFields[sortBy] {
		sortBy = "id"
	}

	// 验证排序方向
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "ASC"
	}

	response, err := h.service.GetStudentCourses(studentIDCard, page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 只返回安全字段
	var safeCourses []map[string]interface{}
	for _, course := range response.Data {
		safeCourses = append(safeCourses, map[string]interface{}{
			"name":           course["name"],
			"remark":         course["remark"],
			"student_maxnum": course["student_maxnum"],
			"hours":          course["hours"],
			"start_date":     course["start_date"],
			"semester":       course["semester"],
			"teacher_id":     course["teacher_id"],
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data":        safeCourses,
		"total":       response.Total,
		"page":        response.Page,
		"page_size":   response.PageSize,
		"total_pages": response.TotalPages,
	})
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

func (h *EnrollmentHandler) GetStudentsByCourseID(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效课程ID"})
		return
	}
	students, err := h.service.GetStudentsByCourseID(int64(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 只返回安全字段
	var safeStudents []map[string]interface{}
	for _, s := range students {
		safeStudents = append(safeStudents, map[string]interface{}{
			"name": s.Name,
		})
	}
	c.JSON(http.StatusOK, gin.H{"students": safeStudents})
}
