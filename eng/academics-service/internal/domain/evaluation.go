package domain

import (
	"github.com/google/uuid"
)

type StudentMark struct {
	BaseEntity
	TenantID uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`

	ExamScheduleID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_marks_unique" json:"exam_schedule_id"`

	StudentID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_marks_unique" json:"student_id"`

	MarksObtained float64 `gorm:"not null" json:"marks_obtained"`
	IsAbsent      bool    `gorm:"default:false" json:"is_absent"`
	Remarks       string  `gorm:"size:255" json:"remarks"`

	GradedBy uuid.UUID        `gorm:"type:uuid" json:"graded_by"`
	Status   EvaluationStatus `gorm:"size:20;default:'DRAFT'" json:"status"`

	ExamSchedule ExamSchedule `gorm:"foreignKey:ExamScheduleID" json:"-"`
	Student      Student      `gorm:"foreignKey:StudentID" json:"student"`
}
