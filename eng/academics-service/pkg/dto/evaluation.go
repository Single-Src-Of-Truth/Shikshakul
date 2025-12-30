package dto

import "github.com/google/uuid"

type SubmitMarksRequest struct {
	ExamScheduleID string             `json:"exam_schedule_id" binding:"required"`
	Marks          []StudentMarkEntry `json:"marks" binding:"required"`
	Finalize       bool               `json:"finalize"`
}

type StudentMarkEntry struct {
	StudentID     string  `json:"student_id" binding:"required"`
	MarksObtained float64 `json:"marks_obtained"`
	IsAbsent      bool    `json:"is_absent"`
	Remarks       string  `json:"remarks"`
}

type EvaluationSheetResponse struct {
	ExamScheduleID uuid.UUID    `json:"exam_schedule_id"`
	SubjectName    string       `json:"subject_name"`
	MaxMarks       float64      `json:"max_marks"`
	Status         string       `json:"status"`
	Entries        []SheetEntry `json:"entries"`
}

type SheetEntry struct {
	StudentID     uuid.UUID `json:"student_id"`
	StudentName   string    `json:"student_name"`
	RollNo        string    `json:"roll_no"`
	MarksObtained float64   `json:"marks_obtained"`
	IsAbsent      bool      `json:"is_absent"`
	Remarks       string    `json:"remarks"`
}
