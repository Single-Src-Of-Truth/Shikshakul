package controller

import (
	"net/http"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResultController struct {
	Service *service.ResultService
}

func NewResultController(svc *service.ResultService) *ResultController {
	return &ResultController{Service: svc}
}

func (ctrl *ResultController) GenerateResults(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.GenerateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctrl.Service.GenerateClassResults(tenantID, req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Results generated successfully", nil)
}

func (ctrl *ResultController) GetReportCard(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	studentID := uuid.MustParse(c.Param("student_id"))
	termID := c.Query("term_id")

	if termID == "" {
		response.Error(c, http.StatusBadRequest, "term_id is required")
		return
	}

	report, err := ctrl.Service.GetReportCard(tenantID, studentID, termID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, report)
}

func (ctrl *ResultController) PublishResults(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.PublishResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctrl.Service.PublishResults(tenantID, req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	msg := "Results published"
	if !req.Publish {
		msg = "Results unpublished"
	}
	response.SuccessWithMsg(c, msg, nil)
}
