package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/service"
)

type CourseHandler struct {
	courseService *service.CourseService
}

func NewCourseHandler(courseService *service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	teacherID := c.GetString("user_id")

	var input service.CreateCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseService.CreateCourse(teacherID, input)
	if err != nil {
		switch err {
		case service.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case service.ErrTeacherNotFound, service.ErrInvalidDateFormat, service.ErrPastStartDate:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (h *CourseHandler) GetCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := strings.ToUpper(c.DefaultQuery("sort_order", "ASC"))

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

	input := service.GetCoursesInput{
		Pagination: model.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Fields:    fields,
	}

	response, err := h.courseService.GetCourses(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	teacherID := c.GetString("user_id")
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效课程ID"})
		return
	}

	err = h.courseService.DeleteCourse(teacherID, int64(courseID))
	if err != nil {
		switch err {
		case service.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case service.ErrCourseNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case service.ErrCourseHasStudents, service.ErrCourseStarted:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程删除成功"})
}

func (h *CourseHandler) GetTeacherCourses(c *gin.Context) {
	teacherID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := strings.ToUpper(c.DefaultQuery("sort_order", "ASC"))

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

	input := service.GetCoursesInput{
		Pagination: model.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Fields:    fields,
	}

	response, err := h.courseService.GetTeacherCourses(teacherID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CourseHandler) GetCoursesByTeacherName(c *gin.Context) {
	teacherName := c.Param("teachername")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := strings.ToUpper(c.DefaultQuery("sort_order", "ASC"))
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

	input := service.GetCoursesInput{
		Pagination: model.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Fields:    fields,
	}

	response, err := h.courseService.GetCoursesByTeacherName(teacherName, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CourseHandler) GetCoursesByCourseName(c *gin.Context) {
	courseName := c.Param("coursename")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := strings.ToUpper(c.DefaultQuery("sort_order", "ASC"))
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

	input := service.GetCoursesInput{
		Pagination: model.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Fields:    fields,
	}

	response, err := h.courseService.GetCoursesByCourseName(courseName, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	teacherID := c.GetString("user_id")
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的课程ID"})
		return
	}

	var input service.UpdateCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseService.UpdateCourse(teacherID, int64(courseID), input)
	if err != nil {
		switch err {
		case service.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case service.ErrCourseNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case service.ErrInvalidDateFormat, service.ErrInvalidStudentNum:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, course)
}
