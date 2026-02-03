package dto

import "time"

type GenerateIDCardRequest struct {
	ClassID string `json:"class_id" binding:"required"`
}

type IDCardResponse struct {
	StudentID     string `json:"student_id"`
	AdmissionNo   string `json:"admission_no"`
	FullName      string `json:"full_name"`
	ClassDetails  string `json:"class_details"`
	DOB           string `json:"dob"`
	BloodGroup    string `json:"blood_group"`
	FatherName    string `json:"father_name"`
	ContactNumber string `json:"contact_number"`
	Address       string `json:"address"`
	PhotoURL      string `json:"photo_url"`
}

type IssueTCRequest struct {
	StudentID    string    `json:"student_id" binding:"required"`
	LeavingDate  time.Time `json:"leaving_date" binding:"required"`
	Reason       string    `json:"reason" binding:"required"`
	Conduct      string    `json:"conduct" binding:"required"`
	MarkInactive bool      `json:"mark_inactive"`
}

type TCResponse struct {
	CertificateNo string `json:"certificate_no"`
	IssueDate     string `json:"issue_date"`

	StudentName string `json:"student_name"`
	FatherName  string `json:"father_name"`
	MotherName  string `json:"mother_name"`
	AdmissionNo string `json:"admission_no"`
	DOB         string `json:"dob"`
	Nationality string `json:"nationality"`

	ClassJoined  string `json:"class_joined"`
	LastClass    string `json:"last_class"`
	ResultStatus string `json:"result_status"`

	Reason  string `json:"reason"`
	Conduct string `json:"conduct"`
}

type IssueBonafideRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Purpose   string `json:"purpose" binding:"required"`
}

type BonafideResponse struct {
	CertificateNo string `json:"certificate_no"`
	IssueDate     string `json:"issue_date"`
	StudentName   string `json:"student_name"`
	AdmissionNo   string `json:"admission_no"`
	ClassDetails  string `json:"class_details"`
	AcademicYear  string `json:"academic_year"`
	Purpose       string `json:"purpose"`
}
