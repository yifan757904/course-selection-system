package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liuyifan1996/course-selection-system/api/model"
	"gorm.io/gorm"
)

type EnrollmentHandler struct {
	db *gorm.DB
}

func NewEnrollmentHandler(db *gorm.DB) *EnrollmentHandler {
	return &EnrollmentHandler{db: db}
}

func (h *EnrollmentHandler) Enroll(c *gin.Context) {
	// 检查学生是否存在
	var student model.User
	studentID_card := c.GetString("user_id")
	if err := h.db.Where("id_card = ? AND role = 'student'", studentID_card).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "学生不存在"})
		return
	}

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效课程ID"})
		return
	}

	studentID := student.ID
	// 检查是否已选课
	var existing model.Enrollment
	if h.db.Where("student_id = ? AND course_id = ?", studentID, courseID).First(&existing).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已选过该课程"})
		return
	}

	// 检查课程是否存在
	var course model.Course
	if err := h.db.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	// 检查课程是否已开始
	if course.StartDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "课程已开始，不能选课"})
		return
	}

	// 检查课程是否已满
	var count int64
	h.db.Model(&model.Enrollment{}).Where("course_id = ?", courseID).Count(&count)
	if count >= int64(course.StudentMaxNum) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "课程人数已满"})
		return
	}

	enrollment := model.Enrollment{
		StudentID: int64(studentID),
		CourseID:  int64(courseID),
	}

	if err := h.db.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "选课成功"})
}

func (h *EnrollmentHandler) GetStudentCourses(c *gin.Context) {
	var student model.User
	studentID_card := c.GetString("user_id")
	//检查学生是否存在
	if err := h.db.Where("id_card = ? AND role = 'student'", studentID_card).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "学生不存在"})
		return
	}
	studentID := student.ID

	//查找学生的选课记录
	var enrollments []model.Enrollment
	if err := h.db.Where("student_id = ?", studentID).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//保存所选课程id
	var courseIDs []int64
	for _, e := range enrollments {
		courseIDs = append(courseIDs, e.CourseID)
	}

	// 绑定分页参数
	var pagination model.Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	page_size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	pagination.Page = page
	pagination.PageSize = page_size
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := strings.ToUpper(c.DefaultQuery("sort_order", "ASC"))

	allowedSortFields := model.AllowedSortFields

	// 验证排序字段
	if !allowedSortFields[sortBy] {
		sortBy = "id"
	}

	// 验证排序方向
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "ASC"
	}

	var total int64

	// 获取总数
	h.db.Model(&enrollments).Where("student_id = ?", studentID).Count(&total)

	var courses []model.Course

	// 分页查询
	if err := h.db.Offset(pagination.Offset()).
		Limit(pagination.Limit()).
		Where("id IN ?", courseIDs).
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Find(&courses).Error; err != nil {
		c.JSON(500, gin.H{"error": "查询失败"})
		return
	}

	// 构建响应
	response := model.PaginatedResponse[model.Course]{
		Data:       courses,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)),
	}

	c.JSON(200, response)

}

func (h *EnrollmentHandler) DeleteEnroll(c *gin.Context) {
	// 检查学生是否存在
	var student model.User
	studentID_card := c.GetString("user_id")
	if err := h.db.Where("id_card = ? AND role = 'student'", studentID_card).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "学生不存在"})
		return
	}

	// 获取课程ID
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效课程ID"})
		return
	}
	studentID := student.ID

	// 检查是否已选课
	var existing model.Enrollment
	if h.db.Where("student_id = ? AND course_id = ?", studentID, courseID).First(&existing).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择该课程"})
		return
	}

	var course model.Course
	if err := h.db.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}
	// 检查课程是否已开始
	if course.StartDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "课程已开始，不能退选"})
		return
	}

	// 删除选课
	if err := h.db.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程退选成功"})
}
