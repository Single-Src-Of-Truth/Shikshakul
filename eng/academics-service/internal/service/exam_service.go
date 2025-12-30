package service

import (
	"errors"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
)

type ExamService struct {
	Repo *repository.ExamRepository
}

func NewExamService(repo *repository.ExamRepository) *ExamService {
	return &ExamService{Repo: repo}
}

func (s *ExamService) CreateTerm(tenantID uuid.UUID, req dto.CreateExamTermRequest) (*domain.ExamTerm, error) {
	if req.EndDate.Before(req.StartDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	term := &domain.ExamTerm{
		TenantID:       tenantID,
		AcademicYearID: uuid.MustParse(req.AcademicYearID),
		Name:           req.Name,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		IsActive:       true,
	}
	err := s.Repo.CreateTerm(term)
	return term, err
}

func (s *ExamService) GetTermsByYear(yearIDStr string) ([]dto.ExamTermResponse, error) {
	yearID := uuid.MustParse(yearIDStr)
	terms, err := s.Repo.FindTermsByYear(yearID)
	if err != nil {
		return nil, err
	}

	var dtos []dto.ExamTermResponse
	for _, t := range terms {
		dtos = append(dtos, dto.ExamTermResponse{
			ID: t.ID, Name: t.Name, StartDate: t.StartDate, EndDate: t.EndDate, IsActive: t.IsActive,
		})
	}
	return dtos, nil
}

func (s *ExamService) CreateSchedule(tenantID uuid.UUID, req dto.CreateExamScheduleRequest) (*domain.ExamSchedule, error) {
	termID := uuid.MustParse(req.ExamTermID)
	classID := uuid.MustParse(req.ClassID)

	term, err := s.Repo.FindTermByID(termID)
	if err != nil {
		return nil, errors.New("exam term not found")
	}

	if req.ExamDate.Before(term.StartDate) || req.ExamDate.After(term.EndDate) {
		return nil, errors.New("exam date is outside the exam term range")
	}

	conflict, _ := s.Repo.CheckScheduleConflict(classID, req.ExamDate)
	if conflict {
		return nil, errors.New("class already has an exam scheduled on this date")
	}

	schedule := &domain.ExamSchedule{
		TenantID:    tenantID,
		ExamTermID:  termID,
		ClassID:     classID,
		SubjectID:   uuid.MustParse(req.SubjectID),
		ExamDate:    req.ExamDate,
		StartTime:   req.StartTime,
		DurationMin: req.DurationMin,
		MaxMarks:    req.MaxMarks,
		PassMarks:   req.PassMarks,
	}

	err = s.Repo.CreateSchedule(schedule)
	return schedule, err
}

func (s *ExamService) GetSchedule(termIDStr, classIDStr string) ([]dto.ExamScheduleResponse, error) {
	termID := uuid.MustParse(termIDStr)
	classID := uuid.MustParse(classIDStr)

	schedules, err := s.Repo.FindSchedulesByTermAndClass(termID, classID)
	if err != nil {
		return nil, err
	}

	var dtos []dto.ExamScheduleResponse
	for _, sc := range schedules {
		dtos = append(dtos, dto.ExamScheduleResponse{
			ID: sc.ID, SubjectName: sc.Subject.Name, ExamDate: sc.ExamDate, StartTime: sc.StartTime, DurationMin: sc.DurationMin, MaxMarks: sc.MaxMarks,
		})
	}
	return dtos, nil
}

func (s *ExamService) DeleteSchedule(tenantID, id uuid.UUID) error {
	return s.Repo.DeleteSchedule(id)
}
