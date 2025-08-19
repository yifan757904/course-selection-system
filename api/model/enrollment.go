package model

type Enrollment struct {
	CourseID  int64 `gorm:"primarykey"`
	StudentID int64 `gorm:"primarykey"`
}
