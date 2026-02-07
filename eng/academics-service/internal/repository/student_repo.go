package repository

import (
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentRepository struct {
	DB *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

func (r *StudentRepository) CreateSchema(schema *domain.AdmissionSchema) error {
	return r.DB.Create(schema).Error
}

func (r *StudentRepository) UpdateSchema(schema *domain.AdmissionSchema) error {
	return r.DB.Save(schema).Error
}

func (r *StudentRepository) DeactivateOldSchemas(tenantID uuid.UUID) error {
	return r.DB.Model(&domain.AdmissionSchema{}).
		Where("tenant_id = ?", tenantID).
		Update("is_active", false).Error
}

func (r *StudentRepository) FindActiveSchema(tenantID uuid.UUID) (*domain.AdmissionSchema, error) {
	var schema domain.AdmissionSchema
	err := r.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true).First(&schema).Error
	return &schema, err
}

func (r *StudentRepository) DeleteSchema(tenantID uuid.UUID) error {
	return r.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true).Delete(&domain.AdmissionSchema{}).Error
}

func (r *StudentRepository) CreateStudent(student *domain.Student) error {
	return r.DB.Create(student).Error
}

func (r *StudentRepository) FindStudentByID(tenantID, id uuid.UUID) (*domain.Student, error) {
	var student domain.Student
	err := r.DB.Preload("Class").Preload("Section").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&student).Error
	return &student, err
}

func (r *StudentRepository) UpdateStudent(student *domain.Student) error {
	return r.DB.Save(student).Error
}

func (r *StudentRepository) DeleteStudent(tenantID, id uuid.UUID) error {
	return r.DB.Where("tenant_id = ?", tenantID).Delete(&domain.Student{}, "id = ?", id).Error
}

func (r *StudentRepository) FindAllStudents(tenantID uuid.UUID, classID, status string) ([]domain.Student, error) {
	var students []domain.Student
	query := r.DB.Preload("Class").Preload("Section").Where("tenant_id = ?", tenantID)

	if classID != "" {
		query = query.Where("class_id = ?", classID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at desc").Find(&students).Error
	return students, err
}

// FindActiveStudents retrieves all approved students with a non-empty admission number.
func (r *StudentRepository) FindActiveStudents(tenantID uuid.UUID, classID string) ([]domain.Student, error) {
	var students []domain.Student
	query := r.DB.Preload("Class").Preload("Section").
		Where("tenant_id = ?", tenantID).
		Where("status = ?", domain.StatusApproved).
		Where("admission_no != ''")

	if classID != "" {
		query = query.Where("class_id = ?", classID)
	}

	err := query.Order("admission_no asc").Find(&students).Error
	return students, err
}

func (r *StudentRepository) UpsertSequence(seq *domain.AdmissionSequence) error {
	return r.DB.Save(seq).Error
}

func (r *StudentRepository) GetSequenceByYear(tenantID, yearID uuid.UUID) (*domain.AdmissionSequence, error) {
	var seq domain.AdmissionSequence
	err := r.DB.Where("tenant_id = ? AND academic_year_id = ?", tenantID, yearID).First(&seq).Error
	return &seq, err
}

func (r *StudentRepository) FindActiveStudentsByClass(classID uuid.UUID) ([]domain.Student, error) {
	var students []domain.Student
	err := r.DB.Where("class_id = ? AND status = ?", classID, domain.StatusApproved).Find(&students).Error
	return students, err
}
