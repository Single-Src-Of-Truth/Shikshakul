package dto

import (
	"time"

	"github.com/google/uuid"
)

type IDResponse struct {
	ID string `json:"id"`
}

type CreateYearRequest struct {
	Name      string    `json:"name" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	IsCurrent bool      `json:"is_current"`
}

type UpdateYearRequest struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsCurrent *bool     `json:"is_current"`
}

type AcademicYearResponse struct {
	ID        uuid.UUID `json:"academic_year_id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsCurrent bool      `json:"is_current"`
}

type CreateClassRequest struct {
	Name      string `json:"name" binding:"required"`
	SortOrder int    `json:"sort_order" binding:"required"`
}

type UpdateClassRequest struct {
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

type ClassResponse struct {
	ID        uuid.UUID `json:"class_id"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
}

type CreateSectionRequest struct {
	Name     string `json:"name" binding:"required"`
	Capacity int    `json:"capacity"`
}

type UpdateSectionRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type SectionResponse struct {
	ID       uuid.UUID `json:"section_id"`
	ClassID  uuid.UUID `json:"class_id"`
	Name     string    `json:"name"`
	Capacity int       `json:"capacity"`
}

type CreateSubjectRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code"`
	Type string `json:"type"`
}

type UpdateSubjectRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Type string `json:"type"`
}

type SubjectResponse struct {
	ID   uuid.UUID `json:"subject_id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
	Type string    `json:"type"`
}

type AssignSubjectsRequest struct {
	Subjects []SubjectMappingRequest `json:"subjects" binding:"required"`
}

type SubjectMappingRequest struct {
	SubjectID      string `json:"subject_id" binding:"required"`
	IsOptional     bool   `json:"is_optional"`
	WeeklyLectures int    `json:"weekly_lectures"`
}
