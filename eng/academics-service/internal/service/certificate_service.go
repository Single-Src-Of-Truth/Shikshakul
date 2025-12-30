package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/google/uuid"
)

type CertificateService struct {
	Repo        *repository.CertificateRepository
	StudentRepo *repository.StudentRepository
	AcadRepo    *repository.AcademicRepository
}

func NewCertificateService(repo *repository.CertificateRepository, sRepo *repository.StudentRepository, aRepo *repository.AcademicRepository) *CertificateService {
	return &CertificateService{Repo: repo, StudentRepo: sRepo, AcadRepo: aRepo}
}

func (s *CertificateService) GetIDCardData(tenantID uuid.UUID, classIDStr string) ([]dto.IDCardResponse, error) {
	classID := uuid.MustParse(classIDStr)
	students, err := s.StudentRepo.FindActiveStudentsByClass(classID)
	if err != nil {
		return nil, err
	}

	var cards []dto.IDCardResponse
	for _, st := range students {
		profile := make(map[string]interface{})
		_ = json.Unmarshal(st.ProfileData, &profile)

		fatherName, _ := profile["father_name"].(string)
		bloodGroup, _ := profile["blood_group"].(string)
		address, _ := profile["address"].(string)
		photo, _ := profile["photo_url"].(string)

		cards = append(cards, dto.IDCardResponse{
			StudentID:     st.ID.String(),
			AdmissionNo:   st.AdmissionNo,
			FullName:      st.FirstName + " " + st.LastName,
			ClassDetails:  fmt.Sprintf("%s - %s", st.Class.Name, st.Section.Name),
			DOB:           getString(profile, "dob"),
			BloodGroup:    bloodGroup,
			FatherName:    fatherName,
			ContactNumber: st.Mobile,
			Address:       address,
			PhotoURL:      photo,
		})
	}
	return cards, nil
}

func (s *CertificateService) GenerateTC(tenantID, issuerID uuid.UUID, req dto.IssueTCRequest) (*dto.TCResponse, error) {
	studentID := uuid.MustParse(req.StudentID)

	student, err := s.Repo.GetStudentFullDetails(studentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	certNo := s.Repo.GenerateCertificateNo(domain.CertTypeTC)
	record := &domain.CertificateRecord{
		TenantID:      tenantID,
		StudentID:     studentID,
		Type:          domain.CertTypeTC,
		CertificateNo: certNo,
		IssueDate:     time.Now(),
		Remarks:       req.Reason,
		IssuedBy:      issuerID,
	}

	if err := s.Repo.CreateRecord(record); err != nil {
		return nil, err
	}

	if req.MarkInactive {
		_ = s.Repo.UpdateStudentStatus(studentID, domain.StatusInactive)
	}

	profile := make(map[string]interface{})
	_ = json.Unmarshal(student.ProfileData, &profile)

	return &dto.TCResponse{
		CertificateNo: certNo,
		IssueDate:     time.Now().Format("02-01-2006"),
		StudentName:   student.FirstName + " " + student.LastName,
		FatherName:    getString(profile, "father_name"),
		MotherName:    getString(profile, "mother_name"),
		AdmissionNo:   student.AdmissionNo,
		DOB:           getString(profile, "dob"),
		Nationality:   getString(profile, "nationality"),
		LastClass:     student.Class.Name,
		ResultStatus:  "Passed", // Logic needed if based on marks
		Reason:        req.Reason,
		Conduct:       req.Conduct,
	}, nil
}

func (s *CertificateService) GenerateBonafide(tenantID, issuerID uuid.UUID, req dto.IssueBonafideRequest) (*dto.BonafideResponse, error) {
	studentID := uuid.MustParse(req.StudentID)
	student, err := s.Repo.GetStudentFullDetails(studentID)
	if err != nil {
		return nil, err
	}

	certNo := s.Repo.GenerateCertificateNo(domain.CertTypeBonafide)
	record := &domain.CertificateRecord{
		TenantID: tenantID, StudentID: studentID, Type: domain.CertTypeBonafide,
		CertificateNo: certNo, IssueDate: time.Now(), Remarks: req.Purpose, IssuedBy: issuerID,
	}
	s.Repo.CreateRecord(record)

	// Fetch Year Name (Assuming you can get it, or student has YearID)
	// For simplicity, hardcoding or fetching via repo if needed
	yearName := "2024-2025"

	return &dto.BonafideResponse{
		CertificateNo: certNo,
		IssueDate:     time.Now().Format("02-01-2006"),
		StudentName:   student.FirstName + " " + student.LastName,
		AdmissionNo:   student.AdmissionNo,
		ClassDetails:  student.Class.Name,
		AcademicYear:  yearName,
		Purpose:       req.Purpose,
	}, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
