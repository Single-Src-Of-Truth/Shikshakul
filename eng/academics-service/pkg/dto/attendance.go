package dto

import (
	"time"

	"github.com/google/uuid"
)

type MarkAttendanceRequest struct {
	ClassID   string            `json:"class_id" binding:"required"`
	SectionID string            `json:"section_id" binding:"required"`
	Date      time.Time         `json:"date" binding:"required"`
	Students  []AttendanceEntry `json:"students" binding:"required"`
}

type AttendanceEntry struct {
	StudentID string `json:"student_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
	Remarks   string `json:"remarks"`
}

type AttendanceRegisterResponse struct {
	Date      string         `json:"date"`
	SectionID uuid.UUID      `json:"section_id"`
	Records   []RecordDetail `json:"records"`
	Summary   AttendanceStat `json:"summary"`
}

type RecordDetail struct {
	StudentID   uuid.UUID `json:"student_id"`
	StudentName string    `json:"student_name"`
	RollNo      string    `json:"roll_no"`
	Status      string    `json:"status"`
	Remarks     string    `json:"remarks"`
}

type AttendanceStat struct {
	Total   int `json:"total"`
	Present int `json:"present"`
	Absent  int `json:"absent"`
	Late    int `json:"late"`
}

type StudentAttendanceHistory struct {
	Date    string `json:"date"`
	Status  string `json:"status"`
	Remarks string `json:"remarks"`
}
