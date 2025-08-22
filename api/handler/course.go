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

	course, httpcode, err := h.courseService.CreateCourse(teacherID, input)
	if err != nil {
		c.JSON(httpcode, gin.H{"error": err.Error()})
	}

	c.JSON(httpcode, course)
}

func (h *CourseHandler) GetCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := strings.ToUpper(c.DefaultQuery("sort_order", "ASC"))

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
	}

	response, err := h.courseService.GetCourses(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
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
	}

	response, err := h.courseService.GetTeacherCourses(teacherID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
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

type SearchByName struct {
	Name string `json:"name" binding:"required"`
}

func (h *CourseHandler) GetCoursesByTeacherName(c *gin.Context) {
	var teacherName SearchByName

	if err := c.ShouldBindJSON(&teacherName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := strings.ToUpper(c.DefaultQuery("sort_order", "ASC"))

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
	}

	response, err := h.courseService.GetCoursesByTeacherName(teacherName.Name, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
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

func (h *CourseHandler) GetCoursesByCourseName(c *gin.Context) {
	var courseName SearchByName
	if err := c.ShouldBindJSON(&courseName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := strings.ToUpper(c.DefaultQuery("sort_order", "ASC"))

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
	}

	response, err := h.courseService.GetCoursesByCourseName(courseName.Name, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
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

	course, httpcode, err := h.courseService.UpdateCourse(teacherID, int64(courseID), input)
	if err != nil {
		c.JSON(httpcode, gin.H{"error": err.Error()})
	}

	c.JSON(httpcode, SafeResponseCourses([]model.Course{*course}))
}

func SafeResponseCourses(courses []model.Course) []map[string]interface{} {
	// 只返回安全字段
	var safeCourses []map[string]interface{}
	for _, s := range courses {
		safeCourses = append(safeCourses, map[string]interface{}{
			"name":           s.Name,
			"remark":         s.Remark,
			"student_maxnum": s.StudentMaxNum,
			"hours":          s.Hours,
			"start_date":     s.StartDate,
			"semester":       s.Semester,
			"teacher_id":     s.TeacherID,
		})
	}
	return safeCourses
}
