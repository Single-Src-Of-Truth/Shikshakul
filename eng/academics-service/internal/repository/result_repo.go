package repository

import (
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ResultRepository struct {
	DB *gorm.DB
}

func NewResultRepository(db *gorm.DB) *ResultRepository {
	return &ResultRepository{DB: db}
}

func (r *ResultRepository) FetchStudentMarksForTerm(studentID, termID uuid.UUID) ([]domain.StudentMark, error) {
	var marks []domain.StudentMark

	err := r.DB.Preload("ExamSchedule").Preload("ExamSchedule.Subject").
		Joins("JOIN exam_schedules ON exam_schedules.id = student_marks.exam_schedule_id").
		Where("student_marks.student_id = ? AND exam_schedules.exam_term_id = ?", studentID, termID).
		Find(&marks).Error

	return marks, err
}

func (r *ResultRepository) SaveResult(result *domain.ExamResult) error {
	return r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "exam_term_id"}, {Name: "student_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"total_marks", "obtained_marks", "percentage", "grade", "result_status", "updated_at"}),
	}).Create(result).Error
}

func (r *ResultRepository) FindResult(studentID, termID uuid.UUID) (*domain.ExamResult, error) {
	var res domain.ExamResult
	err := r.DB.Preload("Student").Preload("Student.Class").Preload("ExamTerm").
		Where("student_id = ? AND exam_term_id = ?", studentID, termID).
		First(&res).Error
	return &res, err
}

func (r *ResultRepository) BulkPublishResults(termID, classID uuid.UUID, publish bool) error {
	return r.DB.Model(&domain.ExamResult{}).
		Where("exam_term_id = ? AND class_id = ?", termID, classID).
		Update("is_published", publish).Error
}
