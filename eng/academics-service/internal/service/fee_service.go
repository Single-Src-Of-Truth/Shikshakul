package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/database"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
)

type FeeService struct {
	Repo        *repository.FeeRepository
	StudentRepo *repository.StudentRepository
}

func NewFeeService(repo *repository.FeeRepository, studRepo *repository.StudentRepository) *FeeService {
	return &FeeService{Repo: repo, StudentRepo: studRepo}
}

func (s *FeeService) CreateFeeHead(tenantID uuid.UUID, req dto.CreateFeeHeadRequest) (*dto.FeeHeadResponse, error) {
	head := &domain.FeeHead{
		TenantID: tenantID,
		Name:     req.Name,
		Type:     req.Type,
	}
	if err := s.Repo.CreateFeeHead(head); err != nil {
		return nil, err
	}
	return &dto.FeeHeadResponse{ID: head.ID, Name: head.Name, Type: head.Type}, nil
}

func (s *FeeService) GetAllFeeHeads(tenantID uuid.UUID) ([]dto.FeeHeadResponse, error) {
	heads, err := s.Repo.GetAllFeeHeads(tenantID)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.FeeHeadResponse, len(heads))
	for i, h := range heads {
		dtos[i] = dto.FeeHeadResponse{ID: h.ID, Name: h.Name, Type: h.Type}
	}
	return dtos, nil
}

func (s *FeeService) UpdateFeeHead(tenantID, id uuid.UUID, req dto.UpdateFeeHeadRequest) (*dto.FeeHeadResponse, error) {
	head, err := s.Repo.FindHeadByID(id)
	if err != nil || head.TenantID != tenantID {
		return nil, ErrNotFound
	}

	if req.Name != "" {
		head.Name = req.Name
	}
	if req.Type != "" {
		head.Type = req.Type
	}

	if err := s.Repo.UpdateFeeHead(head); err != nil {
		return nil, err
	}
	return &dto.FeeHeadResponse{ID: head.ID, Name: head.Name, Type: head.Type}, nil
}

func (s *FeeService) DeleteFeeHead(tenantID, id uuid.UUID) error {
	head, err := s.Repo.FindHeadByID(id)
	if err != nil || head.TenantID != tenantID {
		return ErrNotFound
	}

	count, _ := s.Repo.CountStructuresByHead(id)
	if count > 0 {
		return fmt.Errorf("%w: fee head is used in fee structures", ErrConflict)
	}
	return s.Repo.DeleteFeeHead(id)
}

func (s *FeeService) CreateFeeStructure(tenantID uuid.UUID, req dto.CreateFeeStructureRequest) (*dto.FeeStructureResponse, error) {
	structure := &domain.FeeStructure{
		TenantID:       tenantID,
		AcademicYearID: uuid.MustParse(req.AcademicYearID),
		ClassID:        uuid.MustParse(req.ClassID),
		FeeHeadID:      uuid.MustParse(req.FeeHeadID),
		Amount:         req.Amount,
		Frequency:      domain.FeeFrequency(req.Frequency),
		DueDateDay:     req.DueDateDay,
	}
	if err := s.Repo.CreateFeeStructure(structure); err != nil {
		return nil, err
	}
	return mapStructureToDTO(structure), nil
}

func (s *FeeService) GetStructuresByClass(tenantID, classID, yearID uuid.UUID) ([]dto.FeeStructureResponse, error) {
	structures, err := s.Repo.FindStructuresByClass(classID, yearID)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.FeeStructureResponse, len(structures))
	for i, st := range structures {
		dtos[i] = *mapStructureToDTO(&st)
	}
	return dtos, nil
}

func (s *FeeService) UpdateFeeStructure(tenantID, id uuid.UUID, req dto.UpdateFeeStructureRequest) (*dto.FeeStructureResponse, error) {
	st, err := s.Repo.FindStructureByID(id)
	if err != nil || st.TenantID != tenantID {
		return nil, ErrNotFound
	}

	if req.Amount != 0 {
		st.Amount = req.Amount
	}
	if req.DueDateDay != 0 {
		st.DueDateDay = req.DueDateDay
	}

	if err := s.Repo.UpdateFeeStructure(st); err != nil {
		return nil, err
	}
	return mapStructureToDTO(st), nil
}

func (s *FeeService) DeleteFeeStructure(tenantID, id uuid.UUID) error {
	st, err := s.Repo.FindStructureByID(id)
	if err != nil || st.TenantID != tenantID {
		return ErrNotFound
	}

	count, _ := s.Repo.CountFeesByStructure(id)
	if count > 0 {
		return fmt.Errorf("%w: demands have already been generated for this structure", ErrConflict)
	}
	return s.Repo.DeleteFeeStructure(id)
}

