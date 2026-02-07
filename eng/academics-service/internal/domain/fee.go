package domain

import (
	"time"

	"github.com/google/uuid"
)

type FeeHead struct {
	BaseEntity
	TenantID uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name     string    `gorm:"size:100;not null" json:"name"`
	Type     string    `gorm:"size:20;default:'RECURRING'" json:"type"`
}

type FeeStructure struct {
	BaseEntity
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	AcademicYearID uuid.UUID `gorm:"type:uuid;not null;index" json:"academic_year_id"`
	ClassID        uuid.UUID `gorm:"type:uuid;not null;index" json:"class_id"`
	FeeHeadID      uuid.UUID `gorm:"type:uuid;not null;index" json:"fee_head_id"`

	Amount     float64      `gorm:"not null" json:"amount"`
	Frequency  FeeFrequency `gorm:"size:20;default:'MONTHLY'" json:"frequency"`
	DueDateDay int          `gorm:"default:10" json:"due_date_day"`

	FeeHead FeeHead `gorm:"foreignKey:FeeHeadID" json:"fee_head"`
}

type StudentFee struct {
	BaseEntity
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	StudentID      uuid.UUID `gorm:"type:uuid;not null;index" json:"student_id"`
	FeeStructureID uuid.UUID `gorm:"type:uuid;not null;index" json:"fee_structure_id"`

	Title   string    `gorm:"size:100" json:"title"`
	DueDate time.Time `json:"due_date"`

	AmountDue  float64   `gorm:"not null" json:"amount_due"`
	AmountPaid float64   `gorm:"default:0" json:"amount_paid"`
	Status     FeeStatus `gorm:"size:20;default:'PENDING'" json:"status"`
}

type FeeTransaction struct {
	BaseEntity
	TenantID     uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	StudentID    uuid.UUID `gorm:"type:uuid;not null;index" json:"student_id"`
	StudentFeeID uuid.UUID `gorm:"type:uuid;not null;index" json:"student_fee_id"`

	Amount          float64     `gorm:"not null" json:"amount"`
	PaymentMode     PaymentMode `gorm:"size:20;not null" json:"payment_mode"`
	TransactionDate time.Time   `gorm:"not null" json:"transaction_date"`
	ReceiptNo       string      `gorm:"size:50;not null" json:"receipt_no"`
	Remarks         string      `gorm:"size:255" json:"remarks"`

	CollectedBy uuid.UUID `gorm:"type:uuid" json:"collected_by"`
}
