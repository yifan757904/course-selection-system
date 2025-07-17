package model

type Course struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	TeacherID      string `json:"teacher_id"`
	TeacherName    string `json:"teacher_name"`
	Remarks        string `json:"remarks"`
	Student_maxnum int    `json:"student_maxnum"`
	Time_max       int    `json:"time_max"`
	Time_min       int    `json:"time_min"`
}
