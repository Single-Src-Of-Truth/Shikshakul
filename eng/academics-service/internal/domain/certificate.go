package domain

import (
	"time"

	"github.com/google/uuid"
)

type CertificateRecord struct {
	BaseEntity
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	StudentID uuid.UUID `gorm:"type:uuid;not null;index" json:"student_id"`

	Type          CertificateType `gorm:"size:50;not null" json:"type"`
	CertificateNo string          `gorm:"size:50;not null;index" json:"certificate_no"`
	IssueDate     time.Time       `json:"issue_date"`

	Remarks  string    `gorm:"size:255" json:"remarks"`
	IssuedBy uuid.UUID `gorm:"type:uuid" json:"issued_by"`

	Student Student `gorm:"foreignKey:StudentID" json:"student"`
}
