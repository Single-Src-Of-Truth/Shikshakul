package domain

import (
	"time"

	"github.com/google/uuid"
)

type CalendarEvent struct {
	BaseEntity
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	AcademicYearID uuid.UUID `gorm:"type:uuid;not null;index" json:"academic_year_id"`

	Title       string    `gorm:"size:100;not null" json:"title"`
	Description string    `gorm:"size:255" json:"description"`
	EventType   EventType `gorm:"size:20;not null" json:"event_type"`

	StartDate time.Time `gorm:"not null;index" json:"start_date"`
	EndDate   time.Time `gorm:"not null" json:"end_date"`

	IsHoliday bool `gorm:"default:false" json:"is_holiday"`
}
