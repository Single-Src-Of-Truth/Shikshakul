package controller

import (
	"net/http"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExamController struct {
	Service *service.ExamService
}

func NewExamController(svc *service.ExamService) *ExamController {
	return &ExamController{Service: svc}
}

func (ctrl *ExamController) CreateTerm(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.CreateExamTermRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	term, err := ctrl.Service.CreateTerm(tenantID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, "Exam term created successfully", term)
}

func (ctrl *ExamController) GetTerms(c *gin.Context) {
	yearID := c.Query("academic_year_id")
	if yearID == "" {
		response.Error(c, http.StatusBadRequest, "academic_year_id is required")
		return
	}

	terms, err := ctrl.Service.GetTermsByYear(yearID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, terms)
}

func (ctrl *ExamController) CreateSchedule(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.CreateExamScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	schedule, err := ctrl.Service.CreateSchedule(tenantID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, "Exam scheduled successfully", schedule)
}

func (ctrl *ExamController) GetSchedule(c *gin.Context) {
	termID := c.Query("term_id")
	classID := c.Query("class_id")

	if termID == "" || classID == "" {
		response.Error(c, http.StatusBadRequest, "term_id and class_id are required")
		return
	}

	schedules, err := ctrl.Service.GetSchedule(termID, classID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, schedules)
}

func (ctrl *ExamController) DeleteSchedule(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))

	if err := ctrl.Service.DeleteSchedule(tenantID, id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Exam schedule deleted", nil)
}
