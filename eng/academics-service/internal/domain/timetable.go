package domain

import (
	"github.com/google/uuid"
)

type WeeklyRoutine struct {
	BaseEntity
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	AcademicYearID uuid.UUID `gorm:"type:uuid;not null;index" json:"academic_year_id"`

	ClassID   uuid.UUID `gorm:"type:uuid;not null;index" json:"class_id"`
	SectionID uuid.UUID `gorm:"type:uuid;not null;index" json:"section_id"`
	SubjectID uuid.UUID `gorm:"type:uuid;not null;index" json:"subject_id"`
	TeacherID uuid.UUID `gorm:"type:uuid;not null;index" json:"teacher_id"`

	DayOfWeek  DayOfWeek `gorm:"size:20;not null" json:"day_of_week"`
	StartTime  string    `gorm:"size:5;not null" json:"start_time"`
	EndTime    string    `gorm:"size:5;not null" json:"end_time"`
	RoomNumber string    `gorm:"size:20" json:"room_number"`

	Class   Class   `gorm:"foreignKey:ClassID" json:"-"`
	Section Section `gorm:"foreignKey:SectionID" json:"section"`
	Subject Subject `gorm:"foreignKey:SubjectID" json:"subject"`
	Teacher Teacher `gorm:"foreignKey:TeacherID" json:"teacher"`
}
