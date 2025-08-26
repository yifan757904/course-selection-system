package model

type Admin struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	JobNo    string `gorm:"unique;size:20;not null"`
	Password string `gorm:"size:60;not null"`
}
