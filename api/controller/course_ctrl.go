package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/service"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type CourseController struct {
	courseService *service.CourseService
}

func NewCourseController(courseService *service.CourseService) *CourseController {
	return &CourseController{courseService: courseService}
}

type CreateCourseRequest struct {
	Name          string `json:"name" binding:"required"`
	TeacherName   string `json:"teacher_name" binding:"required"`
	Remarks       string `json:"remarks"`
	StudentMaxNum int    `json:"student_max_num" binding:"required"`
	TimeMax       int    `json:"time_max" binding:"required"`
	TimeMin       int    `json:"time_min" binding:"required"`
}

func (c *CourseController) CreateCourse(ctx *gin.Context) {
	var req CreateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherID := ctx.MustGet("userID").(string)

	course := &model.Course{
		Name:          req.Name,
		TeacherID:     teacherID,
		Remark:        req.Remarks,
		StudentMaxNum: req.StudentMaxNum,
	}

	if err := c.courseService.CreateCourse(course); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, course)
}

func (c *CourseController) GetAllCourses(ctx *gin.Context) {
	courses, err := c.courseService.GetAllCourses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, courses)
}

func (c *CourseController) EnrollCourse(ctx *gin.Context) {
	courseID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	studentID := ctx.MustGet("userID").(int64)

	if err := c.courseService.Enroll(studentID, int64(courseID)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "enrolled successfully"})
}

func (c *CourseController) DeleteCourse(ctx *gin.Context) {
	courseID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	teacherID := ctx.MustGet("userID").(int64)

	if err := c.courseService.DeleteCourse(int64(courseID), teacherID); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "course deleted successfully"})
}

func (c *CourseController) GetMyCourses(ctx *gin.Context) {
	teacherID := ctx.MustGet("userID").(int64)

	courses, err := c.courseService.GetCoursesByTeacher(teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, courses)
}
