package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateExamTermRequest struct {
	AcademicYearID string    `json:"academic_year_id" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	StartDate      time.Time `json:"start_date" binding:"required"`
	EndDate        time.Time `json:"end_date" binding:"required"`
}

type ExamTermResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
}

type CreateExamScheduleRequest struct {
	ExamTermID  string    `json:"exam_term_id" binding:"required"`
	ClassID     string    `json:"class_id" binding:"required"`
	SubjectID   string    `json:"subject_id" binding:"required"`
	ExamDate    time.Time `json:"exam_date" binding:"required"`
	StartTime   string    `json:"start_time" binding:"required"`
	DurationMin int       `json:"duration_min" binding:"required"`
	MaxMarks    float64   `json:"max_marks" binding:"required"`
	PassMarks   float64   `json:"pass_marks" binding:"required"`
}

type ExamScheduleResponse struct {
	ID          uuid.UUID `json:"id"`
	SubjectName string    `json:"subject_name"`
	ExamDate    time.Time `json:"exam_date"`
	StartTime   string    `json:"start_time"`
	DurationMin int       `json:"duration_min"`
	MaxMarks    float64   `json:"max_marks"`
}
