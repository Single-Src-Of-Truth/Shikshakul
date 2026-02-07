package service

import (
	"errors"
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
)

type CalendarService struct {
	Repo     *repository.CalendarRepository
	AcadRepo *repository.AcademicRepository
}

func NewCalendarService(repo *repository.CalendarRepository, acadRepo *repository.AcademicRepository) *CalendarService {
	return &CalendarService{Repo: repo, AcadRepo: acadRepo}
}

func (s *CalendarService) CreateEvent(tenantID uuid.UUID, req dto.CreateEventRequest) (*dto.EventResponse, error) {
	if req.StartDate.After(req.EndDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	if req.StartDate.Before(time.Now().Add(-24 * time.Hour)) {
		return nil, errors.New("cannot create events in the past")
	}

	yearUUID := uuid.MustParse(req.AcademicYearID)
	academicYear, err := s.AcadRepo.FindYearByID(yearUUID)
	if err != nil {
		return nil, errors.New("academic year not found")
	}

	if academicYear.TenantID != tenantID {
		return nil, errors.New("unauthorized access to academic year")
	}

	if req.StartDate.Before(academicYear.StartDate) || req.EndDate.After(academicYear.EndDate) {
		return nil, errors.New("event dates must fall within the selected academic year duration")
	}

	isOverlap, _ := s.Repo.CheckEventOverlap(tenantID, req.StartDate, req.EndDate, uuid.Nil)
	if isOverlap {
		return nil, errors.New("an event already exists during this time range")
	}

	event := &domain.CalendarEvent{
		TenantID:       tenantID,
		AcademicYearID: yearUUID,
		Title:          req.Title,
		Description:    req.Description,
		EventType:      domain.EventType(req.EventType),
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		IsHoliday:      req.IsHoliday,
	}

	if err := s.Repo.CreateEvent(event); err != nil {
		return nil, err
	}

	return mapEventToDTO(event), nil
}

func (s *CalendarService) UpdateEvent(tenantID, eventID uuid.UUID, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	event, err := s.Repo.FindEventByID(eventID)
	if err != nil || event.TenantID != tenantID {
		return nil, errors.New("event not found")
	}

	newStart := event.StartDate
	newEnd := event.EndDate
	datesChanged := false

	if !req.StartDate.IsZero() {
		newStart = req.StartDate
		datesChanged = true
	}
	if !req.EndDate.IsZero() {
		newEnd = req.EndDate
		datesChanged = true
	}

	if datesChanged {
		if newStart.After(newEnd) {
			return nil, errors.New("start date cannot be after end date")
		}

		academicYear, err := s.AcadRepo.FindYearByID(event.AcademicYearID)
		if err == nil {
			if newStart.Before(academicYear.StartDate) || newEnd.After(academicYear.EndDate) {
				return nil, errors.New("updated dates must remain within the academic year duration")
			}
		}

		isOverlap, _ := s.Repo.CheckEventOverlap(tenantID, newStart, newEnd, eventID)
		if isOverlap {
			return nil, errors.New("updated dates conflict with another event")
		}

		event.StartDate = newStart
		event.EndDate = newEnd
	}

	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.EventType != "" {
		event.EventType = domain.EventType(req.EventType)
	}
	if req.IsHoliday != nil {
		event.IsHoliday = *req.IsHoliday
	}

	if err := s.Repo.UpdateEvent(event); err != nil {
		return nil, err
	}
	return mapEventToDTO(event), nil
}

func (s *CalendarService) GetMonthEvents(tenantID uuid.UUID, filter dto.CalendarFilter) ([]dto.EventResponse, error) {
	start := time.Date(filter.Year, time.Month(filter.Month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1).Add(23 * time.Hour)

	events, err := s.Repo.FindEventsByRange(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	var dtos []dto.EventResponse
	for _, e := range events {
		dtos = append(dtos, *mapEventToDTO(&e))
	}
	return dtos, nil
}

func (s *CalendarService) DeleteEvent(tenantID, id uuid.UUID) error {
	event, err := s.Repo.FindEventByID(id)
	if err != nil || event.TenantID != tenantID {
		return errors.New("event not found")
	}
	return s.Repo.DeleteEvent(id)
}

func mapEventToDTO(e *domain.CalendarEvent) *dto.EventResponse {
	return &dto.EventResponse{
		ID:          e.ID.String(),
		Title:       e.Title,
		Description: e.Description,
		EventType:   string(e.EventType),
		StartDate:   e.StartDate.Format(time.RFC3339),
		EndDate:     e.EndDate.Format(time.RFC3339),
		IsHoliday:   e.IsHoliday,
	}
}
