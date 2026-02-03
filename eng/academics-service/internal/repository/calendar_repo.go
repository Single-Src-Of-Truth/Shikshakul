package repository

import (
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CalendarRepository struct {
	DB *gorm.DB
}

func NewCalendarRepository(db *gorm.DB) *CalendarRepository {
	return &CalendarRepository{DB: db}
}

func (r *CalendarRepository) CreateEvent(event *domain.CalendarEvent) error {
	return r.DB.Create(event).Error
}

func (r *CalendarRepository) UpdateEvent(event *domain.CalendarEvent) error {
	return r.DB.Save(event).Error
}

func (r *CalendarRepository) FindEventsByRange(tenantID uuid.UUID, start, end time.Time) ([]domain.CalendarEvent, error) {
	var events []domain.CalendarEvent
	err := r.DB.Where("tenant_id = ? AND start_date >= ? AND start_date <= ?", tenantID, start, end).
		Order("start_date asc").
		Find(&events).Error
	return events, err
}

func (r *CalendarRepository) DeleteEvent(id uuid.UUID) error {
	return r.DB.Delete(&domain.CalendarEvent{}, "id = ?", id).Error
}

func (r *CalendarRepository) FindEventByID(id uuid.UUID) (*domain.CalendarEvent, error) {
	var event domain.CalendarEvent
	err := r.DB.First(&event, "id = ?", id).Error
	return &event, err
}

func (r *CalendarRepository) CheckEventOverlap(tenantID uuid.UUID, start, end time.Time, excludeEventID uuid.UUID) (bool, error) {
	var count int64
	query := r.DB.Model(&domain.CalendarEvent{}).
		Where("tenant_id = ?", tenantID).
		Where("start_date < ? AND end_date > ?", end, start)

	if excludeEventID != uuid.Nil {
		query = query.Where("id != ?", excludeEventID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
