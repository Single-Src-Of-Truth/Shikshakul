package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

type TeacherSchema struct {
	BaseEntity
	TenantID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Structure json.RawMessage `gorm:"type:jsonb;not null" json:"structure"`
	Version   int             `gorm:"default:1" json:"version"`
	IsActive  bool            `gorm:"default:true" json:"is_active"`
}

type Teacher struct {
	BaseEntity
	TenantID uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	UserID   *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`

	FirstName string `gorm:"size:100;not null" json:"first_name"`
	LastName  string `gorm:"size:100;not null" json:"last_name"`
	Email     string `gorm:"size:100;not null" json:"email"`
	Mobile    string `gorm:"size:20" json:"mobile"`

	ProfileData json.RawMessage `gorm:"type:jsonb" json:"profile_data"`
	Documents   json.RawMessage `gorm:"type:jsonb" json:"documents"`

	Status         Status `gorm:"size:20;default:'PENDING'" json:"status"`
	EmployeeID     string `gorm:"size:50;index" json:"employee_id"`
	IsClassTeacher bool   `gorm:"default:false" json:"is_class_teacher"`
}

type SubjectAllocation struct {
	BaseEntity
	TenantID uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`

	AcademicYearID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_alloc_unique" json:"academic_year_id"`
	SectionID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_alloc_unique" json:"section_id"`
	SubjectID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_alloc_unique" json:"subject_id"`

	TeacherID uuid.UUID `gorm:"type:uuid;not null;index" json:"teacher_id"`

	Teacher Teacher `gorm:"foreignKey:TeacherID" json:"teacher"`
	Subject Subject `gorm:"foreignKey:SubjectID" json:"subject"`
}
