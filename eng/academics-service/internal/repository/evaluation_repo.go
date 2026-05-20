package repository

import (
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EvaluationRepository struct {
	DB *gorm.DB
}

func NewEvaluationRepository(db *gorm.DB) *EvaluationRepository {
	return &EvaluationRepository{DB: db}
}

func (r *EvaluationRepository) UpsertMarks(marks []domain.StudentMark) error {
	return r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "exam_schedule_id"}, {Name: "student_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"marks_obtained", "is_absent", "remarks", "graded_by", "status", "updated_at"}),
	}).Create(&marks).Error
}

func (r *EvaluationRepository) FindMarksBySchedule(scheduleID uuid.UUID) ([]domain.StudentMark, error) {
	var marks []domain.StudentMark
	err := r.DB.Where("exam_schedule_id = ?", scheduleID).Find(&marks).Error
	return marks, err
}

func (r *EvaluationRepository) FindStudentsByExamSchedule(scheduleID uuid.UUID) ([]domain.Student, error) {
	var students []domain.Student

	var schedule domain.ExamSchedule
	if err := r.DB.Select("class_id").First(&schedule, "id = ?", scheduleID).Error; err != nil {
		return nil, err
	}

	err := r.DB.Where("class_id = ? AND status = ?", schedule.ClassID, domain.StatusApproved).
		Order("first_name asc").
		Find(&students).Error

	return students, err
}

func (r *EvaluationRepository) FindScheduleByID(id uuid.UUID) (*domain.ExamSchedule, error) {
	var s domain.ExamSchedule
	err := r.DB.Preload("Subject").First(&s, "id = ?", id).Error
	return &s, err
}
