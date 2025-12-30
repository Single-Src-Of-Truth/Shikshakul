package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateFeeHeadRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type UpdateFeeHeadRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type FeeHeadResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

type CreateFeeStructureRequest struct {
	AcademicYearID string  `json:"academic_year_id" binding:"required"`
	ClassID        string  `json:"class_id" binding:"required"`
	FeeHeadID      string  `json:"fee_head_id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	Frequency      string  `json:"frequency" binding:"required"`
	DueDateDay     int     `json:"due_date_day"`
}

type UpdateFeeStructureRequest struct {
	Amount     float64 `json:"amount"`
	DueDateDay int     `json:"due_date_day"`
}

type FeeStructureResponse struct {
	ID          uuid.UUID `json:"id"`
	ClassID     uuid.UUID `json:"class_id"`
	FeeHeadID   uuid.UUID `json:"fee_head_id"`
	FeeHeadName string    `json:"fee_head_name,omitempty"`
	Amount      float64   `json:"amount"`
	Frequency   string    `json:"frequency"`
}

type GenerateDemandRequest struct {
	ClassID        string `json:"class_id" binding:"required"`
	AcademicYearID string `json:"academic_year_id" binding:"required"`
	Month          int    `json:"month" binding:"required"`
	Year           int    `json:"year" binding:"required"`
}

type CollectFeeRequest struct {
	StudentFeeID string  `json:"student_fee_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	PaymentMode  string  `json:"payment_mode" binding:"required"`
	Remarks      string  `json:"remarks"`
}

type StudentDueResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	DueDate   time.Time `json:"due_date"`
	AmountDue float64   `json:"amount_due"`
	Paid      float64   `json:"paid"`
	Balance   float64   `json:"balance"`
	Status    string    `json:"status"`
}

type TransactionHistoryResponse struct {
	ID              uuid.UUID `json:"id"`
	ReceiptNo       string    `json:"receipt_no"`
	Amount          float64   `json:"amount"`
	PaymentMode     string    `json:"payment_mode"`
	TransactionDate time.Time `json:"transaction_date"`
	FeeTitle        string    `json:"fee_title"`
}
