package routes

import (
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/bootstrap"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, app *bootstrap.Application) {
	r.GET("/", app.HealthCtrl.HealthCheck)
	r.GET("/health", app.HealthCtrl.HealthCheck)

	api := r.Group("/api/v1/academics")
	api.Use(middleware.AuthMiddleware())
	{
		setupRoutes(api, app)
		studentRoutes(api, app)
		teacherRoutes(api, app)
		feeRoutes(api, app)
		timetableRoutes(api, app)
		examRoutes(api, app)
		resultRoutes(api, app)
		attendanceRoutes(api, app)
		certificateRoutes(api, app)
		calendarRoutes(api, app)
	}
}

func setupRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/setup")
	g.POST("/academic-years", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.CreateYear)
	g.GET("/academic-years", app.AcademicCtrl.GetAllYears)
	g.GET("/academic-years/current", app.AcademicCtrl.GetCurrentYear)
	g.PUT("/academic-years/:id", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.UpdateYear)
	g.DELETE("/academic-years/:id", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.DeleteYear)

	g.POST("/classes", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.CreateClass)
	g.GET("/classes", app.AcademicCtrl.GetAllClasses)
	g.PUT("/classes/:class_id", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.UpdateClass)
	g.DELETE("/classes/:class_id", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.DeleteClass)

	g.POST("/classes/:class_id/sections", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.CreateSection)
	g.GET("/classes/:class_id/sections", app.AcademicCtrl.GetSectionsByClass)
	g.PUT("/sections/:id", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.UpdateSection)
	g.DELETE("/sections/:id", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.DeleteSection)

	g.POST("/subjects", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.CreateSubject)
	g.GET("/subjects", app.AcademicCtrl.GetAllSubjects)
	g.PUT("/subjects/:id", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.UpdateSubject)
	g.DELETE("/subjects/:id", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.DeleteSubject)

	g.POST("/classes/:class_id/subjects", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.AssignSubjectsToClass)
	g.DELETE("/classes/:class_id/subjects/:subject_id", middleware.RoleGuard("ADMIN"), app.AcademicCtrl.UnassignSubject)
}

func studentRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/students")
	g.POST("/schema", middleware.RoleGuard("ADMIN"), app.StudentCtrl.CreateSchema)
	g.GET("/schema", app.StudentCtrl.GetActiveSchema)
	g.DELETE("/schema", middleware.RoleGuard("ADMIN"), app.StudentCtrl.DeleteSchema)
	g.POST("/setup/admission-sequence", middleware.RoleGuard("ADMIN"), app.StudentCtrl.ConfigureSequence)

	g.POST("/onboard", app.StudentCtrl.SubmitApplication)
	g.GET("", middleware.RoleGuard("ADMIN", "MANAGEMENT"), app.StudentCtrl.GetAllStudents)
	g.GET("/:id", middleware.RoleGuard("ADMIN", "MANAGEMENT"), app.StudentCtrl.GetStudentByID)
	g.POST("/:id/approve", middleware.RoleGuard("ADMIN"), app.StudentCtrl.ApproveStudent)
	g.PUT("/:id", middleware.RoleGuard("ADMIN"), app.StudentCtrl.UpdateStudent)
	g.DELETE("/:id", middleware.RoleGuard("ADMIN"), app.StudentCtrl.DeleteStudent)
	g.GET("/active", middleware.RoleGuard("ADMIN", "TEACHER", "MANAGEMENT"), app.StudentCtrl.GetActiveStudents)
}

func teacherRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/teachers")
	g.POST("/schema", middleware.RoleGuard("ADMIN"), app.TeacherCtrl.CreateSchema)
	g.GET("/schema", app.TeacherCtrl.GetActiveSchema)
	g.DELETE("/schema", middleware.RoleGuard("ADMIN"), app.TeacherCtrl.DeleteSchema)

	g.POST("/onboard", app.TeacherCtrl.OnboardTeacher)
	g.GET("", middleware.RoleGuard("ADMIN", "MANAGEMENT"), app.TeacherCtrl.GetAllTeachers)
	g.GET("/:id", middleware.RoleGuard("ADMIN", "MANAGEMENT"), app.TeacherCtrl.GetTeacherByID)
	g.PUT("/:id", middleware.RoleGuard("ADMIN"), app.TeacherCtrl.UpdateTeacher)
	g.DELETE("/:id", middleware.RoleGuard("ADMIN"), app.TeacherCtrl.DeleteTeacher)
	g.POST("/:id/approve", middleware.RoleGuard("ADMIN"), app.TeacherCtrl.ApproveTeacher)

	rg.POST("/sections/:section_id/class-teacher", middleware.RoleGuard("ADMIN"), app.TeacherCtrl.AssignClassTeacher)
	rg.POST("/sections/:section_id/subject-teachers", middleware.RoleGuard("ADMIN"), app.TeacherCtrl.AssignSubjectTeacher)
	rg.GET("/sections/:section_id/allocations", app.TeacherCtrl.GetSectionAllocations)
}

func feeRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/fees")
	g.POST("/heads", middleware.RoleGuard("ADMIN"), app.FeeCtrl.CreateFeeHead)
	g.GET("/heads", middleware.RoleGuard("ADMIN"), app.FeeCtrl.GetAllFeeHeads)
	g.PUT("/heads/:id", middleware.RoleGuard("ADMIN"), app.FeeCtrl.UpdateFeeHead)
	g.DELETE("/heads/:id", middleware.RoleGuard("ADMIN"), app.FeeCtrl.DeleteFeeHead)

	g.POST("/structure", middleware.RoleGuard("ADMIN"), app.FeeCtrl.CreateFeeStructure)
	g.GET("/structure", middleware.RoleGuard("ADMIN"), app.FeeCtrl.GetStructuresByClass)
	g.PUT("/structure/:id", middleware.RoleGuard("ADMIN"), app.FeeCtrl.UpdateFeeStructure)
	g.DELETE("/structure/:id", middleware.RoleGuard("ADMIN"), app.FeeCtrl.DeleteFeeStructure)

	g.POST("/generate-demands", middleware.RoleGuard("ADMIN"), app.FeeCtrl.GenerateDemands)
	g.POST("/collect", middleware.RoleGuard("ADMIN"), app.FeeCtrl.CollectFee)
	g.GET("/students/:student_id/dues", app.FeeCtrl.GetStudentDues)
	g.GET("/students/:student_id/history", app.FeeCtrl.GetStudentHistory)
}

func timetableRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/timetables")
	g.POST("", middleware.RoleGuard("ADMIN", "MANAGEMENT"), app.TimetableCtrl.CreateRoutine)
	g.GET("", app.TimetableCtrl.GetTimetable)
	g.DELETE("/:id", middleware.RoleGuard("ADMIN"), app.TimetableCtrl.DeleteRoutine)
}

func examRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/exams")
	g.POST("/terms", middleware.RoleGuard("ADMIN"), app.ExamCtrl.CreateTerm)
	g.GET("/terms", app.ExamCtrl.GetTerms)
	g.POST("/schedules", middleware.RoleGuard("ADMIN"), app.ExamCtrl.CreateSchedule)
	g.GET("/schedules", app.ExamCtrl.GetSchedule)
	g.DELETE("/schedules/:id", middleware.RoleGuard("ADMIN"), app.ExamCtrl.DeleteSchedule)
	g.GET("/schedules/:schedule_id/marks", middleware.RoleGuard("TEACHER", "ADMIN"), app.EvalCtrl.GetMarkSheet)
	g.POST("/marks", middleware.RoleGuard("TEACHER", "ADMIN"), app.EvalCtrl.SubmitMarks)
}

func resultRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/results")
	g.POST("/generate", middleware.RoleGuard("ADMIN"), app.ResultCtrl.GenerateResults)
	g.POST("/publish", middleware.RoleGuard("ADMIN"), app.ResultCtrl.PublishResults)
	g.GET("/students/:student_id", app.ResultCtrl.GetReportCard)
}

func attendanceRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/attendance")
	g.POST("", middleware.RoleGuard("TEACHER", "ADMIN"), app.AttCtrl.MarkAttendance)
	g.GET("/section/:section_id", middleware.RoleGuard("TEACHER", "ADMIN"), app.AttCtrl.GetClassRegister)
	g.GET("/student/:student_id", app.AttCtrl.GetStudentHistory)
}

func certificateRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/certificates")
	g.GET("/id-cards", middleware.RoleGuard("ADMIN", "TEACHER"), app.CertCtrl.GetIDCards)
	g.POST("/tc", middleware.RoleGuard("ADMIN"), app.CertCtrl.GenerateTC)
	g.POST("/bonafide", middleware.RoleGuard("ADMIN", "TEACHER"), app.CertCtrl.GenerateBonafide)
}

func calendarRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/calendar")
	g.POST("/events", middleware.RoleGuard("ADMIN"), app.CalendarCtrl.CreateEvent)
	g.GET("/events", app.CalendarCtrl.GetEvents)
	g.DELETE("/events/:id", middleware.RoleGuard("ADMIN"), app.CalendarCtrl.DeleteEvent)
	g.PUT("/events/:id", middleware.RoleGuard("ADMIN"), app.CalendarCtrl.UpdateEvent)
}
