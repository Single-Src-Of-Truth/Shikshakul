package domain

import (
	"time"

	"github.com/google/uuid"
)

type StudentAttendance struct {
	BaseEntity
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	AcademicYearID uuid.UUID `gorm:"type:uuid;not null;index" json:"academic_year_id"`

	ClassID   uuid.UUID `gorm:"type:uuid;not null;index" json:"class_id"`
	SectionID uuid.UUID `gorm:"type:uuid;not null;index" json:"section_id"`

	StudentID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_att_student_date" json:"student_id"`

	Date    time.Time        `gorm:"type:date;not null;uniqueIndex:idx_att_student_date" json:"date"`
	Status  AttendanceStatus `gorm:"size:20;not null" json:"status"`
	Remarks string           `gorm:"size:255" json:"remarks"`

	MarkedBy uuid.UUID `gorm:"type:uuid" json:"marked_by"`

	Student Student `gorm:"foreignKey:StudentID" json:"student"`
}
