package domain

import (
	"time"

	"github.com/google/uuid"
)

type AcademicYear struct {
	BaseEntity
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsCurrent bool      `gorm:"default:false" json:"is_current"`
}

type Class struct {
	BaseEntity
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	SortOrder int       `gorm:"not null" json:"sort_order"`
	Sections  []Section `gorm:"foreignKey:ClassID" json:"sections,omitempty"`
}

type Section struct {
	BaseEntity
	ClassID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"class_id"`
	Name           string     `gorm:"size:50;not null" json:"name"`
	Capacity       int        `json:"capacity"`
	ClassTeacherID *uuid.UUID `gorm:"type:uuid;index" json:"class_teacher_id,omitempty"`
}

type Subject struct {
	BaseEntity
	TenantID uuid.UUID   `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name     string      `gorm:"size:100;not null" json:"name"`
	Code     string      `gorm:"size:20" json:"code"`
	Type     SubjectType `gorm:"size:20;default:'THEORY'" json:"type"`
}

type ClassSubject struct {
	BaseEntity
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ClassID        uuid.UUID `gorm:"type:uuid;not null;index" json:"class_id"`
	SubjectID      uuid.UUID `gorm:"type:uuid;not null;index" json:"subject_id"`
	IsOptional     bool      `gorm:"default:false" json:"is_optional"`
	WeeklyLectures int       `gorm:"default:0" json:"weekly_lectures"`

	Subject Subject `gorm:"foreignKey:SubjectID" json:"subject"`
}
