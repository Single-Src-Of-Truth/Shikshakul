package repository

import (
	"errors"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimetableRepository struct {
	DB *gorm.DB
}

func NewTimetableRepository(db *gorm.DB) *TimetableRepository {
	return &TimetableRepository{DB: db}
}

func (r *TimetableRepository) CreateRoutine(routine *domain.WeeklyRoutine) error {
	return r.DB.Create(routine).Error
}

func (r *TimetableRepository) DeleteRoutine(id uuid.UUID) error {
	return r.DB.Delete(&domain.WeeklyRoutine{}, "id = ?", id).Error
}

func (r *TimetableRepository) FindRoutines(tenantID uuid.UUID, sectionID, teacherID uuid.UUID, day string) ([]domain.WeeklyRoutine, error) {
	var routines []domain.WeeklyRoutine
	query := r.DB.Preload("Subject").Preload("Teacher").
		Where("tenant_id = ?", tenantID)

	if sectionID != uuid.Nil {
		query = query.Where("section_id = ?", sectionID)
	}
	if teacherID != uuid.Nil {
		query = query.Where("teacher_id = ?", teacherID)
	}
	if day != "" {
		query = query.Where("day_of_week = ?", day)
	}

	err := query.Order("day_of_week, start_time").Find(&routines).Error
	return routines, err
}

func (r *TimetableRepository) CheckConflicts(tenantID, teacherID, sectionID uuid.UUID, day, start, end string) error {
	var count int64

	err := r.DB.Model(&domain.WeeklyRoutine{}).
		Where("tenant_id = ? AND day_of_week = ?", tenantID, day).
		Where("start_time < ? AND end_time > ?", end, start).
		Where(r.DB.Where("teacher_id = ?", teacherID).Or("section_id = ?", sectionID)).
		Count(&count).Error

	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("conflict detected: teacher or section is already occupied at this time")
	}
	return nil
}
