package service

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/repository"
)

// http.StatusBadRequest:ErrTeacherNotFound,ErrInvalidDateFormat,ErrPastStartDate
// http.StatusUnauthorized:ErrUnauthorized
var (
	ErrUnauthorized      = errors.New("未认证用户")
	ErrTeacherNotFound   = errors.New("教师不存在或权限不足")
	ErrInvalidDateFormat = errors.New("日期格式不正确，请使用YYYY-MM-DD格式")
	ErrPastStartDate     = errors.New("课程开始日期不能早于今天")
	ErrCourseNotFound    = errors.New("课程不存在或权限不足")
	ErrCourseHasStudents = errors.New("课程已有学生选课，不能删除")
	ErrCourseStarted     = errors.New("课程已开始，不能删除")
	ErrInvalidStudentNum = errors.New("新人数限制不能小于当前报名人数")
	ErrInvalidCourseID   = errors.New("无效的课程ID")
)

type CourseService struct {
	courseRepo repository.CourseRepository
	userRepo   repository.AuthRepository
}

func NewCourseService(courseRepo repository.CourseRepository, userRepo repository.AuthRepository) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
		userRepo:   userRepo,
	}
}

type CreateCourseInput struct {
	Name          string `json:"name"`
	Remark        string `json:"remark"`
	StudentMaxNum int    `json:"student_maxnum"`
	Hours         int    `json:"hours"`
	StartDate     string `json:"start_date"`
	Semester      string `json:"semester"`
}

func (s *CourseService) CreateCourse(teacherID string, input CreateCourseInput) (*model.Course, int, error) {
	if teacherID == "" {
		return nil, http.StatusUnauthorized, ErrUnauthorized
	}

	// 验证教师是否存在
	teacher, err := s.userRepo.FindByIDCard(teacherID)
	if err != nil || teacher == nil || teacher.Role != "teacher" {
		return nil, http.StatusBadRequest, ErrTeacherNotFound
	}

	// 解析和验证开始日期
	var startDate time.Time
	switch v := any(input.StartDate).(type) {
	case string:
		startDate, err = time.Parse("2006-01-02", v)
		if err != nil {
			return nil, http.StatusBadRequest, ErrInvalidDateFormat
		}
	case time.Time:
		startDate = v
	default:
		return nil, http.StatusBadRequest, ErrInvalidDateFormat
	}

	semester := model.GetSemesterByDate(startDate, model.DefaultSemesterConfig)
	if input.Semester != "" {
		if err := model.ValidateSemester(input.Semester); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("invalid semester: %v", err)
		}
		semester = input.Semester
	}

	course := &model.Course{
		Name:          input.Name,
		TeacherID:     teacher.IDCard,
		Remark:        input.Remark,
		StudentMaxNum: input.StudentMaxNum,
		Hours:         input.Hours,
		StartDate:     startDate,
		Semester:      semester,
	}

	if err := s.courseRepo.Create(course); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return course, http.StatusCreated, nil
}

type GetCoursesInput struct {
	Pagination model.Pagination
	SortBy     string
	SortOrder  string
}

func (s *CourseService) GetCourses(input GetCoursesInput) (*model.PaginatedResponse[map[string]interface{}], error) {
	semester := model.GetSemesterByDate(time.Now(), model.DefaultSemesterConfig)
	courses, total, err := s.courseRepo.GetAll(input.Pagination, input.SortBy, input.SortOrder, semester)
	if err != nil {
		return nil, err
	}

	return &model.PaginatedResponse[map[string]interface{}]{
		Data:       courses,
		Total:      total,
		Page:       input.Pagination.Page,
		PageSize:   input.Pagination.PageSize,
		TotalPages: int((total + int64(input.Pagination.PageSize) - 1) / int64(input.Pagination.PageSize)),
	}, nil
}

func (s *CourseService) DeleteCourse(teacherID string, courseID int64) error {
	if teacherID == "" {
		return ErrUnauthorized
	}

	// 检查课程是否属于该教师
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil || course.TeacherID != teacherID {
		return ErrCourseNotFound
	}

	// 检查是否有学生选课
	count, err := s.courseRepo.GetEnrollmentCount(courseID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrCourseHasStudents
	}

	// 检查课程开始时间
	if course.StartDate.Before(time.Now()) {
		return ErrCourseStarted
	}

	return s.courseRepo.Delete(courseID)
}

func (s *CourseService) GetTeacherCourses(teacherID string, input GetCoursesInput) (*model.PaginatedResponse[map[string]interface{}], error) {
	semester := model.GetSemesterByDate(time.Now(), model.DefaultSemesterConfig)
	courses, total, err := s.courseRepo.GetByTeacherID(teacherID, input.Pagination, input.SortBy, input.SortOrder, semester)
	if err != nil {
		return nil, err
	}

	return &model.PaginatedResponse[map[string]interface{}]{
		Data:       courses,
		Total:      total,
		Page:       input.Pagination.Page,
		PageSize:   input.Pagination.PageSize,
		TotalPages: int((total + int64(input.Pagination.PageSize) - 1) / int64(input.Pagination.PageSize)),
	}, nil
}

