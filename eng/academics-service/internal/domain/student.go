package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

type AdmissionSchema struct {
	BaseEntity
	TenantID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Structure json.RawMessage `gorm:"type:jsonb;not null" json:"structure"`
	Version   int             `gorm:"default:1" json:"version"`
	IsActive  bool            `gorm:"default:true" json:"is_active"`
}

type Student struct {
	BaseEntity
	TenantID uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	UserID   *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`

	AcademicYearID uuid.UUID  `gorm:"type:uuid;not null;index" json:"academic_year_id"`
	ClassID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"class_id"`
	SectionID      *uuid.UUID `gorm:"type:uuid;index" json:"section_id,omitempty"`

	FirstName string `gorm:"size:100;not null" json:"first_name"`
	LastName  string `gorm:"size:100;not null" json:"last_name"`
	Email     string `gorm:"size:100" json:"email"`
	Mobile    string `gorm:"size:20" json:"mobile"`

	ProfileData json.RawMessage `gorm:"type:jsonb" json:"profile_data"`
	Documents   json.RawMessage `gorm:"type:jsonb" json:"documents"`

	Status      Status `gorm:"size:20;default:'PENDING'" json:"status"`
	AdmissionNo string `gorm:"size:50;index" json:"admission_no"`

	Class   Class   `gorm:"foreignKey:ClassID" json:"class"`
	Section Section `gorm:"foreignKey:SectionID" json:"section"`
}
