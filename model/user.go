package model

type User struct {
	ID   string `gorm:"primarykey;type:varchar(20)"`
	Name string `gorm:"type:varchar(60);not null"`
	Rule string `gorm:"type:enum('student','teacher');not null"`

	Courses []Course `gorm:"many2many:enrollments;foreignKey:ID;joinForeignKey:StudentID;References:ID;joinReferences:CourseID"`
}
