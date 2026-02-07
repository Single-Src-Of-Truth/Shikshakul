package repository

import (
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamRepository struct {
	DB *gorm.DB
}

func NewExamRepository(db *gorm.DB) *ExamRepository {
	return &ExamRepository{DB: db}
}

func (r *ExamRepository) CreateTerm(term *domain.ExamTerm) error {
	return r.DB.Create(term).Error
}

func (r *ExamRepository) FindTermByID(id uuid.UUID) (*domain.ExamTerm, error) {
	var term domain.ExamTerm
	err := r.DB.First(&term, "id = ?", id).Error
	return &term, err
}

func (r *ExamRepository) FindTermsByYear(yearID uuid.UUID) ([]domain.ExamTerm, error) {
	var terms []domain.ExamTerm
	err := r.DB.Where("academic_year_id = ?", yearID).Find(&terms).Error
	return terms, err
}

func (r *ExamRepository) DeleteTerm(id uuid.UUID) error {
	return r.DB.Delete(&domain.ExamTerm{}, "id = ?", id).Error
}

func (r *ExamRepository) CreateSchedule(schedule *domain.ExamSchedule) error {
	return r.DB.Create(schedule).Error
}

func (r *ExamRepository) FindSchedulesByTermAndClass(termID, classID uuid.UUID) ([]domain.ExamSchedule, error) {
	var schedules []domain.ExamSchedule
	err := r.DB.Preload("Subject").
		Where("exam_term_id = ? AND class_id = ?", termID, classID).
		Order("exam_date, start_time").
		Find(&schedules).Error
	return schedules, err
}

func (r *ExamRepository) DeleteSchedule(id uuid.UUID) error {
	return r.DB.Delete(&domain.ExamSchedule{}, "id = ?", id).Error
}

// CheckConflicts: Ensures class doesn't have another exam overlapping
// Simplified: Check if Class has ANY exam on that Date
func (r *ExamRepository) CheckScheduleConflict(classID uuid.UUID, date time.Time) (bool, error) {
	var count int64
	// Truncate time part for strict date checking if needed, but standard Time comparison works if stored correctly.
	// Assuming date passed is 00:00:00
	err := r.DB.Model(&domain.ExamSchedule{}).
		Where("class_id = ? AND exam_date = ?", classID, date).
		Count(&count).Error
	return count > 0, err
}