func (s *CourseService) GetCoursesByTeacherName(teacherName string, input GetCoursesInput) (*model.PaginatedResponse[map[string]interface{}], error) {
	semester := model.GetSemesterByDate(time.Now(), model.DefaultSemesterConfig)
	courses, total, err := s.courseRepo.GetByTeacherName(teacherName, input.Pagination, input.SortBy, input.SortOrder, semester)
	if err != nil {
		return nil, err
	}

	return &model.PaginatedResponse[map[string]interface{}]{
		Data:       courses,
		Total:      total,
		Page:       input.Pagination.Page,
		PageSize:   input.Pagination.PageSize,
		TotalPages: int((total + int64(input.Pagination.PageSize) - 1) / int64(input.Pagination.PageSize)),
	}, nil
}

func (s *CourseService) GetCoursesByCourseName(courseName string, input GetCoursesInput) (*model.PaginatedResponse[map[string]interface{}], error) {
	semester := model.GetSemesterByDate(time.Now(), model.DefaultSemesterConfig)
	courses, total, err := s.courseRepo.GetByCourseName(courseName, input.Pagination, input.SortBy, input.SortOrder, semester)
	if err != nil {
		return nil, err
	}

	return &model.PaginatedResponse[map[string]interface{}]{
		Data:       courses,
		Total:      total,
		Page:       input.Pagination.Page,
		PageSize:   input.Pagination.PageSize,
		TotalPages: int((total + int64(input.Pagination.PageSize) - 1) / int64(input.Pagination.PageSize)),
	}, nil
}

type UpdateCourseInput struct {
	Name          *string `json:"name"`
	Remark        *string `json:"remark"`
	StudentMaxNum *int    `json:"student_maxnum"`
	Hours         *int    `json:"hours"`
	StartDate     *string `json:"start_date"`
	Semester      *string `json:"semester"`
}

func (s *CourseService) UpdateCourse(teacherID string, courseID int64, input UpdateCourseInput) (*model.Course, int, error) {
	if teacherID == "" {
		return nil, http.StatusUnauthorized, ErrUnauthorized
	}

	// 检查课程是否存在且属于该教师
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil || course.TeacherID != teacherID {
		return nil, http.StatusBadRequest, ErrCourseNotFound
	}

	updateData := make(map[string]interface{})
	if input.Name != nil {
		updateData["name"] = *input.Name
	}
	if input.Remark != nil {
		updateData["remark"] = *input.Remark
	}
	if input.StudentMaxNum != nil {
		count, err := s.courseRepo.GetEnrollmentCount(courseID)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		if *input.StudentMaxNum < int(count) {
			return nil, http.StatusBadRequest, fmt.Errorf("%w: 新人数限制(%d)不能小于当前报名人数(%d)",
				ErrInvalidStudentNum, *input.StudentMaxNum, count)
		}
		updateData["student_max_num"] = *input.StudentMaxNum
	}
	if input.Hours != nil {
		updateData["hours"] = *input.Hours
	}
	if input.StartDate != nil {
		var startDate time.Time
		switch v := any(*input.StartDate).(type) {
		case string:
			startDate, err = time.Parse("2006-01-02", v)
			if err != nil {
				return nil, http.StatusBadRequest, ErrInvalidDateFormat
			}
		case time.Time:
			startDate = v
		default:
			return nil, http.StatusBadRequest, ErrInvalidDateFormat
		}
		// 处理时区
		validatedTime, err := s.parseAndValidateTime(&startDate, os.Getenv("TIME_LOCAL"))
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		updateData["start_date"] = *validatedTime

		semester := model.GetSemesterByDate(startDate, model.DefaultSemesterConfig)
		if input.Semester != nil {
			if err := model.ValidateSemester(*input.Semester); err != nil {
				return nil, http.StatusBadRequest, fmt.Errorf("invalid semester: %v", err)
			}
			semester = *input.Semester
		}
		updateData["semester"] = semester
	}

	if err := s.courseRepo.Update(course, updateData); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	course, _ = s.courseRepo.GetByID(courseID)

	return course, http.StatusOK, nil
}

// 统一的时区处理函数
func (s *CourseService) parseAndValidateTime(inputTime *time.Time, timezone string) (*time.Time, error) {
	if inputTime == nil {
		return nil, nil
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	// 转换为目标时区时间
	localTime := inputTime.In(loc)

	// 验证时间是否在未来
	if localTime.Before(time.Now().In(loc)) {
		return nil, ErrPastStartDate
	}

	return &localTime, nil
}