func (s *FeeService) GenerateMonthlyDemands(tenantID uuid.UUID, req dto.GenerateDemandRequest) (int, error) {
	classUUID, _ := uuid.Parse(req.ClassID)
	yearUUID, _ := uuid.Parse(req.AcademicYearID)

	structures, err := s.Repo.FindStructuresByClass(classUUID, yearUUID)
	if err != nil || len(structures) == 0 {
		return 0, errors.New("no fee structure found")
	}

	students, err := s.StudentRepo.FindActiveStudentsByClass(classUUID)
	if err != nil {
		return 0, err
	}
	if len(students) == 0 {
		return 0, errors.New("no active students in class")
	}

	var feesToInsert []domain.StudentFee
	for _, student := range students {
		for _, st := range structures {
			if st.Frequency != domain.FeeFrequencyMonthly {
				continue
			}

			title := fmt.Sprintf("%s - %d/%d", st.FeeHead.Name, req.Month, req.Year)
			dueDate := time.Date(req.Year, time.Month(req.Month), st.DueDateDay, 0, 0, 0, 0, time.Local)

			feesToInsert = append(feesToInsert, domain.StudentFee{
				TenantID:       tenantID,
				StudentID:      student.ID,
				FeeStructureID: st.ID,
				Title:          title,
				DueDate:        dueDate,
				AmountDue:      st.Amount,
				Status:         domain.FeeStatusPending,
			})
		}
	}

	if len(feesToInsert) > 0 {
		err = s.Repo.BatchCreateStudentFees(feesToInsert)
	}
	return len(feesToInsert), err
}

func (s *FeeService) CollectFee(tenantID, collectorID uuid.UUID, req dto.CollectFeeRequest) (*domain.FeeTransaction, error) {
	feeID, _ := uuid.Parse(req.StudentFeeID)

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	feeBill, err := s.Repo.FindFeeWithLock(tx, feeID)
	if err != nil {
		tx.Rollback()
		return nil, ErrNotFound
	}

	pending := feeBill.AmountDue - feeBill.AmountPaid
	if req.Amount > pending {
		tx.Rollback()
		return nil, errors.New("amount exceeds pending dues")
	}

	receiptNo := fmt.Sprintf("RCPT-%d-%d", time.Now().Year(), rand.Intn(100000))
	txn := &domain.FeeTransaction{
		TenantID:        tenantID,
		StudentID:       feeBill.StudentID,
		StudentFeeID:    feeBill.ID,
		Amount:          req.Amount,
		PaymentMode:     domain.PaymentMode(req.PaymentMode),
		TransactionDate: time.Now(),
		ReceiptNo:       receiptNo,
		Remarks:         req.Remarks,
		CollectedBy:     collectorID,
	}

	if err := s.Repo.CreateTransaction(tx, txn); err != nil {
		tx.Rollback()
		return nil, err
	}

	feeBill.AmountPaid += req.Amount
	feeBill.Status = domain.FeeStatusPartial
	if feeBill.AmountPaid >= feeBill.AmountDue {
		feeBill.Status = domain.FeeStatusPaid
	}

	if err := s.Repo.UpdateStudentFee(tx, feeBill); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return txn, nil
}

func (s *FeeService) GetStudentDues(tenantID, studentID uuid.UUID) ([]dto.StudentDueResponse, error) {
	fees, err := s.Repo.FindPendingDuesByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var dtos []dto.StudentDueResponse
	for _, f := range fees {
		dtos = append(dtos, dto.StudentDueResponse{
			ID:        f.ID,
			Title:     f.Title,
			DueDate:   f.DueDate,
			AmountDue: f.AmountDue,
			Paid:      f.AmountPaid,
			Balance:   f.AmountDue - f.AmountPaid,
			Status:    string(f.Status),
		})
	}
	return dtos, nil
}

func (s *FeeService) GetStudentHistory(tenantID, studentID uuid.UUID) ([]dto.TransactionHistoryResponse, error) {
	txns, err := s.Repo.FindTransactionsByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var dtos []dto.TransactionHistoryResponse
	for _, t := range txns {
		// N+1 Query (acceptable for single student history)
		title := s.Repo.GetFeeTitle(t.StudentFeeID)

		dtos = append(dtos, dto.TransactionHistoryResponse{
			ID:              t.ID,
			ReceiptNo:       t.ReceiptNo,
			Amount:          t.Amount,
			PaymentMode:     string(t.PaymentMode),
			TransactionDate: t.TransactionDate,
			FeeTitle:        title,
		})
	}
	return dtos, nil
}

func mapStructureToDTO(st *domain.FeeStructure) *dto.FeeStructureResponse {
	return &dto.FeeStructureResponse{
		ID:          st.ID,
		ClassID:     st.ClassID,
		FeeHeadID:   st.FeeHeadID,
		FeeHeadName: st.FeeHead.Name,
		Amount:      st.Amount,
		Frequency:   string(st.Frequency),
	}
}
