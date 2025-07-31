package routes

import (
	"fmt"
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

	// 解析开始日期
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

	// 检查课程开始日期是否早于今天
	if course.StartDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "课程开始日期不能早于今天"})
		return
	}

	if err := h.db.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (h *CourseHandler) GetCourses(c *gin.Context) {
	// 获取所有课程
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

	// 检查课程开始时间是否早于当前时间
	if course.StartDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "课程已开始，不能删除"})
		return
	}

	// 删除课程
	if err := h.db.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程删除成功"})
}

func (h *CourseHandler) GetTeacherCourses(c *gin.Context) {
	teacherID := c.GetString("user_id")

	// 检查教师是否存在
	var courses []model.Course
	if err := h.db.Where("teacher_id = ?", teacherID).First(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "该教师不存在"})
		return
	}

	// 获取该教师的所有课程
	if err := h.db.Where("teacher_id = ?", teacherID).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, courses)
}

type UpdateCourseRequest struct {
	Name          *string `json:"name"`
	Remark        *string `json:"remark"`
	StudentMaxNum *int    `json:"student_maxnum"`
	Hours         *int    `json:"hours"`
	StartDate     *string `json:"start_date"`
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	// 获取课程ID
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的课程ID"})
		return
	}

	// 获取当前教师ID
	teacherID := c.GetString("user_id")

	// 先查询要更新的课程
	var course model.Course
	if err := h.db.Where("id = ? AND teacher_id = ?", courseID, teacherID).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在或权限不足"})
		return
	}

	// 解析请求
	var input UpdateCourseRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 准备更新字段
	updateData := make(map[string]interface{})
	if input.Name != nil {
		updateData["name"] = *input.Name
	}
	if input.Remark != nil {
		updateData["remark"] = *input.Remark
	}
	if input.StudentMaxNum != nil {
		// 检查新人数是否小于当前报名人数
		var count int64
		h.db.Model(&model.Enrollment{}).Where("course_id = ?", courseID).Count(&count)
		if *input.StudentMaxNum < int(count) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("新人数限制(%d)不能小于当前报名人数(%d)",
					*input.StudentMaxNum, count),
			})
			return
		}
		updateData["student_max_num"] = *input.StudentMaxNum
	}
	if input.Hours != nil {
		updateData["hours"] = *input.Hours
	}
	if input.StartDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *input.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "日期格式不正确，请使用YYYY-MM-DD格式",
			})
			return
		}
		updateData["start_date"] = parsedDate
	}

	// 执行更新 - 关键修改点
	if err := h.db.Model(&course).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回更新后的课程数据
	c.JSON(http.StatusOK, course)
}
