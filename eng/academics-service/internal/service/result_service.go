package service

import (
	"errors"
	"math"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
)

type ResultService struct {
	Repo        *repository.ResultRepository
	StudentRepo *repository.StudentRepository
}

func NewResultService(repo *repository.ResultRepository, studRepo *repository.StudentRepository) *ResultService {
	return &ResultService{Repo: repo, StudentRepo: studRepo}
}

func (s *ResultService) GenerateClassResults(tenantID uuid.UUID, req dto.GenerateResultRequest) error {
	classID := uuid.MustParse(req.ClassID)
	termID := uuid.MustParse(req.ExamTermID)

	students, err := s.StudentRepo.FindActiveStudentsByClass(classID)
	if err != nil {
		return err
	}

	for _, student := range students {
		marks, err := s.Repo.FetchStudentMarksForTerm(student.ID, termID)
		if err != nil {
			continue
		}

		var totalMax, totalObtained float64
		hasFail := false

		for _, m := range marks {
			totalMax += m.ExamSchedule.MaxMarks
			totalObtained += m.MarksObtained

			if (m.MarksObtained / m.ExamSchedule.MaxMarks * 100) < m.ExamSchedule.PassMarks {
				hasFail = true
			}
		}

		if totalMax == 0 {
			continue
		}

		percentage := (totalObtained / totalMax) * 100
		grade := calculateGrade(percentage)
		status := "PASS"
		if hasFail {
			status = "FAIL"
		}

		result := &domain.ExamResult{
			TenantID:      tenantID,
			ExamTermID:    termID,
			StudentID:     student.ID,
			ClassID:       classID,
			TotalMarks:    totalMax,
			ObtainedMarks: totalObtained,
			Percentage:    math.Round(percentage*100) / 100,
			Grade:         grade,
			ResultStatus:  status,
			IsPublished:   false,
		}

		s.Repo.SaveResult(result)
	}
	return nil
}

func (s *ResultService) GetReportCard(tenantID, studentID uuid.UUID, termIDStr string) (*dto.ReportCardResponse, error) {
	termID := uuid.MustParse(termIDStr)

	result, err := s.Repo.FindResult(studentID, termID)
	if err != nil {
		return nil, errors.New("result not generated yet")
	}

	if result.TenantID != tenantID {
		return nil, errors.New("unauthorized")
	}

	marks, _ := s.Repo.FetchStudentMarksForTerm(studentID, termID)

	resp := &dto.ReportCardResponse{
		StudentInfo: dto.StudentSummary{
			Name:        result.Student.FirstName + " " + result.Student.LastName,
			AdmissionNo: result.Student.AdmissionNo,
			ClassName:   result.Student.Class.Name,
		},
		ExamInfo: dto.ExamSummary{
			TermName: result.ExamTerm.Name,
		},
		Summary: dto.ResultSummary{
			TotalMax:      result.TotalMarks,
			TotalObtained: result.ObtainedMarks,
			Percentage:    result.Percentage,
			Grade:         result.Grade,
			Status:        result.ResultStatus,
		},
	}

	for _, m := range marks {
		resp.Subjects = append(resp.Subjects, dto.SubjectMark{
			SubjectName:   m.ExamSchedule.Subject.Name,
			MaxMarks:      m.ExamSchedule.MaxMarks,
			MarksObtained: m.MarksObtained,
			Grade:         calculateGrade((m.MarksObtained / m.ExamSchedule.MaxMarks) * 100),
			IsAbsent:      m.IsAbsent,
			Remarks:       m.Remarks,
		})
	}

	return resp, nil
}

func (s *ResultService) PublishResults(tenantID uuid.UUID, req dto.PublishResultRequest) error {
	return s.Repo.BulkPublishResults(uuid.MustParse(req.ExamTermID), uuid.MustParse(req.ClassID), req.Publish)
}

func calculateGrade(percentage float64) string {
	switch {
	case percentage >= 90:
		return "A+"
	case percentage >= 80:
		return "A"
	case percentage >= 70:
		return "B+"
	case percentage >= 60:
		return "B"
	case percentage >= 50:
		return "C"
	case percentage >= 40:
		return "D"
	default:
		return "F"
	}
}
