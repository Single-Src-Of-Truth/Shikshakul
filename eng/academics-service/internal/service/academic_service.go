package service

import (
	"errors"
	"fmt"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
)

var (
	ErrConflict       = errors.New("resource has dependent records")
	ErrNotFound       = errors.New("resource not found")
	ErrDateValidation = errors.New("end date cannot be before start date")
)

type AcademicService struct {
	Repo *repository.AcademicRepository
}

func NewAcademicService(repo *repository.AcademicRepository) *AcademicService {
	return &AcademicService{Repo: repo}
}

func (s *AcademicService) CreateAcademicYear(tenantID uuid.UUID, req dto.CreateYearRequest) (*dto.AcademicYearResponse, error) {
	if req.EndDate.Before(req.StartDate) {
		return nil, ErrDateValidation
	}

	if req.IsCurrent {
		_ = s.Repo.UnsetOtherCurrentYears(tenantID)
	}

	year := &domain.AcademicYear{
		TenantID:  tenantID,
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsCurrent: req.IsCurrent,
	}

	if err := s.Repo.CreateYear(year); err != nil {
		return nil, err
	}
	return mapYearToDTO(year), nil
}

func (s *AcademicService) UpdateAcademicYear(tenantID, id uuid.UUID, req dto.UpdateYearRequest) (*dto.AcademicYearResponse, error) {
	year, err := s.Repo.FindYearByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	if year.TenantID != tenantID {
		return nil, ErrNotFound
	}

	if req.Name != "" {
		year.Name = req.Name
	}
	if !req.StartDate.IsZero() {
		year.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		year.EndDate = req.EndDate
	}
	if req.IsCurrent != nil {
		year.IsCurrent = *req.IsCurrent
		if *req.IsCurrent {
			_ = s.Repo.UnsetOtherCurrentYears(tenantID)
		}
	}

	if err := s.Repo.UpdateYear(year); err != nil {
		return nil, err
	}
	return mapYearToDTO(year), nil
}

func (s *AcademicService) DeleteAcademicYear(tenantID, id uuid.UUID) error {
	year, err := s.Repo.FindYearByID(id)
	if err != nil {
		return ErrNotFound
	}
	if year.TenantID != tenantID {
		return ErrNotFound
	}

	return s.Repo.DeleteYear(id)
}

func (s *AcademicService) GetAllYears(tenantID uuid.UUID) ([]dto.AcademicYearResponse, error) {
	years, err := s.Repo.FindAllYears(tenantID)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.AcademicYearResponse, len(years))
	for i, y := range years {
		dtos[i] = *mapYearToDTO(&y)
	}
	return dtos, nil
}

func (s *AcademicService) GetCurrentYear(tenantID uuid.UUID) (*dto.AcademicYearResponse, error) {
	year, err := s.Repo.FindCurrentYear(tenantID)
	if err != nil {
		return nil, err
	}
	return mapYearToDTO(year), nil
}

func (s *AcademicService) CreateClass(tenantID uuid.UUID, req dto.CreateClassRequest) (*dto.ClassResponse, error) {
	class := &domain.Class{
		TenantID:  tenantID,
		Name:      req.Name,
		SortOrder: req.SortOrder,
	}
	if err := s.Repo.CreateClass(class); err != nil {
		return nil, err
	}
	return mapClassToDTO(class), nil
}

func (s *AcademicService) UpdateClass(tenantID, id uuid.UUID, req dto.UpdateClassRequest) (*dto.ClassResponse, error) {
	class, err := s.Repo.FindClassByID(id)
	if err != nil || class.TenantID != tenantID {
		return nil, ErrNotFound
	}

	if req.Name != "" {
		class.Name = req.Name
	}
	if req.SortOrder != 0 {
		class.SortOrder = req.SortOrder
	}

	if err := s.Repo.UpdateClass(class); err != nil {
		return nil, err
	}
	return mapClassToDTO(class), nil
}

func (s *AcademicService) DeleteClass(tenantID, id uuid.UUID) error {
	class, err := s.Repo.FindClassByID(id)
	if err != nil || class.TenantID != tenantID {
		return ErrNotFound
	}

	count, _ := s.Repo.CountSectionsByClass(id)
	if count > 0 {
		return fmt.Errorf("%w: class has active sections", ErrConflict)
	}

	return s.Repo.DeleteClass(id)
}

func (s *AcademicService) GetAllClasses(tenantID uuid.UUID) ([]dto.ClassResponse, error) {
	classes, err := s.Repo.FindAllClasses(tenantID)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.ClassResponse, len(classes))
	for i, c := range classes {
		dtos[i] = *mapClassToDTO(&c)
	}
	return dtos, nil
}

func (s *AcademicService) CreateSection(classIDStr string, req dto.CreateSectionRequest) (*dto.SectionResponse, error) {
	classUUID, err := uuid.Parse(classIDStr)
	if err != nil {
		return nil, errors.New("invalid class ID")
	}

	_, err = s.Repo.FindClassByID(classUUID)
	if err != nil {
		return nil, errors.New("class not found")
	}

	section := &domain.Section{
		ClassID:  classUUID,
		Name:     req.Name,
		Capacity: req.Capacity,
	}
	if err := s.Repo.CreateSection(section); err != nil {
		return nil, err
	}
	return mapSectionToDTO(section), nil
}

