package repository

import (
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AttendanceRepository struct {
	DB *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{DB: db}
}

func (r *AttendanceRepository) UpsertAttendance(records []domain.StudentAttendance) error {
	return r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "student_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "remarks", "marked_by", "updated_at"}),
	}).Create(&records).Error
}

func (r *AttendanceRepository) FindBySectionAndDate(sectionID uuid.UUID, date time.Time) ([]domain.StudentAttendance, error) {
	var records []domain.StudentAttendance
	dateOnly := date.Format("2006-01-02")

	err := r.DB.Preload("Student").
		Where("section_id = ? AND date = ?", sectionID, dateOnly).
		Order("created_at asc").
		Find(&records).Error
	return records, err
}

func (r *AttendanceRepository) FindByStudent(studentID uuid.UUID, startDate, endDate time.Time) ([]domain.StudentAttendance, error) {
	var records []domain.StudentAttendance
	err := r.DB.Where("student_id = ? AND date BETWEEN ? AND ?", studentID, startDate, endDate).
		Order("date desc").
		Find(&records).Error
	return records, err
}
