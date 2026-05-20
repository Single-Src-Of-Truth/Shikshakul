package dto

import "time"

type CreateEventRequest struct {
	AcademicYearID string    `json:"academic_year_id" binding:"required"`
	Title          string    `json:"title" binding:"required"`
	Description    string    `json:"description"`
	EventType      string    `json:"event_type" binding:"required"`
	StartDate      time.Time `json:"start_date" binding:"required"`
	EndDate        time.Time `json:"end_date" binding:"required"`
	IsHoliday      bool      `json:"is_holiday"`
}

type UpdateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventType   string    `json:"event_type"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	IsHoliday   *bool     `json:"is_holiday"`
}

type EventResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EventType   string `json:"event_type"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	IsHoliday   bool   `json:"is_holiday"`
}

type CalendarFilter struct {
	Month int `form:"month"`
	Year  int `form:"year"`
}
