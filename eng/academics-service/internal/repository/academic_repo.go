package repository

import (
	"errors"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AcademicRepository struct {
	DB *gorm.DB
}

func NewAcademicRepository(db *gorm.DB) *AcademicRepository {
	return &AcademicRepository{DB: db}
}

func (r *AcademicRepository) CreateYear(year *domain.AcademicYear) error {
	return r.DB.Create(year).Error
}

func (r *AcademicRepository) FindYearByID(id uuid.UUID) (*domain.AcademicYear, error) {
	var year domain.AcademicYear
	err := r.DB.First(&year, "id = ?", id).Error
	return &year, err
}

func (r *AcademicRepository) UpdateYear(year *domain.AcademicYear) error {
	return r.DB.Save(year).Error
}

func (r *AcademicRepository) DeleteYear(id uuid.UUID) error {
	return r.DB.Delete(&domain.AcademicYear{}, "id = ?", id).Error
}

func (r *AcademicRepository) FindCurrentYear(tenantID uuid.UUID) (*domain.AcademicYear, error) {
	var year domain.AcademicYear
	err := r.DB.Where("tenant_id = ? AND is_current = ?", tenantID, true).First(&year).Error
	return &year, err
}

func (r *AcademicRepository) UnsetOtherCurrentYears(tenantID uuid.UUID) error {
	return r.DB.Model(&domain.AcademicYear{}).
		Where("tenant_id = ?", tenantID).
		Update("is_current", false).Error
}

func (r *AcademicRepository) FindAllYears(tenantID uuid.UUID) ([]domain.AcademicYear, error) {
	var years []domain.AcademicYear
	err := r.DB.Where("tenant_id = ?", tenantID).Order("start_date desc").Find(&years).Error
	return years, err
}

func (r *AcademicRepository) CreateClass(class *domain.Class) error {
	return r.DB.Create(class).Error
}

func (r *AcademicRepository) FindClassByID(classID uuid.UUID) (*domain.Class, error) {
	var class domain.Class
	err := r.DB.First(&class, "id = ?", classID).Error
	return &class, err
}

func (r *AcademicRepository) FindAllClasses(tenantID uuid.UUID) ([]domain.Class, error) {
	var classes []domain.Class
	err := r.DB.Preload("Sections").Where("tenant_id = ?", tenantID).Order("sort_order asc").Find(&classes).Error
	return classes, err
}

func (r *AcademicRepository) UpdateClass(class *domain.Class) error {
	return r.DB.Save(class).Error
}

func (r *AcademicRepository) DeleteClass(id uuid.UUID) error {
	return r.DB.Delete(&domain.Class{}, "id = ?", id).Error
}

func (r *AcademicRepository) CountSectionsByClass(classID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.Section{}).Where("class_id = ?", classID).Count(&count).Error
	return count, err
}

func (r *AcademicRepository) CreateSection(section *domain.Section) error {
	return r.DB.Create(section).Error
}

func (r *AcademicRepository) FindSectionByID(id uuid.UUID) (*domain.Section, error) {
	var section domain.Section
	err := r.DB.First(&section, "id = ?", id).Error
	return &section, err
}

func (r *AcademicRepository) FindSectionsByClass(classID uuid.UUID) ([]domain.Section, error) {
	var sections []domain.Section
	err := r.DB.Where("class_id = ?", classID).Find(&sections).Error
	return sections, err
}

func (r *AcademicRepository) UpdateSection(section *domain.Section) error {
	return r.DB.Save(section).Error
}

func (r *AcademicRepository) DeleteSection(id uuid.UUID) error {
	return r.DB.Delete(&domain.Section{}, "id = ?", id).Error
}

func (r *AcademicRepository) CountStudentsBySection(sectionID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.Student{}).Where("section_id = ?", sectionID).Count(&count).Error
	return count, err
}

func (r *AcademicRepository) CreateSubject(subject *domain.Subject) error {
	return r.DB.Create(subject).Error
}

func (r *AcademicRepository) FindSubjectByID(id uuid.UUID) (*domain.Subject, error) {
	var subject domain.Subject
	err := r.DB.First(&subject, "id = ?", id).Error
	return &subject, err
}

func (r *AcademicRepository) FindAllSubjects(tenantID uuid.UUID) ([]domain.Subject, error) {
	var subjects []domain.Subject
	err := r.DB.Where("tenant_id = ?", tenantID).Find(&subjects).Error
	return subjects, err
}

func (r *AcademicRepository) UpdateSubject(subject *domain.Subject) error {
	return r.DB.Save(subject).Error
}

func (r *AcademicRepository) DeleteSubject(id uuid.UUID) error {
	return r.DB.Delete(&domain.Subject{}, "id = ?", id).Error
}

func (r *AcademicRepository) CountClassSubjectMappings(subjectID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.ClassSubject{}).Where("subject_id = ?", subjectID).Count(&count).Error
	return count, err
}

func (r *AcademicRepository) RemoveSubjectFromClass(classID, subjectID uuid.UUID) error {
	result := r.DB.Where("class_id = ? AND subject_id = ?", classID, subjectID).Delete(&domain.ClassSubject{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("mapping not found")
	}
	return nil
}
