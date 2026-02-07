package service

import (
	"errors"
	"fmt"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
)

type EvaluationService struct {
	Repo        *repository.EvaluationRepository
	TeacherRepo *repository.TeacherRepository
}

func NewEvaluationService(repo *repository.EvaluationRepository, tRepo *repository.TeacherRepository) *EvaluationService {
	return &EvaluationService{Repo: repo, TeacherRepo: tRepo}
}

func (s *EvaluationService) GetMarkSheet(tenantID uuid.UUID, scheduleIDStr string) (*dto.EvaluationSheetResponse, error) {
	scheduleID := uuid.MustParse(scheduleIDStr)

	schedule, err := s.Repo.FindScheduleByID(scheduleID)
	if err != nil {
		return nil, errors.New("exam schedule not found")
	}

	students, err := s.Repo.FindStudentsByExamSchedule(scheduleID)
	if err != nil {
		return nil, err
	}

	marks, err := s.Repo.FindMarksBySchedule(scheduleID)
	if err != nil {
		return nil, err
	}

	marksMap := make(map[uuid.UUID]domain.StudentMark)
	status := "DRAFT"
	for _, m := range marks {
		marksMap[m.StudentID] = m
		if m.Status == domain.EvaluationStatusSubmitted {
			status = "SUBMITTED"
		}
	}

	var entries []dto.SheetEntry
	for _, stud := range students {
		entry := dto.SheetEntry{
			StudentID:   stud.ID,
			StudentName: stud.FirstName + " " + stud.LastName,
			RollNo:      stud.AdmissionNo,
		}

		if val, ok := marksMap[stud.ID]; ok {
			entry.MarksObtained = val.MarksObtained
			entry.IsAbsent = val.IsAbsent
			entry.Remarks = val.Remarks
		}
		entries = append(entries, entry)
	}

	return &dto.EvaluationSheetResponse{
		ExamScheduleID: schedule.ID,
		SubjectName:    schedule.Subject.Name,
		MaxMarks:       schedule.MaxMarks,
		Status:         status,
		Entries:        entries,
	}, nil
}

func (s *EvaluationService) SubmitMarks(tenantID, teacherID uuid.UUID, req dto.SubmitMarksRequest) error {
	scheduleID := uuid.MustParse(req.ExamScheduleID)

	schedule, err := s.Repo.FindScheduleByID(scheduleID)
	if err != nil {
		return errors.New("exam schedule not found")
	}

	// Optional: Check if this Teacher teaches this Subject in this Class?
	// (You can use TeacherRepo logic here for strict security)

	var marksBatch []domain.StudentMark
	status := domain.EvaluationStatusDraft
	if req.Finalize {
		status = domain.EvaluationStatusSubmitted
	}

	for _, entry := range req.Marks {
		if entry.MarksObtained > schedule.MaxMarks {
			return fmt.Errorf("student %s marks (%v) exceed max marks (%v)", entry.StudentID, entry.MarksObtained, schedule.MaxMarks)
		}

		marksBatch = append(marksBatch, domain.StudentMark{
			TenantID:       tenantID,
			ExamScheduleID: scheduleID,
			StudentID:      uuid.MustParse(entry.StudentID),
			MarksObtained:  entry.MarksObtained,
			IsAbsent:       entry.IsAbsent,
			Remarks:        entry.Remarks,
			GradedBy:       teacherID,
			Status:         status,
		})
	}

	if len(marksBatch) == 0 {
		return nil
	}

	return s.Repo.UpsertMarks(marksBatch)
}
