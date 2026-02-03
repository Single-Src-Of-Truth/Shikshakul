package service

import (
	"errors"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
)

type TimetableService struct {
	Repo *repository.TimetableRepository
}

func NewTimetableService(repo *repository.TimetableRepository) *TimetableService {
	return &TimetableService{Repo: repo}
}

func (s *TimetableService) CreateRoutine(tenantID uuid.UUID, req dto.CreateRoutineRequest) (*dto.RoutineResponse, error) {
	if req.StartTime >= req.EndTime {
		return nil, errors.New("start time must be before end time")
	}

	teacherID := uuid.MustParse(req.TeacherID)
	sectionID := uuid.MustParse(req.SectionID)

	if err := s.Repo.CheckConflicts(tenantID, teacherID, sectionID, req.DayOfWeek, req.StartTime, req.EndTime); err != nil {
		return nil, err
	}

	routine := &domain.WeeklyRoutine{
		TenantID:       tenantID,
		AcademicYearID: uuid.MustParse(req.AcademicYearID),
		ClassID:        uuid.MustParse(req.ClassID),
		SectionID:      sectionID,
		SubjectID:      uuid.MustParse(req.SubjectID),
		TeacherID:      teacherID,
		DayOfWeek:      domain.DayOfWeek(req.DayOfWeek),
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		RoomNumber:     req.RoomNumber,
	}

	if err := s.Repo.CreateRoutine(routine); err != nil {
		return nil, err
	}

	return &dto.RoutineResponse{
		ID:        routine.ID,
		DayOfWeek: string(routine.DayOfWeek),
		StartTime: routine.StartTime,
		EndTime:   routine.EndTime,
	}, nil
}

func (s *TimetableService) GetTimetable(tenantID uuid.UUID, filter dto.TimetableFilter) ([]dto.RoutineResponse, error) {
	var secID, teachID uuid.UUID
	if filter.SectionID != "" {
		secID = uuid.MustParse(filter.SectionID)
	}
	if filter.TeacherID != "" {
		teachID = uuid.MustParse(filter.TeacherID)
	}

	routines, err := s.Repo.FindRoutines(tenantID, secID, teachID, filter.Day)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.RoutineResponse, len(routines))
	for i, r := range routines {
		dtos[i] = dto.RoutineResponse{
			ID:          r.ID,
			DayOfWeek:   string(r.DayOfWeek),
			StartTime:   r.StartTime,
			EndTime:     r.EndTime,
			SubjectName: r.Subject.Name,
			TeacherName: r.Teacher.FirstName + " " + r.Teacher.LastName,
			RoomNumber:  r.RoomNumber,
		}
	}
	return dtos, nil
}

func (s *TimetableService) DeleteRoutine(tenantID, id uuid.UUID) error {
	return s.Repo.DeleteRoutine(id)
}
