package domain

import (
	"github.com/google/uuid"
)

type ExamResult struct {
	BaseEntity
	TenantID uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`

	ExamTermID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_result_unique" json:"exam_term_id"`
	StudentID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_result_unique" json:"student_id"`
	ClassID    uuid.UUID `gorm:"type:uuid;not null;index" json:"class_id"`

	TotalMarks    float64 `json:"total_marks"`
	ObtainedMarks float64 `json:"obtained_marks"`
	Percentage    float64 `json:"percentage"`
	Grade         string  `gorm:"size:5" json:"grade"`
	ResultStatus  string  `gorm:"size:20" json:"result_status"`

	Remarks     string `gorm:"size:255" json:"remarks"`
	IsPublished bool   `gorm:"default:false" json:"is_published"`

	ExamTerm ExamTerm `gorm:"foreignKey:ExamTermID" json:"exam_term"`
	Student  Student  `gorm:"foreignKey:StudentID" json:"student"`
}
