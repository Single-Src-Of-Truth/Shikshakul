package domain

import (
	"time"

	"github.com/google/uuid"
)

type ExamTerm struct {
	BaseEntity
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	AcademicYearID uuid.UUID `gorm:"type:uuid;not null;index" json:"academic_year_id"`

	Name      string    `gorm:"size:100;not null" json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	IsPublished bool `gorm:"default:false" json:"is_published"`
	IsActive    bool `gorm:"default:true" json:"is_active"`
}

type ExamSchedule struct {
	BaseEntity
	TenantID   uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ExamTermID uuid.UUID `gorm:"type:uuid;not null;index" json:"exam_term_id"`

	ClassID   uuid.UUID `gorm:"type:uuid;not null;index" json:"class_id"`
	SubjectID uuid.UUID `gorm:"type:uuid;not null;index" json:"subject_id"`

	ExamDate    time.Time `json:"exam_date"`
	StartTime   string    `gorm:"size:5;not null" json:"start_time"`
	DurationMin int       `json:"duration_min"`
	MaxMarks    float64   `gorm:"default:100" json:"max_marks"`
	PassMarks   float64   `gorm:"default:33" json:"pass_marks"`

	ExamTerm ExamTerm `gorm:"foreignKey:ExamTermID" json:"exam_term"`
	Class    Class    `gorm:"foreignKey:ClassID" json:"-"`
	Subject  Subject  `gorm:"foreignKey:SubjectID" json:"subject"`
}
