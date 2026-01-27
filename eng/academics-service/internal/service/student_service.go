package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type StudentService struct {
	Repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{Repo: repo}
}

func (s *StudentService) CreateSchema(tenantID uuid.UUID, structure []map[string]interface{}) (*domain.AdmissionSchema, error) {
	structureBytes, err := json.Marshal(structure)
	if err != nil {
		return nil, errors.New("invalid schema JSON")
	}

	_ = s.Repo.DeactivateOldSchemas(tenantID)

	schema := &domain.AdmissionSchema{
		TenantID:  tenantID,
		Structure: structureBytes,
		Version:   1,
		IsActive:  true,
	}
	err = s.Repo.CreateSchema(schema)
	return schema, err
}

func (s *StudentService) GetActiveSchema(tenantID uuid.UUID) (*domain.AdmissionSchema, error) {
	return s.Repo.FindActiveSchema(tenantID)
}

func (s *StudentService) DeleteSchema(tenantID uuid.UUID) error {
	return s.Repo.DeleteSchema(tenantID)
}

func (s *StudentService) SubmitApplication(tenantID uuid.UUID, req dto.SubmitApplicationRequest) (*domain.Student, error) {
	schema, err := s.Repo.FindActiveSchema(tenantID)
	if err != nil {
		return nil, errors.New("admission schema not configured")
	}

	if err := s.validateDynamicData(schema.Structure, req.ProfileData); err != nil {
		return nil, err
	}

	profileBytes, _ := json.Marshal(req.ProfileData)
	docBytes, _ := json.Marshal(req.Documents)
	academicYearUUID, _ := uuid.Parse(req.AcademicYearID)
	classUUID, _ := uuid.Parse(req.ClassID)

	student := &domain.Student{
		TenantID:       tenantID,
		AcademicYearID: academicYearUUID,
		ClassID:        classUUID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Mobile:         req.Mobile,
		ProfileData:    profileBytes,
		Documents:      docBytes,
		Status:         domain.StatusPending,
	}

	if err := s.Repo.CreateStudent(student); err != nil {
		return nil, err
	}
	return s.Repo.FindStudentByID(tenantID, student.ID)
}

func (s *StudentService) UpdateStudent(tenantID, studentID uuid.UUID, req dto.UpdateStudentRequest) (*domain.Student, error) {
	student, err := s.Repo.FindStudentByID(tenantID, studentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	if req.FirstName != "" {
		student.FirstName = req.FirstName
	}
	if req.LastName != "" {
		student.LastName = req.LastName
	}
	if req.Email != "" {
		student.Email = req.Email
	}
	if req.Mobile != "" {
		student.Mobile = req.Mobile
	}

	if req.ProfileData != nil {
		bytes, _ := json.Marshal(req.ProfileData)
		student.ProfileData = bytes
	}
	if req.Documents != nil {
		bytes, _ := json.Marshal(req.Documents)
		student.Documents = bytes
	}

	err = s.Repo.UpdateStudent(student)
	return student, err
}

func (s *StudentService) DeleteStudent(tenantID, studentID uuid.UUID) error {
	_, err := s.Repo.FindStudentByID(tenantID, studentID)
	if err != nil {
		return errors.New("student not found")
	}

	return s.Repo.DeleteStudent(tenantID, studentID)
}

func (s *StudentService) ProcessApproval(tenantID uuid.UUID, studentID uuid.UUID, req dto.ApproveStudentRequest) (*domain.Student, error) {
	tx := s.Repo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var student domain.Student
	if err := tx.Preload("Class").Where("id = ? AND tenant_id = ?", studentID, tenantID).First(&student).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("student not found")
	}

	if student.Status == domain.StatusApproved {
		tx.Rollback()
		return nil, errors.New("operation failed: student is already approved (use withdrawal process to remove)")
	}

	if req.Action == "PENDING" {
		student.Status = domain.StatusPending
		if err := tx.Save(&student).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		return &student, nil
	}

	if req.Action == "REJECT" {
		student.Status = domain.StatusRejected
		if err := tx.Save(&student).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		return &student, nil
	}

	if req.SectionID == "" {
		tx.Rollback()
		return nil, errors.New("section_id is required")
	}
	sectionUUID := uuid.MustParse(req.SectionID)

	var year domain.AcademicYear
	if err := tx.First(&year, "id = ?", student.AcademicYearID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("academic year not found")
	}

	var seq domain.AdmissionSequence
	if err := tx.Where("tenant_id = ? AND academic_year_id = ?", tenantID, student.AcademicYearID).First(&seq).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("admission sequence not configured for this year")
	}

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&seq, "id = ?", seq.ID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	seq.CurrentCount++
	tx.Save(&seq)

	yearShort := year.StartDate.Year() % 100
	admNo := fmt.Sprintf("%s-%d-%04d", seq.Prefix, yearShort, seq.CurrentCount)

	student.Status = domain.StatusApproved
	student.AdmissionNo = admNo
	student.SectionID = &sectionUUID

	if err := tx.Save(&student).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var finalStudent domain.Student
	if err := tx.Preload("Class").Preload("Section").
		Where("id = ? AND tenant_id = ?", student.ID, tenantID).
		First(&finalStudent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &finalStudent, nil
}

func (s *StudentService) GetAllStudents(tenantID uuid.UUID, filter dto.StudentFilter) ([]domain.Student, error) {
	return s.Repo.FindAllStudents(tenantID, filter.ClassID, filter.Status)
}

func (s *StudentService) GetActiveStudents(tenantID uuid.UUID, classID string) ([]domain.Student, error) {
	return s.Repo.FindActiveStudents(tenantID, classID)
}

func (s *StudentService) GetStudentByID(tenantID, studentID uuid.UUID) (*domain.Student, error) {
	return s.Repo.FindStudentByID(tenantID, studentID)
}

func (s *StudentService) ConfigureAdmissionSequence(tenantID uuid.UUID, req dto.ConfigureSequenceRequest) error {
	yearID := uuid.MustParse(req.AcademicYearID)
	seq, err := s.Repo.GetSequenceByYear(tenantID, yearID)
	if err != nil {
		startingCount := 0
		if req.StartingCount != nil {
			startingCount = *req.StartingCount
		}

		seq = &domain.AdmissionSequence{
			TenantID:       tenantID,
			AcademicYearID: yearID,
			Prefix:         req.Prefix,
			CurrentCount:   startingCount,
		}
		return s.Repo.UpsertSequence(seq)
	}

	if seq.CurrentCount > 0 && seq.Prefix != req.Prefix {
		return errors.New("cannot change prefix after admission numbers have been generated for this year")
	}

	seq.Prefix = req.Prefix

	if req.StartingCount != nil && *req.StartingCount > seq.CurrentCount {
		seq.CurrentCount = *req.StartingCount
	}

	return s.Repo.UpsertSequence(seq)
}

func (s *StudentService) validateDynamicData(schemaJSON []byte, inputData map[string]interface{}) error {
	var fields []map[string]interface{}
	if err := json.Unmarshal(schemaJSON, &fields); err != nil {
		return errors.New("corrupt schema configuration")
	}
	for _, field := range fields {
		key := field["key"].(string)
		isRequired, _ := field["required"].(bool)
		val, exists := inputData[key]
		if isRequired && (!exists || val == "" || val == nil) {
			return fmt.Errorf("missing required field: %s", field["label"])
		}
	}
	return nil
}
