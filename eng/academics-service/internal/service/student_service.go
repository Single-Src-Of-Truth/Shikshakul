package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
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

	err = s.Repo.CreateStudent(student)
	return student, err
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

func (s *StudentService) ProcessApproval(tenantID, studentID uuid.UUID, req dto.ApproveStudentRequest) (*domain.Student, error) {
	student, err := s.Repo.FindStudentByID(tenantID, studentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	if req.Action == "REJECT" {
		student.Status = domain.StatusRejected
		err = s.Repo.UpdateStudent(student)
		return student, err
	}

	if req.SectionID == "" {
		return nil, errors.New("section_id required")
	}

	sectionUUID, err := uuid.Parse(req.SectionID)
	if err != nil {
		return nil, errors.New("invalid section id")
	}

	year := time.Now().Year()
	randNum := rand.Intn(9000) + 1000
	admNo := fmt.Sprintf("ADM-%d-%d", year, randNum)

	student.Status = domain.StatusApproved
	student.AdmissionNo = admNo
	student.SectionID = &sectionUUID

	err = s.Repo.UpdateStudent(student)
	return student, err
}

func (s *StudentService) GetAllStudents(tenantID uuid.UUID, filter dto.StudentFilter) ([]domain.Student, error) {
	return s.Repo.FindAllStudents(tenantID, filter.ClassID, filter.Status)
}

func (s *StudentService) GetStudentByID(tenantID, studentID uuid.UUID) (*domain.Student, error) {
	return s.Repo.FindStudentByID(tenantID, studentID)
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
