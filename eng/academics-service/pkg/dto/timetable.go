package dto

import "github.com/google/uuid"

type CreateRoutineRequest struct {
	AcademicYearID string `json:"academic_year_id" binding:"required"`
	ClassID        string `json:"class_id" binding:"required"`
	SectionID      string `json:"section_id" binding:"required"`
	SubjectID      string `json:"subject_id" binding:"required"`
	TeacherID      string `json:"teacher_id" binding:"required"`

	DayOfWeek  string `json:"day_of_week" binding:"required"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
	RoomNumber string `json:"room_number"`
}

type RoutineResponse struct {
	ID          uuid.UUID `json:"id"`
	DayOfWeek   string    `json:"day_of_week"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	SubjectName string    `json:"subject_name"`
	TeacherName string    `json:"teacher_name"`
	RoomNumber  string    `json:"room_number,omitempty"`
}

type TimetableFilter struct {
	SectionID string `form:"section_id"`
	TeacherID string `form:"teacher_id"`
	Day       string `form:"day"`
}
