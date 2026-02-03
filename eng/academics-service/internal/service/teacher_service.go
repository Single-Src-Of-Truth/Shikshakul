package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
)

type TeacherService struct {
	Repo         *repository.TeacherRepository
	AcademicRepo *repository.AcademicRepository
}

func NewTeacherService(repo *repository.TeacherRepository, acadRepo *repository.AcademicRepository) *TeacherService {
	return &TeacherService{Repo: repo, AcademicRepo: acadRepo}
}

func (s *TeacherService) CreateSchema(tenantID uuid.UUID, structure []map[string]interface{}) (*domain.TeacherSchema, error) {
	structureBytes, err := json.Marshal(structure)
	if err != nil {
		return nil, errors.New("invalid schema JSON")
	}

	_ = s.Repo.DeactivateOldSchemas(tenantID)

	schema := &domain.TeacherSchema{
		TenantID:  tenantID,
		Structure: structureBytes,
		Version:   1,
		IsActive:  true,
	}
	err = s.Repo.CreateSchema(schema)
	return schema, err
}

func (s *TeacherService) GetActiveSchema(tenantID uuid.UUID) (*domain.TeacherSchema, error) {
	return s.Repo.FindActiveSchema(tenantID)
}

func (s *TeacherService) DeleteSchema(tenantID uuid.UUID) error {
	return s.Repo.DeleteSchema(tenantID)
}

func (s *TeacherService) OnboardTeacher(tenantID uuid.UUID, req dto.OnboardTeacherRequest) (*domain.Teacher, error) {
	schema, err := s.Repo.FindActiveSchema(tenantID)
	if err != nil {
		return nil, errors.New("teacher schema not configured")
	}

	if err := s.validateDynamicData(schema.Structure, req.ProfileData); err != nil {
		return nil, err
	}

	profileBytes, _ := json.Marshal(req.ProfileData)
	docBytes, _ := json.Marshal(req.Documents)

	teacher := &domain.Teacher{
		TenantID:    tenantID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Mobile:      req.Mobile,
		ProfileData: profileBytes,
		Documents:   docBytes,
		Status:      domain.StatusPending,
	}

	err = s.Repo.CreateTeacher(teacher)
	return teacher, err
}

func (s *TeacherService) UpdateTeacher(tenantID, id uuid.UUID, req dto.UpdateTeacherRequest) (*domain.Teacher, error) {
	teacher, err := s.Repo.FindTeacherByID(tenantID, id)
	if err != nil {
		return nil, errors.New("teacher not found")
	}

	if req.FirstName != "" {
		teacher.FirstName = req.FirstName
	}
	if req.LastName != "" {
		teacher.LastName = req.LastName
	}
	if req.Email != "" {
		teacher.Email = req.Email
	}
	if req.Mobile != "" {
		teacher.Mobile = req.Mobile
	}

	if req.ProfileData != nil {
		bytes, _ := json.Marshal(req.ProfileData)
		teacher.ProfileData = bytes
	}
	if req.Documents != nil {
		bytes, _ := json.Marshal(req.Documents)
		teacher.Documents = bytes
	}

	err = s.Repo.UpdateTeacher(teacher)
	return teacher, err
}

func (s *TeacherService) DeleteTeacher(tenantID, id uuid.UUID) error {
	_, err := s.Repo.FindTeacherByID(tenantID, id)
	if err != nil {
		return errors.New("teacher not found")
	}

	isClassTeacher, _ := s.Repo.IsClassTeacher(id)
	if isClassTeacher {
		return fmt.Errorf("%w: teacher is assigned as a class teacher", ErrConflict)
	}

	hasAllocations, _ := s.Repo.HasSubjectAllocations(id)
	if hasAllocations {
		return fmt.Errorf("%w: teacher has active subject allocations", ErrConflict)
	}

	return s.Repo.DeleteTeacher(tenantID, id)
}

func (s *TeacherService) GetTeacherByID(tenantID, id uuid.UUID) (*domain.Teacher, error) {
	return s.Repo.FindTeacherByID(tenantID, id)
}

func (s *TeacherService) GetAllTeachers(tenantID uuid.UUID, filter dto.TeacherFilter) ([]domain.Teacher, error) {
	return s.Repo.FindAllTeachers(tenantID, filter.Status)
}

func (s *TeacherService) ApproveTeacher(tenantID, teacherID uuid.UUID, action string) (*domain.Teacher, error) {
	teacher, err := s.Repo.FindTeacherByID(tenantID, teacherID)
	if err != nil {
		return nil, errors.New("teacher not found")
	}

	if action == "REJECT" {
		teacher.Status = domain.StatusRejected
		err = s.Repo.UpdateTeacher(teacher)
		return teacher, err
	}

	empID := fmt.Sprintf("EMP-%d", rand.Intn(90000)+10000)
	teacher.Status = domain.StatusActive
	teacher.EmployeeID = empID

	err = s.Repo.UpdateTeacher(teacher)
	return teacher, err
}

func (s *TeacherService) AssignClassTeacher(tenantID, sectionID, teacherID uuid.UUID) error {
	teacher, err := s.Repo.FindTeacherByID(tenantID, teacherID)
	if err != nil || teacher.Status != domain.StatusActive {
		return errors.New("invalid or inactive teacher")
	}

	err = s.Repo.AssignClassTeacher(sectionID, &teacherID)
	if err == nil {
		teacher.IsClassTeacher = true
		s.Repo.UpdateTeacher(teacher)
	}
	return err
}

func (s *TeacherService) UnassignClassTeacher(sectionID uuid.UUID) error {
	return s.Repo.AssignClassTeacher(sectionID, nil)
}

func (s *TeacherService) AssignSubjectTeacher(tenantID, sectionID uuid.UUID, req dto.AssignSubjectTeacherRequest) error {
	subjectID, _ := uuid.Parse(req.SubjectID)
	yearID, _ := uuid.Parse(req.AcademicYearID)
	teacherID, _ := uuid.Parse(req.TeacherID)

	if _, err := s.Repo.FindTeacherByID(tenantID, teacherID); err != nil {
		return errors.New("teacher not found")
	}

	allocation := &domain.SubjectAllocation{
		TenantID:       tenantID,
		AcademicYearID: yearID,
		SectionID:      sectionID,
		SubjectID:      subjectID,
		TeacherID:      teacherID,
	}

	return s.Repo.UpsertSubjectAllocation(allocation)
}

func (s *TeacherService) UnassignSubjectTeacher(sectionID, subjectID, yearID uuid.UUID) error {
	return s.Repo.RemoveSubjectAllocation(sectionID, subjectID, yearID)
}

func (s *TeacherService) GetAllocations(sectionID, yearID uuid.UUID) ([]domain.SubjectAllocation, error) {
	return s.Repo.FindAllocationsBySection(sectionID, yearID)
}

func (s *TeacherService) validateDynamicData(schemaJSON []byte, inputData map[string]interface{}) error {
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
