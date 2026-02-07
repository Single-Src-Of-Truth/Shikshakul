package controller

import (
	"net/http"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttendanceController struct {
	Service *service.AttendanceService
}

func NewAttendanceController(svc *service.AttendanceService) *AttendanceController {
	return &AttendanceController{Service: svc}
}

func (ctrl *AttendanceController) MarkAttendance(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	teacherID := c.MustGet("userID").(uuid.UUID)

	var req dto.MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctrl.Service.MarkAttendance(tenantID, teacherID, req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Attendance marked successfully", nil)
}

func (ctrl *AttendanceController) GetClassRegister(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	sectionID := c.Param("section_id")
	date := c.Query("date")

	if date == "" {
		response.Error(c, http.StatusBadRequest, "Date parameter is required (YYYY-MM-DD)")
		return
	}

	register, err := ctrl.Service.GetClassRegister(tenantID, sectionID, date)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, register)
}

func (ctrl *AttendanceController) GetStudentHistory(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	studentID := uuid.MustParse(c.Param("student_id"))

	start := c.Query("start")
	end := c.Query("end")

	history, err := ctrl.Service.GetStudentHistory(tenantID, studentID, start, end)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, history)
}
