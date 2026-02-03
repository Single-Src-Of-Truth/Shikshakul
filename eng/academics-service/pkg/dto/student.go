package dto

import "encoding/json"

type CreateSchemaRequest struct {
	Structure []map[string]interface{} `json:"structure" binding:"required"`
}

type SubmitApplicationRequest struct {
	AcademicYearID string                 `json:"academic_year_id" binding:"required"`
	ClassID        string                 `json:"class_id" binding:"required"`
	FirstName      string                 `json:"first_name" binding:"required"`
	LastName       string                 `json:"last_name" binding:"required"`
	Email          string                 `json:"email"`
	Mobile         string                 `json:"mobile"`
	ProfileData    map[string]interface{} `json:"profile_data"`
	Documents      map[string]interface{} `json:"documents"`
}

type UpdateStudentRequest struct {
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	Email       string                 `json:"email"`
	Mobile      string                 `json:"mobile"`
	ProfileData map[string]interface{} `json:"profile_data"`
	Documents   map[string]interface{} `json:"documents"`
}

type StudentResponse struct {
	ID          string          `json:"id"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	Status      string          `json:"status"`
	AdmissionNo string          `json:"admission_no,omitempty"`
	ClassID     string          `json:"class_id"`
	ClassName   string          `json:"class_name,omitempty"`
	ProfileData json.RawMessage `json:"profile_data"`
	Documents   json.RawMessage `json:"documents"`
}

type ApproveStudentRequest struct {
	Action    string `json:"action" binding:"required"`
	SectionID string `json:"section_id"`
}

type StudentFilter struct {
	ClassID string `form:"class_id"`
	Status  string `form:"status"`
}

type ConfigureSequenceRequest struct {
	AcademicYearID string `json:"academic_year_id" binding:"required"`
	Prefix         string `json:"prefix" binding:"required"`
	StartingCount  *int   `json:"starting_count"`
}
