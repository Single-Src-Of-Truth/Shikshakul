package controller

import (
	"net/http"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EvaluationController struct {
	Service *service.EvaluationService
}

func NewEvaluationController(svc *service.EvaluationService) *EvaluationController {
	return &EvaluationController{Service: svc}
}

func (ctrl *EvaluationController) GetMarkSheet(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	scheduleID := c.Param("schedule_id")

	sheet, err := ctrl.Service.GetMarkSheet(tenantID, scheduleID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, sheet)
}

func (ctrl *EvaluationController) SubmitMarks(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	teacherID := c.MustGet("userID").(uuid.UUID)

	var req dto.SubmitMarksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := ctrl.Service.SubmitMarks(tenantID, teacherID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Marks saved successfully", nil)
}
