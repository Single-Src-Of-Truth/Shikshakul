package bootstrap

import (
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/controller"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/database"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/repository"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/service"
)

type Application struct {
	AcademicCtrl  *controller.AcademicController
	StudentCtrl   *controller.StudentController
	TeacherCtrl   *controller.TeacherController
	FeeCtrl       *controller.FeeController
	TimetableCtrl *controller.TimetableController
	ExamCtrl      *controller.ExamController
	EvalCtrl      *controller.EvaluationController
	ResultCtrl    *controller.ResultController
	AttCtrl       *controller.AttendanceController
	CertCtrl      *controller.CertificateController
	CalendarCtrl  *controller.CalendarController
	HealthCtrl    *controller.HealthController
}

func InitializeApp() *Application {
	database.ConnectDB()
	database.ConnectRedis()
	db := database.DB

	acadRepo := repository.NewAcademicRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)
	feeRepo := repository.NewFeeRepository(db)
	timetableRepo := repository.NewTimetableRepository(db)
	examRepo := repository.NewExamRepository(db)
	evalRepo := repository.NewEvaluationRepository(db)
	resultRepo := repository.NewResultRepository(db)
	attRepo := repository.NewAttendanceRepository(db)
	certRepo := repository.NewCertificateRepository(db)
	calRepo := repository.NewCalendarRepository(db)

	acadService := service.NewAcademicService(acadRepo)
	studentService := service.NewStudentService(studentRepo)
	teacherService := service.NewTeacherService(teacherRepo, acadRepo)
	feeService := service.NewFeeService(feeRepo, studentRepo)
	timetableService := service.NewTimetableService(timetableRepo)
	examService := service.NewExamService(examRepo)
	evalService := service.NewEvaluationService(evalRepo, teacherRepo)
	resultService := service.NewResultService(resultRepo, studentRepo)
	attService := service.NewAttendanceService(attRepo, studentRepo, acadRepo)
	certService := service.NewCertificateService(certRepo, studentRepo, acadRepo)
	calService := service.NewCalendarService(calRepo, acadRepo)

	return &Application{
		AcademicCtrl:  controller.NewAcademicController(acadService),
		StudentCtrl:   controller.NewStudentController(studentService),
		TeacherCtrl:   controller.NewTeacherController(teacherService),
		FeeCtrl:       controller.NewFeeController(feeService),
		TimetableCtrl: controller.NewTimetableController(timetableService),
		ExamCtrl:      controller.NewExamController(examService),
		EvalCtrl:      controller.NewEvaluationController(evalService),
		ResultCtrl:    controller.NewResultController(resultService),
		AttCtrl:       controller.NewAttendanceController(attService),
		CertCtrl:      controller.NewCertificateController(certService),
		CalendarCtrl:  controller.NewCalendarController(calService),
		HealthCtrl:    controller.NewHealthController(),
	}
}
