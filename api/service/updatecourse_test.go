package service

import (
	"errors"
	"testing"
	"time"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/repository/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCourse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockCourseRepository(ctrl)
	service := &CourseService{
		courseRepo: mockRepo,
	}

	teacherID := "teacher-123"
	courseID := int64(1)
	now := time.Now()
	validCourse := &model.Course{
		ID:            courseID,
		TeacherID:     teacherID,
		Name:          "Original Name",
		StudentMaxNum: 10,
		StartDate:     now,
	}

	tests := []struct {
		name        string
		input       UpdateCourseInput
		mockSetup   func()
		expected    *model.Course
		expectedErr error
	}{
		{
			name: "success - update name",
			input: UpdateCourseInput{
				Name: stringPtr("New Name"),
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetByID(courseID).Return(validCourse, nil)
				mockRepo.EXPECT().Update(validCourse, map[string]interface{}{
					"name": "New Name",
				}).Return(nil)
				mockRepo.EXPECT().GetByID(courseID).Return(&model.Course{
					ID:            courseID,
					TeacherID:     teacherID,
					Name:          "New Name",
					StudentMaxNum: 10,
					StartDate:     now,
				}, nil)
			},
			expected: &model.Course{
				ID:            courseID,
				TeacherID:     teacherID,
				Name:          "New Name",
				StudentMaxNum: 10,
				StartDate:     now,
			},
			expectedErr: nil,
		},
		{
			name: "success - update start date",
			input: UpdateCourseInput{
				StartDate: &now,
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetByID(courseID).Return(validCourse, nil)
				mockRepo.EXPECT().Update(validCourse, map[string]interface{}{
					"start_date": now.Unix(),
				}).Return(nil)
				mockRepo.EXPECT().GetByID(courseID).Return(&model.Course{
					ID:            courseID,
					TeacherID:     teacherID,
					Name:          "Original Name",
					StudentMaxNum: 10,
					StartDate:     now,
				}, nil)
			},
			expected: &model.Course{
				ID:            courseID,
				TeacherID:     teacherID,
				Name:          "Original Name",
				StudentMaxNum: 10,
				StartDate:     now,
			},
			expectedErr: nil,
		},
		{
			name: "error - unauthorized (empty teacherID)",
			input: UpdateCourseInput{
				Name: stringPtr("New Name"),
			},
			mockSetup:   func() {},
			expected:    nil,
			expectedErr: ErrUnauthorized,
		},
		{
			name: "error - course not found",
			input: UpdateCourseInput{
				Name: stringPtr("New Name"),
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetByID(courseID).Return(nil, errors.New("not found"))
			},
			expected:    nil,
			expectedErr: ErrCourseNotFound,
		},
		{
			name: "error - course not owned by teacher",
			input: UpdateCourseInput{
				Name: stringPtr("New Name"),
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetByID(courseID).Return(&model.Course{
					ID:        courseID,
					TeacherID: "other-teacher",
				}, nil)
			},
			expected:    nil,
			expectedErr: ErrCourseNotFound,
		},
		{
			name: "error - invalid student max num",
			input: UpdateCourseInput{
				StudentMaxNum: intPtr(5),
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetByID(courseID).Return(validCourse, nil)
				mockRepo.EXPECT().GetEnrollmentCount(courseID).Return(int64(8), nil)
			},
			expected:    nil,
			expectedErr: errors.New("invalid student number: 新人数限制(5)不能小于当前报名人数(8)"),
		},
		{
			name: "error - repository update failed",
			input: UpdateCourseInput{
				Name: stringPtr("New Name"),
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetByID(courseID).Return(validCourse, nil)
				mockRepo.EXPECT().Update(validCourse, gomock.Any()).Return(errors.New("update failed"))
			},
			expected:    nil,
			expectedErr: errors.New("update failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			// 根据测试用例决定使用哪个teacherID
			testTeacherID := teacherID
			if tt.name == "error - unauthorized (empty teacherID)" {
				testTeacherID = ""
			}

			result, err := service.UpdateCourse(testTeacherID, courseID, tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions to create pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
