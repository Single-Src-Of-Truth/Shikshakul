package repository

import (
	"fmt"
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CertificateRepository struct {
	DB *gorm.DB
}

func NewCertificateRepository(db *gorm.DB) *CertificateRepository {
	return &CertificateRepository{DB: db}
}

func (r *CertificateRepository) CreateRecord(record *domain.CertificateRecord) error {
	return r.DB.Create(record).Error
}

func (r *CertificateRepository) FindRecordsByStudent(studentID uuid.UUID) ([]domain.CertificateRecord, error) {
	var records []domain.CertificateRecord
	err := r.DB.Where("student_id = ?", studentID).Order("issue_date desc").Find(&records).Error
	return records, err
}

func (r *CertificateRepository) GenerateCertificateNo(cType domain.CertificateType) string {
	return fmt.Sprintf("%s-%d-%d", string(cType)[:2], time.Now().Year(), time.Now().Unix()%100000)
}

func (r *CertificateRepository) GetStudentFullDetails(studentID uuid.UUID) (*domain.Student, error) {
	var s domain.Student
	err := r.DB.Preload("Class").Preload("Section").
		First(&s, "id = ?", studentID).Error
	return &s, err
}

func (r *CertificateRepository) UpdateStudentStatus(studentID uuid.UUID, status domain.Status) error {
	return r.DB.Model(&domain.Student{}).Where("id = ?", studentID).Update("status", status).Error
}
