package routes

import (
	"net/http"
	"strconv"

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
	var student model.User
	studentID_card := c.GetString("user_id")
	if err := h.db.Where("id_card = ? AND rule = 'student'", studentID_card).First(&student).Error; err != nil {
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

	// 检查课程容量
	var course model.Course
	if err := h.db.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

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

func (h *EnrollmentHandler) GetMyCourses(c *gin.Context) {
	studentID := c.GetString("user_id")

	var enrollments []model.Enrollment
	if err := h.db.Where("student_id = ?", studentID).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var courseIDs []int64
	for _, e := range enrollments {
		courseIDs = append(courseIDs, e.CourseID)
	}

	var courses []model.Course
	if err := h.db.Where("id IN ?", courseIDs).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}
