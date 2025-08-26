package model

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	Name          string    `gorm:"size:60;not null"`
	TeacherID     string    `gorm:"size:20;not null"`
	Remark        string    `gorm:"size:200"`
	StudentMaxNum int       `gorm:"not null"`
	Hours         int       `gorm:"not null"`
	StartDate     time.Time `gorm:"type:datetime;not null"`
	Semester      string    `gorm:"size:20;index"`

	// 关联关系
	Students []User `gorm:"many2many:enrollments;foreignKey:ID;joinForeignKey:CourseID;References:ID;joinReferences:StudentID"`
}
