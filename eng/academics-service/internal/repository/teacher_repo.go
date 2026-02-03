package repository

import (
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TeacherRepository struct {
	DB *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) *TeacherRepository {
	return &TeacherRepository{DB: db}
}

func (r *TeacherRepository) CreateSchema(schema *domain.TeacherSchema) error {
	return r.DB.Create(schema).Error
}

func (r *TeacherRepository) FindActiveSchema(tenantID uuid.UUID) (*domain.TeacherSchema, error) {
	var schema domain.TeacherSchema
	err := r.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true).First(&schema).Error
	return &schema, err
}

func (r *TeacherRepository) DeactivateOldSchemas(tenantID uuid.UUID) error {
	return r.DB.Model(&domain.TeacherSchema{}).
		Where("tenant_id = ?", tenantID).
		Update("is_active", false).Error
}

func (r *TeacherRepository) DeleteSchema(tenantID uuid.UUID) error {
	return r.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true).Delete(&domain.TeacherSchema{}).Error
}

func (r *TeacherRepository) CreateTeacher(teacher *domain.Teacher) error {
	return r.DB.Create(teacher).Error
}

func (r *TeacherRepository) FindTeacherByID(tenantID, id uuid.UUID) (*domain.Teacher, error) {
	var teacher domain.Teacher
	err := r.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&teacher).Error
	return &teacher, err
}

func (r *TeacherRepository) FindAllTeachers(tenantID uuid.UUID, status string) ([]domain.Teacher, error) {
	var teachers []domain.Teacher
	query := r.DB.Where("tenant_id = ?", tenantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&teachers).Error
	return teachers, err
}

func (r *TeacherRepository) UpdateTeacher(teacher *domain.Teacher) error {
	return r.DB.Save(teacher).Error
}

func (r *TeacherRepository) DeleteTeacher(tenantID, id uuid.UUID) error {
	return r.DB.Where("tenant_id = ?", tenantID).Delete(&domain.Teacher{}, "id = ?", id).Error
}

func (r *TeacherRepository) IsClassTeacher(teacherID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.Section{}).Where("class_teacher_id = ?", teacherID).Count(&count).Error
	return count > 0, err
}

func (r *TeacherRepository) HasSubjectAllocations(teacherID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.SubjectAllocation{}).Where("teacher_id = ?", teacherID).Count(&count).Error
	return count > 0, err
}

func (r *TeacherRepository) AssignClassTeacher(sectionID uuid.UUID, teacherID *uuid.UUID) error {
	return r.DB.Model(&domain.Section{}).
		Where("id = ?", sectionID).
		Update("class_teacher_id", teacherID).Error
}

func (r *TeacherRepository) UpsertSubjectAllocation(allocation *domain.SubjectAllocation) error {
	return r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "section_id"}, {Name: "subject_id"}, {Name: "academic_year_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"teacher_id", "updated_at"}),
	}).Create(allocation).Error
}

func (r *TeacherRepository) RemoveSubjectAllocation(sectionID, subjectID, yearID uuid.UUID) error {
	return r.DB.Where("section_id = ? AND subject_id = ? AND academic_year_id = ?", sectionID, subjectID, yearID).
		Delete(&domain.SubjectAllocation{}).Error
}

func (r *TeacherRepository) FindAllocationsBySection(sectionID uuid.UUID, yearID uuid.UUID) ([]domain.SubjectAllocation, error) {
	var allocs []domain.SubjectAllocation
	err := r.DB.Preload("Teacher").Preload("Subject").
		Where("section_id = ? AND academic_year_id = ?", sectionID, yearID).
		Find(&allocs).Error
	return allocs, err
}
