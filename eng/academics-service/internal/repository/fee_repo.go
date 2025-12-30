package repository

import (
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FeeRepository struct {
	DB *gorm.DB
}

func NewFeeRepository(db *gorm.DB) *FeeRepository {
	return &FeeRepository{DB: db}
}

func (r *FeeRepository) CreateFeeHead(head *domain.FeeHead) error {
	return r.DB.Create(head).Error
}

func (r *FeeRepository) FindHeadByID(id uuid.UUID) (*domain.FeeHead, error) {
	var head domain.FeeHead
	err := r.DB.First(&head, "id = ?", id).Error
	return &head, err
}

func (r *FeeRepository) GetAllFeeHeads(tenantID uuid.UUID) ([]domain.FeeHead, error) {
	var heads []domain.FeeHead
	err := r.DB.Where("tenant_id = ?", tenantID).Find(&heads).Error
	return heads, err
}

func (r *FeeRepository) UpdateFeeHead(head *domain.FeeHead) error {
	return r.DB.Save(head).Error
}

func (r *FeeRepository) DeleteFeeHead(id uuid.UUID) error {
	return r.DB.Delete(&domain.FeeHead{}, "id = ?", id).Error
}

func (r *FeeRepository) CountStructuresByHead(headID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.FeeStructure{}).Where("fee_head_id = ?", headID).Count(&count).Error
	return count, err
}

func (r *FeeRepository) CreateFeeStructure(structure *domain.FeeStructure) error {
	return r.DB.Create(structure).Error
}

func (r *FeeRepository) FindStructureByID(id uuid.UUID) (*domain.FeeStructure, error) {
	var st domain.FeeStructure
	err := r.DB.First(&st, "id = ?", id).Error
	return &st, err
}

func (r *FeeRepository) FindStructuresByClass(classID uuid.UUID, yearID uuid.UUID) ([]domain.FeeStructure, error) {
	var structures []domain.FeeStructure
	err := r.DB.Preload("FeeHead").
		Where("class_id = ? AND academic_year_id = ?", classID, yearID).
		Find(&structures).Error
	return structures, err
}

func (r *FeeRepository) UpdateFeeStructure(st *domain.FeeStructure) error {
	return r.DB.Save(st).Error
}

func (r *FeeRepository) DeleteFeeStructure(id uuid.UUID) error {
	return r.DB.Delete(&domain.FeeStructure{}, "id = ?", id).Error
}

func (r *FeeRepository) CountFeesByStructure(structureID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.StudentFee{}).Where("fee_structure_id = ?", structureID).Count(&count).Error
	return count, err
}

func (r *FeeRepository) BatchCreateStudentFees(fees []domain.StudentFee) error {
	return r.DB.CreateInBatches(fees, 100).Error
}

func (r *FeeRepository) FindPendingDuesByStudent(studentID uuid.UUID) ([]domain.StudentFee, error) {
	var fees []domain.StudentFee
	err := r.DB.Where("student_id = ? AND status != ?", studentID, domain.FeeStatusPaid).
		Order("due_date asc").
		Find(&fees).Error
	return fees, err
}

func (r *FeeRepository) FindTransactionsByStudent(studentID uuid.UUID) ([]domain.FeeTransaction, error) {
	var txns []domain.FeeTransaction
	// We currently don't have a direct link to StudentFee inside Transaction struct in memory,
	// but we can query it.
	// Ideally, join with StudentFee to get the Title of what was paid.
	err := r.DB.Where("student_id = ?", studentID).
		Order("transaction_date desc").
		Find(&txns).Error
	return txns, err
}

func (r *FeeRepository) GetFeeTitle(feeID uuid.UUID) string {
	var fee domain.StudentFee
	r.DB.Select("title").First(&fee, "id = ?", feeID)
	return fee.Title
}

func (r *FeeRepository) FindFeeWithLock(tx *gorm.DB, feeID uuid.UUID) (*domain.StudentFee, error) {
	var fee domain.StudentFee
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&fee, "id = ?", feeID).Error; err != nil {
		return nil, err
	}
	return &fee, nil
}

func (r *FeeRepository) CreateTransaction(tx *gorm.DB, txn *domain.FeeTransaction) error {
	return tx.Create(txn).Error
}

func (r *FeeRepository) UpdateStudentFee(tx *gorm.DB, fee *domain.StudentFee) error {
	return tx.Save(fee).Error
}
