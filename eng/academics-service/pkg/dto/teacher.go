package dto

import "encoding/json"

type OnboardTeacherRequest struct {
	FirstName   string                 `json:"first_name" binding:"required"`
	LastName    string                 `json:"last_name" binding:"required"`
	Email       string                 `json:"email" binding:"required,email"`
	Mobile      string                 `json:"mobile"`
	ProfileData map[string]interface{} `json:"profile_data"`
	Documents   map[string]interface{} `json:"documents"`
}

type UpdateTeacherRequest struct {
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	Email       string                 `json:"email"`
	Mobile      string                 `json:"mobile"`
	ProfileData map[string]interface{} `json:"profile_data"`
	Documents   map[string]interface{} `json:"documents"`
}

type TeacherResponse struct {
	ID             string          `json:"id"`
	EmployeeID     string          `json:"employee_id,omitempty"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	Email          string          `json:"email"`
	Status         string          `json:"status"`
	IsClassTeacher bool            `json:"is_class_teacher"`
	ProfileData    json.RawMessage `json:"profile_data"`
	Documents      json.RawMessage `json:"documents"`
}

type ApproveTeacherRequest struct {
	Action string `json:"action" binding:"required"`
}

type AssignClassTeacherRequest struct {
	TeacherID string `json:"teacher_id" binding:"required"`
}

type AssignSubjectTeacherRequest struct {
	AcademicYearID string `json:"academic_year_id" binding:"required"`
	SubjectID      string `json:"subject_id" binding:"required"`
	TeacherID      string `json:"teacher_id" binding:"required"`
}

type TeacherFilter struct {
	Status string `form:"status"`
}
