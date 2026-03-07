package routes

import (
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/bootstrap"
	acadMiddleware "github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/middleware"
	coreMiddleware "github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/middleware"
	"github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/versioning"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, app *bootstrap.Application) {
	r.Use(versioning.EnforceVersion())

	r.GET("/", app.HealthCtrl.HealthCheck)
	r.GET("/health", app.HealthCtrl.HealthCheck)

	protected := r.Group("")

	protected.Use(coreMiddleware.RequireGatewayAuth())
	protected.Use(acadMiddleware.UUIDBridge())
	{
		setupRoutes(protected, app)
		studentRoutes(protected, app)
		teacherRoutes(protected, app)
		feeRoutes(protected, app)
		timetableRoutes(protected, app)
		examRoutes(protected, app)
		resultRoutes(protected, app)
		attendanceRoutes(protected, app)
		certificateRoutes(protected, app)
		calendarRoutes(protected, app)
	}
}

func setupRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/setup")

	g.GET("/academic-years", coreMiddleware.RequireGatewayPermission("acad:setup:read"), app.AcademicCtrl.GetAllYears)
	g.GET("/academic-years/current", coreMiddleware.RequireGatewayPermission("acad:setup:read"), app.AcademicCtrl.GetCurrentYear)
	g.POST("/academic-years", coreMiddleware.RequireGatewayPermission("acad:setup:create"), app.AcademicCtrl.CreateYear)
	g.PUT("/academic-years/:id", coreMiddleware.RequireGatewayPermission("acad:setup:update"), app.AcademicCtrl.UpdateYear)
	g.DELETE("/academic-years/:id", coreMiddleware.RequireGatewayPermission("acad:setup:delete"), app.AcademicCtrl.DeleteYear)

	g.GET("/classes", coreMiddleware.RequireGatewayPermission("acad:setup:read"), app.AcademicCtrl.GetAllClasses)
	g.POST("/classes", coreMiddleware.RequireGatewayPermission("acad:setup:create"), app.AcademicCtrl.CreateClass)
	g.PUT("/classes/:class_id", coreMiddleware.RequireGatewayPermission("acad:setup:update"), app.AcademicCtrl.UpdateClass)
	g.DELETE("/classes/:class_id", coreMiddleware.RequireGatewayPermission("acad:setup:delete"), app.AcademicCtrl.DeleteClass)

	g.GET("/classes/:class_id/sections", coreMiddleware.RequireGatewayPermission("acad:setup:read"), app.AcademicCtrl.GetSectionsByClass)
	g.POST("/classes/:class_id/sections", coreMiddleware.RequireGatewayPermission("acad:setup:create"), app.AcademicCtrl.CreateSection)
	g.PUT("/sections/:id", coreMiddleware.RequireGatewayPermission("acad:setup:update"), app.AcademicCtrl.UpdateSection)
	g.DELETE("/sections/:id", coreMiddleware.RequireGatewayPermission("acad:setup:delete"), app.AcademicCtrl.DeleteSection)

	g.GET("/subjects", coreMiddleware.RequireGatewayPermission("acad:setup:read"), app.AcademicCtrl.GetAllSubjects)
	g.POST("/subjects", coreMiddleware.RequireGatewayPermission("acad:setup:create"), app.AcademicCtrl.CreateSubject)
	g.PUT("/subjects/:id", coreMiddleware.RequireGatewayPermission("acad:setup:update"), app.AcademicCtrl.UpdateSubject)
	g.DELETE("/subjects/:id", coreMiddleware.RequireGatewayPermission("acad:setup:delete"), app.AcademicCtrl.DeleteSubject)

	g.POST("/classes/:class_id/subjects", coreMiddleware.RequireGatewayPermission("acad:setup:create"), app.AcademicCtrl.AssignSubjectsToClass)
	g.DELETE("/classes/:class_id/subjects/:subject_id", coreMiddleware.RequireGatewayPermission("acad:setup:delete"), app.AcademicCtrl.UnassignSubject)
}

func studentRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/students")

	g.GET("/schema", coreMiddleware.RequireGatewayPermission("acad:students:read"), app.StudentCtrl.GetActiveSchema)
	g.POST("/schema", coreMiddleware.RequireGatewayPermission("acad:students:create"), app.StudentCtrl.CreateSchema)
	g.DELETE("/schema", coreMiddleware.RequireGatewayPermission("acad:students:delete"), app.StudentCtrl.DeleteSchema)
	g.POST("/setup/admission-sequence", coreMiddleware.RequireGatewayPermission("acad:students:create"), app.StudentCtrl.ConfigureSequence)

	g.POST("/onboard", coreMiddleware.RequireGatewayPermission("acad:students:create"), app.StudentCtrl.SubmitApplication)
	g.GET("", coreMiddleware.RequireGatewayPermission("acad:students:read"), app.StudentCtrl.GetAllStudents)
	g.GET("/active", coreMiddleware.RequireGatewayPermission("acad:students:read"), app.StudentCtrl.GetActiveStudents)
	g.GET("/:id", coreMiddleware.RequireGatewayPermission("acad:students:read"), app.StudentCtrl.GetStudentByID)

	g.POST("/:id/approve", coreMiddleware.RequireGatewayPermission("acad:students:execute"), app.StudentCtrl.ApproveStudent)
	g.PUT("/:id", coreMiddleware.RequireGatewayPermission("acad:students:update"), app.StudentCtrl.UpdateStudent)
	g.DELETE("/:id", coreMiddleware.RequireGatewayPermission("acad:students:delete"), app.StudentCtrl.DeleteStudent)
}

func teacherRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/teachers")

	g.GET("/schema", coreMiddleware.RequireGatewayPermission("acad:teachers:read"), app.TeacherCtrl.GetActiveSchema)
	g.POST("/schema", coreMiddleware.RequireGatewayPermission("acad:teachers:create"), app.TeacherCtrl.CreateSchema)
	g.DELETE("/schema", coreMiddleware.RequireGatewayPermission("acad:teachers:delete"), app.TeacherCtrl.DeleteSchema)

	g.POST("/onboard", coreMiddleware.RequireGatewayPermission("acad:teachers:create"), app.TeacherCtrl.OnboardTeacher)
	g.GET("", coreMiddleware.RequireGatewayPermission("acad:teachers:read"), app.TeacherCtrl.GetAllTeachers)
	g.GET("/:id", coreMiddleware.RequireGatewayPermission("acad:teachers:read"), app.TeacherCtrl.GetTeacherByID)
	g.PUT("/:id", coreMiddleware.RequireGatewayPermission("acad:teachers:update"), app.TeacherCtrl.UpdateTeacher)
	g.DELETE("/:id", coreMiddleware.RequireGatewayPermission("acad:teachers:delete"), app.TeacherCtrl.DeleteTeacher)
	g.POST("/:id/approve", coreMiddleware.RequireGatewayPermission("acad:teachers:execute"), app.TeacherCtrl.ApproveTeacher)

	g.POST("/sections/:section_id/class-teacher", coreMiddleware.RequireGatewayPermission("acad:teachers:allocations:create"), app.TeacherCtrl.AssignClassTeacher)
	g.POST("/sections/:section_id/subject-teachers", coreMiddleware.RequireGatewayPermission("acad:teachers:allocations:create"), app.TeacherCtrl.AssignSubjectTeacher)
	g.GET("/sections/:section_id/allocations", coreMiddleware.RequireGatewayPermission("acad:teachers:allocations:read"), app.TeacherCtrl.GetSectionAllocations)
}

func feeRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/fees")

	g.POST("/heads", coreMiddleware.RequireGatewayPermission("acad:fees:create"), app.FeeCtrl.CreateFeeHead)
	g.GET("/heads", coreMiddleware.RequireGatewayPermission("acad:fees:read"), app.FeeCtrl.GetAllFeeHeads)
	g.PUT("/heads/:id", coreMiddleware.RequireGatewayPermission("acad:fees:update"), app.FeeCtrl.UpdateFeeHead)
	g.DELETE("/heads/:id", coreMiddleware.RequireGatewayPermission("acad:fees:delete"), app.FeeCtrl.DeleteFeeHead)

	g.POST("/structure", coreMiddleware.RequireGatewayPermission("acad:fees:structures:create"), app.FeeCtrl.CreateFeeStructure)
	g.GET("/structure", coreMiddleware.RequireGatewayPermission("acad:fees:structures:read"), app.FeeCtrl.GetStructuresByClass)
	g.PUT("/structure/:id", coreMiddleware.RequireGatewayPermission("acad:fees:structures:update"), app.FeeCtrl.UpdateFeeStructure)
	g.DELETE("/structure/:id", coreMiddleware.RequireGatewayPermission("acad:fees:structures:delete"), app.FeeCtrl.DeleteFeeStructure)

	g.POST("/generate-demands", coreMiddleware.RequireGatewayPermission("acad:fees:demands:execute"), app.FeeCtrl.GenerateDemands)
	g.POST("/collect", coreMiddleware.RequireGatewayPermission("acad:fees:payments:create"), app.FeeCtrl.CollectFee)
	g.GET("/students/:student_id/dues", coreMiddleware.RequireGatewayPermission("acad:fees:payments:read"), app.FeeCtrl.GetStudentDues)
	g.GET("/students/:student_id/history", coreMiddleware.RequireGatewayPermission("acad:fees:payments:read"), app.FeeCtrl.GetStudentHistory)
}

func timetableRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/timetables")

	g.POST("", coreMiddleware.RequireGatewayPermission("acad:timetables:create"), app.TimetableCtrl.CreateRoutine)
	g.GET("", coreMiddleware.RequireGatewayPermission("acad:timetables:read"), app.TimetableCtrl.GetTimetable)
	g.DELETE("/:id", coreMiddleware.RequireGatewayPermission("acad:timetables:delete"), app.TimetableCtrl.DeleteRoutine)
}

func examRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/exams")

	g.POST("/terms", coreMiddleware.RequireGatewayPermission("acad:exams:terms:create"), app.ExamCtrl.CreateTerm)
	g.GET("/terms", coreMiddleware.RequireGatewayPermission("acad:exams:terms:read"), app.ExamCtrl.GetTerms)

	g.POST("/schedules", coreMiddleware.RequireGatewayPermission("acad:exams:schedules:create"), app.ExamCtrl.CreateSchedule)
	g.GET("/schedules", coreMiddleware.RequireGatewayPermission("acad:exams:schedules:read"), app.ExamCtrl.GetSchedule)
	g.DELETE("/schedules/:id", coreMiddleware.RequireGatewayPermission("acad:exams:schedules:delete"), app.ExamCtrl.DeleteSchedule)

	g.GET("/schedules/:schedule_id/marks", coreMiddleware.RequireGatewayPermission("acad:exams:marks:read"), app.EvalCtrl.GetMarkSheet)
	g.POST("/marks", coreMiddleware.RequireGatewayPermission("acad:exams:marks:create"), app.EvalCtrl.SubmitMarks)
}

func resultRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/results")

	g.POST("/generate", coreMiddleware.RequireGatewayPermission("acad:results:execute"), app.ResultCtrl.GenerateResults)
	g.POST("/publish", coreMiddleware.RequireGatewayPermission("acad:results:execute"), app.ResultCtrl.PublishResults)
	g.GET("/students/:student_id", coreMiddleware.RequireGatewayPermission("acad:results:read"), app.ResultCtrl.GetReportCard)
}

func attendanceRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/attendance")

	g.POST("", coreMiddleware.RequireGatewayPermission("acad:attendance:create"), app.AttCtrl.MarkAttendance)
	g.GET("/section/:section_id", coreMiddleware.RequireGatewayPermission("acad:attendance:read"), app.AttCtrl.GetClassRegister)
	g.GET("/student/:student_id", coreMiddleware.RequireGatewayPermission("acad:attendance:read"), app.AttCtrl.GetStudentHistory)
}

func certificateRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/certificates")

	g.GET("/id-cards", coreMiddleware.RequireGatewayPermission("acad:certificates:read"), app.CertCtrl.GetIDCards)
	g.POST("/tc", coreMiddleware.RequireGatewayPermission("acad:certificates:create"), app.CertCtrl.GenerateTC)
	g.POST("/bonafide", coreMiddleware.RequireGatewayPermission("acad:certificates:create"), app.CertCtrl.GenerateBonafide)
}

func calendarRoutes(rg *gin.RouterGroup, app *bootstrap.Application) {
	g := rg.Group("/calendar")

	g.POST("/events", coreMiddleware.RequireGatewayPermission("acad:calendar:create"), app.CalendarCtrl.CreateEvent)
	g.GET("/events", coreMiddleware.RequireGatewayPermission("acad:calendar:read"), app.CalendarCtrl.GetEvents)
	g.DELETE("/events/:id", coreMiddleware.RequireGatewayPermission("acad:calendar:delete"), app.CalendarCtrl.DeleteEvent)
	g.PUT("/events/:id", coreMiddleware.RequireGatewayPermission("acad:calendar:update"), app.CalendarCtrl.UpdateEvent)
}
