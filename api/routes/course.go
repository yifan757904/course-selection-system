package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liuyifan1996/course-selection-system/api/model"
	"gorm.io/gorm"
)

type CourseHandler struct {
	db *gorm.DB
}

func NewCourseHandler(db *gorm.DB) *CourseHandler {
	return &CourseHandler{db: db}
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	var input struct {
		Name          string `json:"name" binding:"required"`
		Remark        string `json:"remark"`
		StudentMaxNum int    `json:"student_maxnum" binding:"required"`
		Hours         int    `json:"hours" binding:"required"`
		StartDate     string `json:"start_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取教师信息
	var teacher model.User
	if err := h.db.Where("id_card = ? AND rule = 'teacher'", teacherID).First(&teacher).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "教师不存在或权限不足"})
		return
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式不正确，请使用YYYY-MM-DD格式"})
		return
	}

	course := model.Course{
		Name:          input.Name,
		TeacherID:     teacher.IDCard,
		Remark:        input.Remark,
		StudentMaxNum: input.StudentMaxNum,
		Hours:         input.Hours,
		StartDate:     startDate,
	}

	if err := h.db.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (h *CourseHandler) GetCourses(c *gin.Context) {
	var courses []model.Course
	if err := h.db.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	teacherID := c.GetString("user_id")
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效课程ID"})
		return
	}

	// 检查课程是否属于该教师
	var course model.Course
	if err := h.db.Where("id = ? AND teacher_id = ?", courseID, teacherID).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在或权限不足"})
		return
	}

	// 检查是否有学生选课
	var count int64
	if h.db.Model(&model.Enrollment{}).Where("course_id = ?", courseID).Count(&count); count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "课程已有学生选课，不能删除"})
		return
	}

	if err := h.db.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程删除成功"})
}