func (s *AcademicService) UpdateSection(id uuid.UUID, req dto.UpdateSectionRequest) (*dto.SectionResponse, error) {
	section, err := s.Repo.FindSectionByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if req.Name != "" {
		section.Name = req.Name
	}
	if req.Capacity != 0 {
		section.Capacity = req.Capacity
	}

	if err := s.Repo.UpdateSection(section); err != nil {
		return nil, err
	}
	return mapSectionToDTO(section), nil
}

func (s *AcademicService) DeleteSection(id uuid.UUID) error {
	_, err := s.Repo.FindSectionByID(id)
	if err != nil {
		return ErrNotFound
	}

	count, _ := s.Repo.CountStudentsBySection(id)
	if count > 0 {
		return fmt.Errorf("%w: section contains active students", ErrConflict)
	}

	return s.Repo.DeleteSection(id)
}

func (s *AcademicService) GetSectionsByClass(classID uuid.UUID) ([]dto.SectionResponse, error) {
	sections, err := s.Repo.FindSectionsByClass(classID)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.SectionResponse, len(sections))
	for i, sec := range sections {
		dtos[i] = *mapSectionToDTO(&sec)
	}
	return dtos, nil
}

func (s *AcademicService) CreateSubject(tenantID uuid.UUID, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
	subject := &domain.Subject{
		TenantID: tenantID,
		Name:     req.Name,
		Code:     req.Code,
		Type:     domain.SubjectType(req.Type),
	}
	if err := s.Repo.CreateSubject(subject); err != nil {
		return nil, err
	}
	return mapSubjectToDTO(subject), nil
}

func (s *AcademicService) UpdateSubject(tenantID, id uuid.UUID, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error) {
	subject, err := s.Repo.FindSubjectByID(id)
	if err != nil || subject.TenantID != tenantID {
		return nil, ErrNotFound
	}

	if req.Name != "" {
		subject.Name = req.Name
	}
	if req.Code != "" {
		subject.Code = req.Code
	}
	if req.Type != "" {
		subject.Type = domain.SubjectType(req.Type)
	}

	if err := s.Repo.UpdateSubject(subject); err != nil {
		return nil, err
	}
	return mapSubjectToDTO(subject), nil
}

func (s *AcademicService) DeleteSubject(tenantID, id uuid.UUID) error {
	subject, err := s.Repo.FindSubjectByID(id)
	if err != nil || subject.TenantID != tenantID {
		return ErrNotFound
	}

	count, _ := s.Repo.CountClassSubjectMappings(id)
	if count > 0 {
		return fmt.Errorf("%w: subject is assigned to classes", ErrConflict)
	}

	return s.Repo.DeleteSubject(id)
}

func (s *AcademicService) GetAllSubjects(tenantID uuid.UUID) ([]dto.SubjectResponse, error) {
	subjects, err := s.Repo.FindAllSubjects(tenantID)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.SubjectResponse, len(subjects))
	for i, sub := range subjects {
		dtos[i] = *mapSubjectToDTO(&sub)
	}
	return dtos, nil
}

func (s *AcademicService) AssignSubjectsToClass(tenantID uuid.UUID, classIDStr string, req dto.AssignSubjectsRequest) error {
	classUUID, _ := uuid.Parse(classIDStr)

	tx := s.Repo.DB.Begin()
	for _, item := range req.Subjects {
		mapping := domain.ClassSubject{
			TenantID:       tenantID,
			ClassID:        classUUID,
			SubjectID:      uuid.MustParse(item.SubjectID),
			IsOptional:     item.IsOptional,
			WeeklyLectures: item.WeeklyLectures,
		}
		if err := tx.Where("class_id = ? AND subject_id = ?", classUUID, mapping.SubjectID).
			FirstOrCreate(&mapping).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (s *AcademicService) UnassignSubject(tenantID, classID, subjectID uuid.UUID) error {
	return s.Repo.RemoveSubjectFromClass(classID, subjectID)
}

func mapYearToDTO(y *domain.AcademicYear) *dto.AcademicYearResponse {
	return &dto.AcademicYearResponse{ID: y.ID, Name: y.Name, StartDate: y.StartDate, EndDate: y.EndDate, IsCurrent: y.IsCurrent}
}
func mapClassToDTO(c *domain.Class) *dto.ClassResponse {
	return &dto.ClassResponse{ID: c.ID, Name: c.Name, SortOrder: c.SortOrder}
}
func mapSectionToDTO(s *domain.Section) *dto.SectionResponse {
	return &dto.SectionResponse{ID: s.ID, ClassID: s.ClassID, Name: s.Name, Capacity: s.Capacity}
}
func mapSubjectToDTO(s *domain.Subject) *dto.SubjectResponse {
	return &dto.SubjectResponse{ID: s.ID, Name: s.Name, Code: s.Code, Type: string(s.Type)}
}
