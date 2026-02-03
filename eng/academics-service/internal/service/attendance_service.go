package service

import (
	"errors"
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
)

type AttendanceService struct {
	Repo        *repository.AttendanceRepository
	StudentRepo *repository.StudentRepository
	AcadRepo    *repository.AcademicRepository
}

func NewAttendanceService(repo *repository.AttendanceRepository, sRepo *repository.StudentRepository, aRepo *repository.AcademicRepository) *AttendanceService {
	return &AttendanceService{Repo: repo, StudentRepo: sRepo, AcadRepo: aRepo}
}

func (s *AttendanceService) MarkAttendance(tenantID, teacherID uuid.UUID, req dto.MarkAttendanceRequest) error {
	today := time.Now().Truncate(24 * time.Hour)
	reqDate := req.Date.Truncate(24 * time.Hour)

	if reqDate.After(today) {
		return errors.New("cannot mark attendance for future dates")
	}

	year, err := s.AcadRepo.FindCurrentYear(tenantID)
	if err != nil {
		return errors.New("no active academic year found")
	}

	classUUID := uuid.MustParse(req.ClassID)
	sectionUUID := uuid.MustParse(req.SectionID)
	var batch []domain.StudentAttendance

	for _, entry := range req.Students {
		batch = append(batch, domain.StudentAttendance{
			TenantID:       tenantID,
			AcademicYearID: year.ID,
			ClassID:        classUUID,
			SectionID:      sectionUUID,
			StudentID:      uuid.MustParse(entry.StudentID),
			Date:           reqDate,
			Status:         domain.AttendanceStatus(entry.Status),
			Remarks:        entry.Remarks,
			MarkedBy:       teacherID,
		})
	}

	if len(batch) == 0 {
		return nil
	}

	return s.Repo.UpsertAttendance(batch)
}

func (s *AttendanceService) GetClassRegister(tenantID uuid.UUID, sectionIDStr string, dateStr string) (*dto.AttendanceRegisterResponse, error) {
	sectionID := uuid.MustParse(sectionIDStr)
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format YYYY-MM-DD")
	}

	// 1. Get Records
	records, err := s.Repo.FindBySectionAndDate(sectionID, date)
	if err != nil {
		return nil, err
	}

	// 2. Get All Active Students (To show who is missing from list or defaulting)
	// (Optional: In a real app, you might want to fetch ALL students and merge with records
	// so the UI shows 'Not Marked' instead of empty rows. For now, we return what is DB.)

	var details []dto.RecordDetail
	stats := dto.AttendanceStat{Total: len(records)}

	for _, r := range records {
		details = append(details, dto.RecordDetail{
			StudentID:   r.StudentID,
			StudentName: r.Student.FirstName + " " + r.Student.LastName,
			RollNo:      r.Student.AdmissionNo,
			Status:      string(r.Status),
			Remarks:     r.Remarks,
		})

		switch r.Status {
		case domain.AttendancePresent:
			stats.Present++
		case domain.AttendanceAbsent:
			stats.Absent++
		case domain.AttendanceLate:
			stats.Late++
		}
	}

	return &dto.AttendanceRegisterResponse{
		Date:      dateStr,
		SectionID: sectionID,
		Records:   details,
		Summary:   stats,
	}, nil
}

func (s *AttendanceService) GetStudentHistory(tenantID, studentID uuid.UUID, startStr, endStr string) ([]dto.StudentAttendanceHistory, error) {
	start := time.Now().AddDate(0, -1, 0)
	end := time.Now()

	if startStr != "" {
		start, _ = time.Parse("2006-01-02", startStr)
	}
	if endStr != "" {
		end, _ = time.Parse("2006-01-02", endStr)
	}

	records, err := s.Repo.FindByStudent(studentID, start, end)
	if err != nil {
		return nil, err
	}

	var history []dto.StudentAttendanceHistory
	for _, r := range records {
		history = append(history, dto.StudentAttendanceHistory{
			Date:    r.Date.Format("2006-01-02"),
			Status:  string(r.Status),
			Remarks: r.Remarks,
		})
	}
	return history, nil
}
