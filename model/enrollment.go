package model

type Enrollment struct {
	CourseID  string `gorm:"primarykey;type:varchar(20)"`
	StudentID string `gorm:"primarykey;type:varchar(20)"`
}
