package database

import (
	"log"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/domain"
)

func Migrate() {
	err := DB.AutoMigrate(
		&domain.AcademicYear{},
		&domain.Class{},
		&domain.Section{},
		&domain.Subject{},
		&domain.ClassSubject{},
		&domain.AdmissionSchema{},
		&domain.Student{},
		&domain.TeacherSchema{},
		&domain.Teacher{},
		&domain.SubjectAllocation{},
		&domain.FeeHead{},
		&domain.FeeStructure{},
		&domain.StudentFee{},
		&domain.FeeTransaction{},
		&domain.WeeklyRoutine{},
		&domain.ExamTerm{},
		&domain.ExamSchedule{},
		&domain.StudentMark{},
		&domain.ExamResult{},
		&domain.StudentAttendance{},
		&domain.CertificateRecord{},
		&domain.CalendarEvent{},
		&domain.AdmissionSequence{},
	)

	if err != nil {
		log.Fatal("Migration Failed:", err)
	}
	log.Println("Database Migration Completed")
}
